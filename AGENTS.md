# Repository Guidelines

## Project Structure & Module Organization
ReviewBot is a Go module built around a single CLI binary. Source code lives in `cmd/reviewbot` for the executable and `pkg` for reusable packages (`pkg/command`, `pkg/form`, `pkg/progress`, `pkg/version`). Provider integrations and prompt assets reside in `ai/`, `llm/`, and `prompt/`. Configuration examples are under `config/`, while release artifacts land in `release/<platform>/`. Tests accompany their packages, e.g. `pkg/command/command_test.go` and `prompt/translation_test.go`. UI assets and demos are stored in `images/`. CI workflows are defined in `.github/workflows/`.

## Build, Test, and Development Commands
- `make build`: runs `go mod tidy`, compiles the binary, and drops it in `bin/reviewbot`.
- `make build_<target>`: cross-compiles into the matching `release/<os>/<arch>/` folder (see Makefile).
- `go test ./...`: executes the entire test suite; use `-run` to scope changes.
- `go run ./cmd/reviewbot --help`: quick local smoke test of CLI options.

## Coding Style & Naming Conventions
Follow idiomatic Go style. Format all Go files with `gofmt` or `go fmt ./...` before committing. Exported types and functions use PascalCase; keep unexported helpers in lowerCamelCase. Prefer table-driven tests and keep files under 400 lines for readability. When adding prompts or configs, mirror existing naming (e.g., `prompt/template/<feature>.tmpl`, `config/reviewbot.yaml`).

## Testing Guidelines
Place tests alongside implementation using the `TestXxx` naming scheme. For features touching external services, stub providers via the abstractions in `pkg/command`. Maintain existing coverage by exercising new branches with table-driven cases. Run `go test -cover ./...` before opening a pull request, and document any gaps.

## Commit & Pull Request Guidelines
Commit messages follow Conventional Commits (`feat:`, `fix:`, `docs:`) as shown in the recent Git history. Group related edits and avoid mixing refactors with feature work. Pull requests should include: a concise summary, reproduction or verification steps (commands, screenshots, or GIFs when UI output changes), updated docs for configuration changes, and references to related issues or discussions.
