package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/loveRyujin/ReviewBot/pkg/config"
	"github.com/stretchr/testify/require"
)

// Run each invocation in a fresh process, just like the installed CLI.
func TestCLIProcess(t *testing.T) {
	i := slices.Index(os.Args, "--")
	if i == -1 {
		return
	}
	rootCmd.SetArgs(os.Args[i+1:])
	Execute()
}

func runCLI(t *testing.T, dir string, input string, args ...string) string {
	t.Helper()
	ctx, cancel := context.WithTimeout(t.Context(), 20*time.Second)
	defer cancel()
	process := exec.CommandContext(ctx, os.Args[0], append([]string{"-test.run=^TestCLIProcess$", "--"}, args...)...)
	process.Dir = dir
	process.Stdin = strings.NewReader(input)
	output, err := process.CombinedOutput()
	require.NoError(t, err, "%s", output)
	return string(output)
}

func TestCLIConfigRoundTrip(t *testing.T) {
	dir := t.TempDir()
	setConfigEnv(t, dir)
	path := filepath.Join(dir, "selected.yaml")
	const other = "git:\n  exclude_list: [vendor, dist]\ncustom:\n  Keep.Me: {}\n"
	const initial = "AI:\n  Model: file-model\n  api_key: file-secret\n"
	require.NoError(t, os.WriteFile(path, []byte(initial+other), 0o600))
	const fallback = "ai:\n  model: fallback-model\n"
	defaultPath := filepath.Join(dir, "reviewbot.yaml")
	require.NoError(t, os.WriteFile(defaultPath, []byte(fallback), 0o600))
	t.Setenv("REVIEWBOT_AI_API_KEY", "")
	t.Setenv("REVIEWBOT_AI_MODEL", "")
	list := runCLI(t, dir, "", "--config", path, "config", "list")
	require.Contains(t, list, "file-model")
	require.Regexp(t, `(?m)^ai\.api_key\s+\*+\s`, list)
	require.NotContains(t, list, "file-secret")
	t.Setenv("REVIEWBOT_AI_API_KEY", "environment-secret")
	t.Setenv("REVIEWBOT_AI_MODEL", "environment-model")
	list = runCLI(t, dir, "", "--config", path, "config", "list")
	require.Contains(t, list, "environment-model")
	require.NotContains(t, list, "file-model")
	require.NotContains(t, list, "file-secret")
	require.NotContains(t, list, "environment-secret")
	runCLI(t, dir, "", "--config", path, "config", "set", "ai.model", "new-model")
	saved, err := os.ReadFile(path)
	require.NoError(t, err)
	require.YAMLEq(t, strings.Replace(initial, "file-model", "new-model", 1)+other, string(saved))
	savedDefault, err := os.ReadFile(defaultPath)
	require.NoError(t, err)
	require.Equal(t, fallback, string(savedDefault))
	t.Setenv("REVIEWBOT_AI_MODEL", "")
	require.Contains(t, runCLI(t, dir, "", "--config", path, "config", "list"), "new-model")
}

func TestCLIConfigSetFirstAPIKey(t *testing.T) {
	dir := t.TempDir()
	setConfigEnv(t, dir)
	t.Setenv("REVIEWBOT_AI_API_KEY", "")
	path := filepath.Join(dir, "settings", "first.yaml")
	output := runCLI(t, dir, "", "--config", path, "config", "set", "ai.api_key", "first-fake-key")
	require.NotContains(t, output, "first-fake-key")
	saved, err := os.ReadFile(path)
	require.NoError(t, err)
	require.YAMLEq(t, "ai:\n  api_key: first-fake-key\n", string(saved))
	if runtime.GOOS != "windows" {
		info, err := os.Stat(path)
		require.NoError(t, err)
		require.Equal(t, os.FileMode(0o600), info.Mode().Perm())
	}
	list := runCLI(t, dir, "", "--config", path, "config", "list")
	require.Regexp(t, `(?m)^ai\.api_key\s+\*+\s`, list)
	require.NotContains(t, list, "first-fake-key")
}

func TestCLIConfigSetHeaders(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "selected.yaml")
	require.NoError(t, os.WriteFile(path, []byte("ai:\n  api_key: fake-key\n"), 0o600))
	runCLI(t, dir, "", "--config", path, "config", "set", "proxy.headers", "Authorization=Bearer a,b, X-Trace=1")
	runCLI(t, dir, "", "--config", path, "config", "set", "git.exclude_list", "vendor,dist")

	cfg, err := config.Load(config.LoadOptions{ExplicitPath: path, EnvPrefix: "REVIEWBOT_SET_TEST"})
	require.NoError(t, err)
	require.Equal(t, []string{"Authorization=Bearer a,b", "X-Trace=1"}, cfg.Proxy.Headers)
	require.Equal(t, []string{"vendor", "dist"}, cfg.Git.ExcludedList)
}

func TestCLIRuntimeSettings(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Stream   bool
			Messages []struct{ Content string }
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Error(err)
			return
		}
		if req.Stream {
			if !strings.Contains(req.Messages[len(req.Messages)-1].Content, "fixture diff") {
				t.Error("configured diff file missing from review prompt")
			}
			w.Header().Set("Content-Type", "text/event-stream")
			fmt.Fprint(w, "data: {\"choices\":[{\"delta\":{\"content\":\"stream 90% done\"}}]}\n\ndata: {\"choices\":[],\"usage\":{\"total_tokens\":3}}\n\ndata: [DONE]\n\n")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"choices":[{"message":{"role":"assistant","content":"fixture response"}}]}`)
	}))
	defer server.Close()
	dir := t.TempDir()
	t.Setenv("REVIEWBOT_AI_BASE_URL", server.URL)
	t.Setenv("REVIEWBOT_AI_API_KEY", "test-key")
	path := filepath.Join(dir, "reviewbot.yaml")
	require.NoError(t, os.WriteFile(path, []byte("runtime:\n  review:\n    mode: external\n    diff_file: patch.txt\n    stream: true\n  commit:\n    preview: true\n"), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "patch.txt"), []byte("fixture diff"), 0o600))
	require.Contains(t, runCLI(t, dir, "", "-c", path, "review"), "stream 90% done")
	require.Contains(t, runCLI(t, dir, "", "-c", path, "review", "--stream=false"), "fixture response")

	git := func(args ...string) string {
		t.Helper()
		command := exec.Command("git", args...)
		command.Dir = dir
		output, err := command.CombinedOutput()
		require.NoError(t, err, "%s", output)
		return string(bytes.TrimSpace(output))
	}
	git("init")
	git("config", "user.name", "Tester")
	git("config", "user.email", "test@example.org")
	git("config", "commit.gpgsign", "false")
	git("commit", "--allow-empty", "-m", "initial")
	before := git("rev-parse", "HEAD")
	git("add", "patch.txt")
	runCLI(t, dir, "n\r", "-c", path, "commit")
	require.Equal(t, before, git("rev-parse", "HEAD"))
	runCLI(t, dir, "", "-c", path, "commit", "--preview=false")
	require.NotEqual(t, before, git("rev-parse", "HEAD"))
}
