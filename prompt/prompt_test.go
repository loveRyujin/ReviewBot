package prompt

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadTemplateFromCustomDir(t *testing.T) {
	dir := t.TempDir()
	custom := filepath.Join(dir, "code_review_file_diff.tmpl")
	if err := os.WriteFile(custom, []byte("custom"), 0o644); err != nil {
		t.Fatalf("write custom template: %v", err)
	}

	SetTemplateDir(dir)
	t.Cleanup(func() { SetTemplateDir("") })

	content, err := GetPromptTmpl(CodeReviewFileDiffTmpl, map[string]any{})
	if err != nil {
		t.Fatalf("GetPromptTmpl returned error: %v", err)
	}

	if content != "custom" {
		t.Fatalf("expected custom content, got %q", content)
	}
}

func TestCustomTemplateFallback(t *testing.T) {
	dir := t.TempDir()
	SetTemplateDir(dir)
	t.Cleanup(func() { SetTemplateDir("") })

	content, err := GetPromptTmpl(CodeReviewFileDiffTmpl, map[string]any{})
	if err != nil {
		t.Fatalf("GetPromptTmpl returned error: %v", err)
	}

	if len(content) == 0 {
		t.Fatal("expected non-empty default template content")
	}
}

func TestCustomTemplatePathTraversal(t *testing.T) {
	dir := t.TempDir()
	SetTemplateDir(dir)
	t.Cleanup(func() { SetTemplateDir("") })

	_, err := GetPromptTmpl("../prompt.go", map[string]any{})
	if err == nil {
		t.Fatal("expected error for path traversal")
	}
}
