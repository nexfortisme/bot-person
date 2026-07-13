# Bot Person

[![Discord](https://img.shields.io/badge/Discord-7289DA?style=for-the-badge&logo=discord&logoColor=white)](https://discord.gg/MtEG5zMtUR)

Bot Person is a Discord bot originally developed by [AltarCrystal](https://github.com/AltarCrystal) and continued by [nexfortisme](https://github.com/nexfortisme). It combines local and hosted AI chat, image generation and editing, web search, conversation replies, usage statistics, and a token economy.

## Features

- Local LLM chat through an OpenAI-compatible endpoint such as Ollama
- OpenAI chat, including optional web search for `/bot-gpt`
- OpenAI image generation and, with a GPT image model, image editing
- Threaded follow-up conversations by replying to bot responses
- Image understanding when supported by the configured local model
- Perplexity web search
- SQLite-backed balances, daily bonuses, streaks, settings, statistics, and conversation history
- Five-minute hot reloads for model and local-endpoint settings in `.env`

Video generation, loot boxes, voice commands, and several experimental commands have source files in the repository but are not currently registered.

## Active slash commands

| Command | Purpose |
| --- | --- |
| `/bot` | Chat with the configured local LLM |
| `/bot-gpt` | Chat with OpenAI, optionally with web search |
| `/image` | Generate an image for one token |
| `/search` | Search the web through Perplexity |
| `/balance` | Check your own or another user's balance |
| `/send` | Send tokens to another user |
| `/bonus` | Claim a daily bonus and manage a missed streak |
| `/burn` | Permanently remove tokens from your balance |
| `/my-stats` | Show personal usage and economy statistics |
| `/bot-stats` | Show bot-wide statistics |
| `/set` | Configure supported per-user settings |
| `/about`, `/help` | Show bot and command information |
| `/invite`, `/donations`, `/broken` | Show community and support information |

You can also mention the bot to chat, attach images for supported multimodal models, or reply to a `/bot`, `/bot-gpt`, or generated-image response to continue working with it.

## Requirements

- Go 1.24 or Docker
- A Discord application and bot token
- OpenAI and Perplexity API credentials
- An OpenAI-compatible local chat-completions server for `/bot` and retorts

## Configuration

Copy the example file and fill in its blank values:

```bash
cp example.env .env
```

| Variable | Description |
| --- | --- |
| `DISCORD_API_KEY` | Discord bot token used in normal mode |
| `DEV_DISCORD_API_KEY` | Discord bot token used with `-dev` or `DEV_MODE=true` |
| `OPEN_AI_API_KEY` | OpenAI API key for hosted chat and images |
| `PERPLEXITY_API_KEY` | Perplexity API key for `/search` |
| `BOT_OPEN_AI_MODEL` | OpenAI model used by legacy hosted bot paths |
| `OPEN_AI_MODEL` | OpenAI model used by `/bot-gpt` |
| `IMAGE_GENERATION_MODEL` | OpenAI image model used by `/image` and image edits |
| `LOCAL_LLM_MODEL` | Local model used by `/bot`; defaults to `gemma4:12b` |
| `RETORT_MACHINE_MODEL` | Local model used for short retorts and image-intent detection |
| `LOCAL_LLM_CHAT_COMPLETIONS_ENDPOINT` | OpenAI-compatible local endpoint; defaults to `http://localhost:11434/v1/chat/completions` |
| `ADMINS` | Comma-separated Discord user IDs allowed to use admin message commands |
| `DB_PATH` | Optional SQLite path; defaults to `db.sqlite` |
| `DEV_MODE` | Set to `true` to use the dev Discord token without the CLI flag |

The application currently validates all OpenAI model settings and both hosted API keys at startup because active commands depend on them. It validates only the Discord token selected for the current mode.

## Run locally

```bash
go run .
```

Useful flags:

```bash
# Use DEV_DISCORD_API_KEY
go run . -dev

# Read configuration only from the process environment
go run . -useEnvFile=false

# Remove this application's registered global slash commands, then exit
go run . -removeCommands
```

The bot registers global slash commands, so Discord may take time to display registration changes.

## Run with Docker

On macOS or Windows with Docker Desktop:

```bash
docker compose up -d --build
```

On native Linux, the override uses host networking so the container can reach a loopback-only local LLM:

```bash
docker compose -f docker-compose.yml -f docker-compose.linux.yml up -d --build
```

Compose mounts `.env` read-only, persists SQLite data under `./data`, and stores generated media under `./img`.

## Development

```bash
gofmt -w <changed-go-files>
go test ./...
go vet ./...
go build ./...
```

There are currently only focused unit tests; most Discord and provider integrations still require manual or mocked integration coverage. See [AGENTS.md](AGENTS.md) for repository-specific contribution guidance.

For support, join the [Discord server](https://discord.gg/MtEG5zMtUR) or open a [GitHub issue](https://github.com/nexfortisme/bot-person/issues).

## License

Bot Person is available under the [MIT License](LICENSE).
