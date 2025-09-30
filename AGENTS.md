# Repository Guidelines

## Project Structure & Module Organization
Backend logic lives in `main.go` and `app.go`, orchestrated by Wails via `wails.json`. Keep Go helpers alongside `app.go`; if you add packages, place them under `internal/`. The Vue UI stays in `frontend/src`, with reusable widgets in `frontend/src/components`, static media in `frontend/src/assets`, and experimental flows in `frontend/src/volt`. Generated Wails bindings live in `frontend/wailsjs`--do not edit them manually. Compiled assets land in `frontend/dist`; packaged binaries live under `build/`.

## Build, Test, and Development Commands
- `wails dev`: run from the repo root to spin up the live dev server.
- `wails build`: builds the Go backend and Vue frontend into `build/bin`.
- `cd frontend && npm install`: install UI dependencies (pnpm also works).
- `cd frontend && npm run dev`: start the Vite-only dev server for frontend iteration.
- `go test ./...`: placeholder until backend tests exist; run before Go changes.

## Coding Style & Naming Conventions
Format Go code with `gofmt` (tabs for indentation) and annotate exported symbols with GoDoc comments. Use descriptive receivers, avoid global mutable state, and surface config through `App` methods. In Vue, keep `<script setup>` components in PascalCase files, camelCase composables, and group Tailwind utilities logically. Shared TypeScript types stay near usage in `frontend/src/volt` or components; run `npx vue-tsc --noEmit` pre-commit.

## Testing Guidelines
Add Go tests alongside code using `_test.go` and table-driven cases. Frontend tests are not wired up yet--use Vitest with `*.spec.ts` next to components when you add coverage. Document manual verification (e.g., League client mocks) in PRs until automation matures.

## Commit & Pull Request Guidelines
Follow Conventional Commits (`feat:`, `fix:`, `chore:`) so tooling can infer changesets. Keep commits focused and note critical build/test output when relevant. Pull requests need a succinct summary, linked issues, UI screenshots, and confirmation that `wails build` succeeds. Request a review whenever you touch authentication or network calls.

## Security & Configuration Tips
Do not ship real Riot tokens in `app.go`; load secrets from environment variables or the Wails secure store and keep them out of version control. Keep the TLS bypass limited to local dev and log firewall or proxy requirements in `README.md`.
