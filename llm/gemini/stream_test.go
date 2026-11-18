package gemini

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"iter"
	"testing"

	"github.com/fatih/color"
	"github.com/stretchr/testify/require"
	"google.golang.org/genai"
)

type streamModels struct {
	*genai.Models
	ctx    context.Context
	stream iter.Seq2[*genai.GenerateContentResponse, error]
}

func (m *streamModels) GenerateContentStream(ctx context.Context, _ string, _ []*genai.Content, _ *genai.GenerateContentConfig) iter.Seq2[*genai.GenerateContentResponse, error] {
	m.ctx = ctx
	return m.stream
}

func textCandidate(parts ...*genai.Part) *genai.Candidate {
	return &genai.Candidate{Content: &genai.Content{Parts: parts}}
}

func TestStreamChatCompletion(t *testing.T) {
	var output bytes.Buffer
	originalOutput := color.Output
	color.Output = &output
	t.Cleanup(func() { color.Output = originalOutput })

	chunks := []*genai.GenerateContentResponse{
		nil, {}, {Candidates: []*genai.Candidate{nil}}, {Candidates: []*genai.Candidate{{}}},
		{Candidates: []*genai.Candidate{textCandidate()}},
		{Candidates: []*genai.Candidate{textCandidate(&genai.Part{Text: "one"})},
			UsageMetadata: &genai.GenerateContentResponseUsageMetadata{TotalTokenCount: 3}},
		{Candidates: []*genai.Candidate{
			textCandidate(nil, &genai.Part{Thought: true, Text: "internal"}, &genai.Part{Text: " "},
				&genai.Part{FunctionCall: &genai.FunctionCall{Name: "tool"}}, &genai.Part{Text: "two"}),
			textCandidate(&genai.Part{Text: "alternative"}),
		}},
		{UsageMetadata: &genai.GenerateContentResponseUsageMetadata{
			PromptTokenCount: 5, CachedContentTokenCount: 2, CandidatesTokenCount: 3, TotalTokenCount: 8,
		}},
		{},
	}
	models := &streamModels{stream: func(yield func(*genai.GenerateContentResponse, error) bool) {
		for _, chunk := range chunks {
			if !yield(chunk, nil) {
				return
			}
		}
	}}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var received []string
	err := (&Client{models: models}).StreamChatCompletion(ctx, "review", func(s string) error {
		received = append(received, s)
		return nil
	})

	require.NoError(t, err)
	require.Same(t, ctx, models.ctx)
	require.Equal(t, []string{"one", " ", "two"}, received)
	require.Contains(t, output.String(), "Prompt tokens: 5 (CachedTokens: 2), Completion tokens: 3, Total tokens: 8")
}

func TestStreamChatCompletionBlocking(t *testing.T) {
	finish := func(reason genai.FinishReason) *genai.GenerateContentResponse {
		return &genai.GenerateContentResponse{Candidates: []*genai.Candidate{{FinishReason: reason}}}
	}
	stop := finish(genai.FinishReasonStop)
	empty := &genai.GenerateContentResponse{
		PromptFeedback: &genai.GenerateContentResponsePromptFeedback{},
		Candidates: []*genai.Candidate{
			textCandidate(&genai.Part{}, &genai.Part{Thought: true, Text: "internal"}),
		},
	}
	blockedPrompt := &genai.GenerateContentResponse{
		PromptFeedback: &genai.GenerateContentResponsePromptFeedback{BlockReason: genai.BlockedReasonOther},
	}
	text := &genai.GenerateContentResponse{
		Candidates: []*genai.Candidate{textCandidate(&genai.Part{Text: "review text"})},
	}

	for _, tt := range []struct {
		name       string
		chunks     []*genai.GenerateContentResponse
		wantText   []string
		wantReason string
	}{
		{name: "empty stream frames", chunks: []*genai.GenerateContentResponse{nil, {}, {Candidates: []*genai.Candidate{}}, empty, stop}},
		{name: "prompt blocked", chunks: []*genai.GenerateContentResponse{blockedPrompt, {}, stop}, wantReason: "OTHER"},
		{name: "safety", chunks: []*genai.GenerateContentResponse{empty, finish(genai.FinishReasonSafety)}, wantReason: "SAFETY"},
		{name: "recitation", chunks: []*genai.GenerateContentResponse{finish(genai.FinishReasonRecitation)}, wantReason: "RECITATION"},
		{name: "blocklist", chunks: []*genai.GenerateContentResponse{finish(genai.FinishReasonBlocklist)}, wantReason: "BLOCKLIST"},
		{name: "prohibited content", chunks: []*genai.GenerateContentResponse{finish(genai.FinishReasonProhibitedContent)}, wantReason: "PROHIBITED_CONTENT"},
		{name: "SPII", chunks: []*genai.GenerateContentResponse{finish(genai.FinishReasonSPII)}, wantReason: "SPII"},
		{name: "text after empty frames", chunks: []*genai.GenerateContentResponse{nil, {}, empty, text}, wantText: []string{"review text"}},
		{name: "stop after text", chunks: []*genai.GenerateContentResponse{text, stop}, wantText: []string{"review text"}},
		{name: "blocked finish after text", chunks: []*genai.GenerateContentResponse{text, finish(genai.FinishReasonSafety)}, wantText: []string{"review text"}},
		{name: "text after blocked finish", chunks: []*genai.GenerateContentResponse{finish(genai.FinishReasonSafety), text}, wantText: []string{"review text"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			models := &streamModels{stream: func(yield func(*genai.GenerateContentResponse, error) bool) {
				for _, chunk := range tt.chunks {
					if !yield(chunk, nil) {
						return
					}
				}
			}}
			var received []string
			err := (&Client{models: models}).StreamChatCompletion(context.Background(), "review", func(s string) error {
				received = append(received, s)
				return nil
			})
			require.Equal(t, tt.wantText, received)
			if tt.wantReason == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tt.wantReason)
			}
		})
	}
}

func TestStreamChatCompletionStops(t *testing.T) {
	handlerErr := errors.New("handler failed")
	upstreamErr := errors.New("upstream failed")
	for _, tt := range []struct {
		name                  string
		handlerErr, streamErr error
		cancelInHandler       bool
		wantErr               error
		wantText              []string
	}{
		{name: "wrapped EOF", streamErr: fmt.Errorf("stream: %w", io.EOF), wantText: []string{"a", "b"}},
		{name: "upstream error", streamErr: upstreamErr, wantErr: upstreamErr, wantText: []string{"a", "b"}},
		{name: "handler error", handlerErr: handlerErr, wantErr: handlerErr, wantText: []string{"a"}},
		{name: "handler cancels context", cancelInHandler: true, wantErr: context.Canceled, wantText: []string{"a"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			chunk := &genai.GenerateContentResponse{Candidates: []*genai.Candidate{
				textCandidate(&genai.Part{Text: "a"}, &genai.Part{Text: "b"}),
			}}
			models := &streamModels{stream: func(yield func(*genai.GenerateContentResponse, error) bool) {
				if !yield(chunk, nil) {
					return
				}
				if tt.streamErr == nil || yield(nil, tt.streamErr) {
					t.Fatal("stream consumed past stopping point")
				}
			}}
			var received []string
			err := (&Client{models: models}).StreamChatCompletion(ctx, "review", func(s string) error {
				received = append(received, s)
				if tt.cancelInHandler {
					cancel()
				}
				return tt.handlerErr
			})
			require.Equal(t, tt.wantErr, err)
			require.Equal(t, tt.wantText, received)
		})
	}
}

func TestStreamChatCompletionCancellation(t *testing.T) {
	for _, when := range []string{"before request", "before chunk", "while awaiting end"} {
		t.Run(when, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			models := &streamModels{stream: func(yield func(*genai.GenerateContentResponse, error) bool) {
				cancel()
				if when == "before chunk" {
					yield(&genai.GenerateContentResponse{Candidates: []*genai.Candidate{
						textCandidate(&genai.Part{Text: "late"}),
					}}, nil)
				}
			}}
			if when == "before request" {
				cancel()
			}
			err := (&Client{models: models}).StreamChatCompletion(ctx, "review", func(string) error {
				t.Fatal("text received after cancellation")
				return nil
			})
			require.ErrorIs(t, err, context.Canceled)
			if when == "before request" {
				require.Nil(t, models.ctx)
			}
		})
	}
}
