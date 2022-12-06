# Bot Person
Bot Person is a farily simple Discord bot written in Go. It started off as a fun side project by [AltarCrystal](https://github.com/AltarCrystal) and is now being maintained by Nex. The main feature that comes with Bot Person is integration with the [OpenAI Davinci](https://beta.openai.com/docs/models/davinci) chat bot. Also through the use of [OpenAI's Dall-E API](https://beta.openai.com/docs/guides/images) users can send prompts to the AI and get back image respoonses back. It allows users to send prompts to the API to see what kind of response it returns.

## Getting Started
If an existing `config.json` doesn't exist, the command line will prompt you for an OpenAI key which can be optained from https://openai.com/api/ and for a Discord Bot token which can be gotten from the Discord developer page. From there the bot will start up as normal and the bot will be able to interact with users in any server its been invited to. The bot will also generate `botTracking.json` which keeps a count of how many interactions the bot has has, how many times its been called good, and how many times its been called bad. Along with the general statistics, it also keeps track of those 3 tallies on a per user basis.

## Flags
- `-dev` - The dev flag will start up the bot in "dev" mode where it will use a dev discord bot token instead of the main token. Similar to the main token and the OpenAI key, it will prompt you for one if one doesn't exist in the `config.json`
- `-removeCommands` - This flag will have the bot unregister all of its slash commands on shut down, making the commands unable to be access by users in servers the bot is in.
- `-disableLogging` - WIP. This stops the bot from logging all interactions to the `logfile`.
- `-disableTracking` - WIP. This will stop the bot from keeping track of the number of messages on a global and per-user basis.

## Notable Features
- Token System
  - Ability for users to get tokens to spend on Dall-E image requests
  - Daily bonus for Tokens that ranges from between .5 and 5 tokens
  - Ability to transferr tokens to other users
  - Lootbox game where a user can spend 2.5 tokens and gets 1 to 250 tokens back on an RNG roll from a randomly generated seed

## Wishlist
- Configurable Commands - Have commands be configurable in the config.json instead of hard coded into the application
  - {command, response, token cost (see below)}
  - A way to specify return elements in the command. ie. {user has <user.messages> sent messages}
- Log File Rotation - The file isn't getting super huge at the moment but having the ability to generate a new file when the existing one reaches a certain length.
- Web Manager - A more friendly way of interacting with the bot configurations through a web page
- Tokens
  - Some kind of games to get more tokens
  - Logging of token balance for users to be able to see how their balance changed
- On Initialization of the bot when a new Discord bot token is added, have it log out an invite link for that bot user with the proper permissions
- Tracking of user activity (in chat channels or voice channels) to give out tokens for use for special commands or other future external applications
- Save Images from Dall-E to disk before serving
  - Images from the Dall-E CDN endpoint have a shockingly short life, by serving them up from the disk, the images will then be handled by Discord's CDN and have an exponentially longer life
- Search Game
  - Every x number of minutes/hours the user can /search and see what they can get
    - They can find tokens, lootboxes, or something else (placeholder for when I actually have ideas of what that something else could be)
- RNG Breakpoints
  - Allow for RNG to be broken up into ranges and from there have it pull a result from a pre-defined loot table
- Better Logging of transactions
  - Allow for users to see where they have spent their tokens and where they have gotten them from
- User Inventory
  - Allow for users to purchase multiple loot boxes or other items (when I figure out what they are) and store them for later
  - Allow for trading of items between players
