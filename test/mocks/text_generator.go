package mocks

import (
	"context"

	"github.com/loveRyujin/ReviewBot/ai"
	"github.com/stretchr/testify/mock"
)

// MockTextGenerator is a mock implementation of ai.TextGenerator for testing.
type MockTextGenerator struct {
	mock.Mock
}

// ChatCompletion mocks the ChatCompletion method.
func (m *MockTextGenerator) ChatCompletion(ctx context.Context, text string) (*ai.Response, error) {
	args := m.Called(ctx, text)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ai.Response), args.Error(1)
}

// StreamChatCompletion mocks the StreamChatCompletion method.
func (m *MockTextGenerator) StreamChatCompletion(ctx context.Context, text string, handler ai.ChunkHandler) error {
	args := m.Called(ctx, text, handler)
	return args.Error(0)
}
