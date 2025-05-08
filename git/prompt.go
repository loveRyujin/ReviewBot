package git

import (
	"bytes"
	"embed"
	"html/template"
)

const (
	CommitMessagePrefix  = "commit_message_prefix"
	CommitMessageTitle   = "commit_message_title"
	CommitMessageSummary = "commit_message_summary"
)

//go:embed template/*
var templateFS embed.FS

func GetCommitMessageTmpl(data map[string]any) (string, error) {
	output, err := templateFS.ReadFile("template/commit_message.tmpl")
	if err != nil {
		return "", err
	}
	tmpl, err := template.New("").Parse(string(output))
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
