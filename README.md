[![Discord](https://img.shields.io/badge/Discord-7289DA?style=for-the-badge&logo=discord&logoColor=white)](https://discord.gg/MtEG5zMtUR)

## What is Bot Person

Bot Person is a feature-rich Discord bot originally developed by [AltarCrystal](https://github.com/AltarCrystal) and continued by [nexfortisme](https://github.com/nexfortisme). The bot integrates AI models for chat, image generation, and video generation alongside an economy system built around tokens.

## Features

- **AI Chat**: Chat via a local LLM or OpenAI's GPT models
- **Image Generation**: Create AI-generated images using DALL-E
- **Video Generation**: Generate short AI videos using OpenAI's Sora
- **Token Economy**: Earn and spend tokens through daily bonuses and commands
- **Statistics Tracking**: Monitor your usage and daily bonus streaks
- **Search**: Search the web via Perplexity directly from Discord

## Commands

### AI & Generation
- `/bot` - Chat with a local LLM (sarcastic, humorous responses)
- `/bot-gpt` - Chat with OpenAI's GPT model
- `/image` - Generate images using DALL-E
- `/slop` - Generate a short AI video using Sora (4, 8, or 12 seconds); costs tokens

### Economy & Rewards
- `/balance` - Check your current token balance
- `/send` - Send tokens to another user
- `/bonus` - Claim your daily token bonus; includes a streak system with Save/Reset streak buttons
- `/burn` - Burn tokens

### Statistics & Information
- `/my-stats` - View your personal statistics
- `/bot-stats` - View bot-wide statistics
- `/about` - Learn about the bot
- `/help` - Get help with commands
- `/search` - Search for information via Perplexity

### Utility
- `/invite` - Get the bot's invite link
- `/donations` - See who has contributed to keeping the bot running
- `/broken` - Get information on how to report a bug or issue
- `/set` - Configure bot settings

## Setup

Copy `example.env` to `.env` and fill in the required values:

```
OPEN_AI_API_KEY=        # OpenAI API key (for /bot-gpt, /image, /slop)
DISCORD_API_KEY=        # Discord bot token
DEV_DISCORD_API_KEY=    # Discord bot token for dev mode
PERPLEXITY_API_KEY=     # Perplexity API key (for /search)

BOT_OPEN_AI_MODEL=      # Model used by /bot-gpt
OPEN_AI_MODEL=          # Model used for chat
IMAGE_GENERATION_MODEL= # Model used for /image

ADMINS=                 # Comma-separated list of admin user IDs
```

The `/bot` command requires a locally-running LLM server (OpenAI-compatible API).

## Running

```bash
# Normal mode
go run main.go

# Dev mode (uses DEV_DISCORD_API_KEY)
go run main.go -dev

# Remove registered slash commands on startup
go run main.go -removeCommands
```

Or use the provided `Dockerfile` to run the bot in a container.

## Support & Development

This bot is actively maintained and improved. For support, join our [Discord server](https://discord.gg/MtEG5zMtUR) or open an issue on [GitHub](https://github.com/nexfortisme/bot-person/issues).

## Contributing

Contributions are welcome! Feel free to submit issues and pull requests to help improve the bot.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
