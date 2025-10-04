# Repository Guidelines

## Project Structure & Module Organization
ReviewBot is a Go module built around the `reviewbot` CLI. Entry commands live in `cmd/`, reusable helpers in `pkg/`, and provider integrations under `ai/` and `llm/`. Prompt templates and translation helpers are in `prompt/`. Configuration samples reside in `config/`, release builds in `release/<platform>/`, and demo assets in `images/`. Tests sit beside their packages, for example `pkg/command/command_test.go` and `prompt/translation_test.go`. CI configs live in `.github/workflows/`.

## Build, Test, and Development Commands
- `make build`: runs `go mod tidy`, compiles the binary, and emits `bin/reviewbot`.
- `make install`: builds and installs the CLI to `GOBIN` (or `GOPATH/bin`) using the same ldflags as `make build`.
- `make build_<target>`: cross-compiles into `release/<os>/<arch>/` (see Makefile targets).
- `go test ./...`: executes the full test suite; add `-run <pattern>` for narrower checks.
- `go run ./cmd/reviewbot --help`: fast manual verification of CLI wiring.

## Coding Style & Naming Conventions
Follow idiomatic Go style. Format Go files with `gofmt` or `go fmt ./...`. Exported types and functions use PascalCase; keep unexported helpers camelCase. Prefer table-driven tests and keep files under 400 lines for readability. When adding prompts or configs, mirror the existing naming (e.g., `prompt/template/<feature>.tmpl`, `config/reviewbot.yaml`).

## Testing Guidelines
Place tests alongside implementation using the `TestXxx` naming scheme. For features touching external services, stub providers via the adapters in `pkg/command`. Maintain coverage by exercising new branches with table-driven cases. Run `go test -cover ./...` before opening a pull request, and record any unavoidable gaps.

## Commit & Pull Request Guidelines
Commit messages follow Conventional Commits (`feat:`, `fix:`, `docs:`). Group related edits and avoid mixing refactors with feature work. Pull requests should include: a concise summary, verification steps (commands, screenshots, or GIFs when output changes), updated docs for configuration changes, and links to related issues or discussions. Document translation behaviors or language additions when touching `--output_lang`.
