package prompt

import (
	"bytes"
	"embed"
	"html/template"
)

const (
	FileDiff     = "file_diff"
	SummaryPoint = "summary_point"
)

//go:embed template/*
var templateFS embed.FS

func GetFileDiffSummaryTmpl(placeHolder string, data any) (string, error) {
	output, err := templateFS.ReadFile("template/summarize_file_diff.tmpl")
	if err != nil {
		return "", err
	}
	tmpl, err := template.New("").Parse(string(output))
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, map[string]any{placeHolder: data}); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func GetCommitMessagePrefixTmpl(placeHolder string, data any) (string, error) {
	output, err := templateFS.ReadFile("template/conventional_commit.tmpl")
	if err != nil {
		return "", err
	}
	tmpl, err := template.New("").Parse(string(output))
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, map[string]interface{}{placeHolder: data}); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func GetCommitMessageTitleTmpl(placeHolder string, data any) (string, error) {
	output, err := templateFS.ReadFile("template/summarize_title.tmpl")
	if err != nil {
		return "", err
	}
	tmpl, err := template.New("").Parse(string(output))
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, map[string]interface{}{placeHolder: data}); err != nil {
		return "", err
	}

	return buf.String(), nil
}
