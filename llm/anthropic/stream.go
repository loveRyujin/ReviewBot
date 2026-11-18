package anthropic

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/fatih/color"
	"github.com/loveRyujin/ReviewBot/ai"
)

func (c *Client) StreamChatCompletion(ctx context.Context, text string, handler ai.ChunkHandler) error {
	resp, err := c.request(ctx, text, true)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	color.Yellow("================Review Summary====================\n\n")
	tokens, err := readStream(ctx, resp.Body, handler)
	if err != nil {
		return err
	}
	color.Yellow("\n==================================================")
	color.Magenta(tokens.String())
	return nil
}

func readStream(ctx context.Context, body io.Reader, handler ai.ChunkHandler) (ai.TokenUsage, error) {
	reader := bufio.NewReader(body)
	var data []string
	var totals usage
	for {
		if err := ctx.Err(); err != nil {
			return ai.TokenUsage{}, err
		}
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			return ai.TokenUsage{}, fmt.Errorf("anthropic stream ended before message_stop: %w", io.ErrUnexpectedEOF)
		}
		if err != nil {
			return ai.TokenUsage{}, err
		}
		line = strings.TrimRight(line, "\r\n")
		if field, value, _ := strings.Cut(line, ":"); field == "data" {
			data = append(data, strings.TrimPrefix(value, " "))
		}
		if line != "" || len(data) == 0 {
			continue
		}

		var event struct {
			Type    string `json:"type"`
			Content block  `json:"content_block"`
			Delta   block  `json:"delta"`
			Message struct {
				Usage usage `json:"usage"`
			} `json:"message"`
			Usage json.RawMessage                `json:"usage"`
			Error struct{ Type, Message string } `json:"error"`
		}
		if err := json.Unmarshal([]byte(strings.Join(data, "\n")), &event); err != nil {
			return ai.TokenUsage{}, fmt.Errorf("decode anthropic stream event: %w", err)
		}
		data = data[:0]
		var chunk string
		switch event.Type {
		case "message_start":
			totals = event.Message.Usage
		case "content_block_start":
			if event.Content.Type == "text" {
				chunk = event.Content.Text
			}
		case "content_block_delta":
			if event.Delta.Type == "text_delta" {
				chunk = event.Delta.Text
			}
		case "message_delta":
			// Counters are cumulative. Preserve fields omitted by a delta.
			if event.Usage != nil {
				if err := json.Unmarshal(event.Usage, &totals); err != nil {
					return ai.TokenUsage{}, fmt.Errorf("decode anthropic stream usage: %w", err)
				}
			}
		case "message_stop":
			return totals.tokens(), nil
		case "error":
			return ai.TokenUsage{}, fmt.Errorf("anthropic stream %s: %s", event.Error.Type, event.Error.Message)
		}
		if chunk != "" {
			if err := handler(chunk); err != nil {
				return ai.TokenUsage{}, err
			}
		}
	}
}
