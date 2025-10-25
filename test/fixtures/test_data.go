package fixtures

import "github.com/loveRyujin/ReviewBot/ai"

// Common test data and fixtures for testing.

// SampleDiff returns a sample git diff for testing.
func SampleDiff() string {
	return `diff --git a/example.go b/example.go
index 1234567..abcdefg 100644
--- a/example.go
+++ b/example.go
@@ -1,5 +1,8 @@
 package example
 
+// Add adds two integers.
 func Add(a, b int) int {
-	return a + b
+	// TODO: add validation
+	result := a + b
+	return result
 }
`
}

// SampleCommitMessage returns a sample commit message.
func SampleCommitMessage() string {
	return "feat: add integer addition function\n\nImplement Add function with basic validation TODO"
}

// SampleAIResponse returns a sample AI response.
func SampleAIResponse() *ai.Response {
	return &ai.Response{
		Text: SampleCommitMessage(),
		TokenUsage: ai.TokenUsage{
			PromptTokens:     100,
			CompletionTokens: 50,
			TotalTokens:      150,
		},
	}
}

// SampleReviewResponse returns a sample code review response.
func SampleReviewResponse() *ai.Response {
	return &ai.Response{
		Text: `# Code Review

## Issues Found
1. Missing input validation
2. TODO comment should be resolved

## Suggestions
- Add proper error handling
- Consider edge cases`,
		TokenUsage: ai.TokenUsage{
			PromptTokens:     200,
			CompletionTokens: 100,
			TotalTokens:      300,
		},
	}
}
