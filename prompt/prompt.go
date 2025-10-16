package prompt

import (
	"bytes"
	"embed"
	"html/template"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

const (
	// Prompt template file names
	CodeReviewFileDiffTmpl  = "code_review_file_diff.tmpl"
	CommitMessagePrefixTmpl = "conventional_commit.tmpl"
	CommitMessageTitleTmpl  = "summarize_title.tmpl"
	CommitFileDiffTmpl      = "summarize_file_diff.tmpl"
	TranslationTmpl         = "translation.tmpl"

	// PlaceHolders
	FileDiff      = "file_diffs"
	SummaryPoint  = "summary_points"
	OutputLang    = "output_language"
	OutputMessage = "output_message"
)

//go:embed template/*
var templateFS embed.FS

var (
	customTemplateDir string
	customDirLock     sync.RWMutex
)

// SetTemplateDir configures a filesystem path for overriding embedded prompts.
func SetTemplateDir(dir string) {
	customDirLock.Lock()
	defer customDirLock.Unlock()

	customTemplateDir = strings.TrimSpace(dir)
}

// GetPromptTmpl reads a template file and executes it with the provided data.
func GetPromptTmpl(file string, data map[string]any) (string, error) {
	buf, err := processTmpl(file, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func processTmpl(file string, data map[string]any) (*bytes.Buffer, error) {
	output, err := loadTemplate(file)
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

func loadTemplate(file string) ([]byte, error) {
	customDirLock.RLock()
	dir := customTemplateDir
	customDirLock.RUnlock()

	if dir != "" {
		content, err := readCustomTemplate(dir, file)
		if err == nil {
			return content, nil
		}
		if !os.IsNotExist(err) {
			return nil, err
		}
	}

	return templateFS.ReadFile(path.Join("template", file))
}

func readCustomTemplate(dir, file string) ([]byte, error) {
	cleanedDir := filepath.Clean(dir)
	if cleanedDir == "" {
		return nil, os.ErrNotExist
	}

	tmplPath := filepath.Join(cleanedDir, file)
	tmplPath = filepath.Clean(tmplPath)
	if !strings.HasPrefix(tmplPath, cleanedDir+string(os.PathSeparator)) && tmplPath != cleanedDir {
		return nil, os.ErrPermission
	}

	return os.ReadFile(tmplPath)
}
