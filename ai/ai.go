package ai

import (
	"context"
	"strconv"

	"github.com/sashabaranov/go-openai"
)

type TokenUsage struct {
	PromptTokens            int
	CompletionTokens        int
	TotalTokens             int
	PromptTokensDetails     *openai.PromptTokensDetails
	CompletionTokensDetails *openai.CompletionTokensDetails
}

func (u TokenUsage) String() string {
	s := "Prompt tokens: " + strconv.Itoa(u.PromptTokens)
	if u.PromptTokensDetails != nil && u.PromptTokensDetails.CachedTokens > 0 {
		s += " (CachedTokens: " + strconv.Itoa(u.PromptTokensDetails.CachedTokens) + ")"
	}
	s += ", Completion tokens: " + strconv.Itoa(u.CompletionTokens)
	if u.CompletionTokensDetails != nil && u.CompletionTokensDetails.ReasoningTokens > 0 {
		s += " (ReasoningTokens: " + strconv.Itoa(u.CompletionTokensDetails.ReasoningTokens) + ")"
	}
	s += ", Total tokens: " + strconv.Itoa(u.TotalTokens)
	return s
}

type Response struct {
	Text       string
	TokenUsage TokenUsage
}

type TextGenerator interface {
	ChatCompletion(ctx context.Context, text string) (*Response, error)
}
