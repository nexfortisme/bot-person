# Repository guidance

## Project overview

Bot Person is a Go 1.24 Discord bot. `main.go` creates the Discord session, registers active slash commands, and dispatches interactions. Command implementations live in `pkg/commands`; message and reply handling lives in `pkg/messages`; external model integrations live in `pkg/external`; SQLite access lives in `pkg/persistance`.

The misspelled package path `persistance` is established throughout the codebase. Do not rename it as part of an unrelated change.

## Making changes

- Keep the command definition in `slashCommands` and its executor in `commandHandlers` in sync. A command is not active merely because a file exists under `pkg/commands`.
- Keep README command and configuration documentation aligned with the commands actually registered in `main.go`.
- Treat Discord's three-second initial interaction deadline as a hard constraint. Acknowledge an interaction before slow network or model work, then edit the response.
- Use the shared HTTP client for local LLM calls. Add timeouts, close response bodies, check HTTP status codes, and handle empty API response arrays before indexing them.
- Keep secrets in `.env`; update `example.env` with names and safe example values only. Never commit `.env`, tokens, generated media, or SQLite databases.
- Preserve user data compatibility. Database schema or economy changes need special care around existing SQLite files and concurrent Discord interactions.
- Use `gofmt` on changed Go files. Avoid reformatting unrelated legacy files in the same change.

## Verification

Run these from the repository root:

```bash
gofmt -w <changed-go-files>
go test ./...
go vet ./...
go build ./...
```

If the environment blocks the default Go build cache, set `GOCACHE` to a writable temporary directory, for example:

```bash
GOCACHE=/tmp/bot-person-go-cache go test ./...
```

Tests must not require live Discord, OpenAI, Perplexity, or local LLM credentials. Isolate parsing and response-handling logic so it can be tested locally.

## Runtime notes

- Normal mode uses `DISCORD_API_KEY`; dev mode (`-dev` or `DEV_MODE=true`) uses `DEV_DISCORD_API_KEY`.
- The default database is `db.sqlite`; `DB_PATH` overrides it. Docker Compose stores it under `./data`.
- `.env` model settings are reloaded every five minutes. Runtime overrides for the local LLM endpoint are deliberately preserved.
- `-removeCommands` removes global slash commands and exits. Use it deliberately against the configured Discord application.
