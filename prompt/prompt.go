package prompt

import (
	"bytes"
	"embed"
	"html/template"
	"path/filepath"
)

const (
	// Prompt template file names
	CodeReviewFileDiffTmpl  = "code_review_file_diff.tmpl"
	CommitMessagePrefixTmpl = "conventional_commit.tmpl"
	CommitMessageTitleTmpl  = "summarize_title.tmpl"
	CommitFileDiffTmpl      = "summarize_file_diff.tmpl"

	// PlaceHolders
	FileDiff     = "file_diffs"
	SummaryPoint = "summary_points"
)

//go:embed template/*
var templateFS embed.FS

// GetPromptTmpl reads a template file and executes it with the provided data.
func GetPromptTmpl(file string, data map[string]any) (string, error) {
	buf, err := processTmpl(file, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func processTmpl(file string, data map[string]any) (*bytes.Buffer, error) {
	output, err := templateFS.ReadFile(filepath.Join("template", file))
	if err != nil {
		return nil, err
	}
	tmpl, err := template.New("").Parse(string(output))
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, err
	}

	return &buf, nil
}
