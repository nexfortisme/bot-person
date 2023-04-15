# Bot Person
Bot Person is a farily simple Discord bot written in Go. It started off as a fun side project by [AltarCrystal](https://github.com/AltarCrystal) and is now being maintained by Nex. The main feature that comes with Bot Person is integration with the [OpenAI's DaVinci](https://beta.openai.com/docs/models/davinci) and [OpenAI's ChatGPT](https://platform.openai.com/docs/api-reference/chat/create) LLMs. Along with the chat prompting, through the use of [OpenAI's Dall-E API](https://beta.openai.com/docs/guides/images), users can send prompts to the AI and get back image responses back. It allows users to send prompts to the API to see what kind of response it returns.

## Getting Started
If an existing `config.json` doesn't exist, the command line will prompt you for the Keys necessary for the bot to function. The three keys you need are from: OpenAI (for chatbot features), Discord (to have the bot work in general), and FinnHub (for the stocks game). The OpenAI key can be obtained from: https://openai.com/api/, the Discord Key can be obtained from: https://discord.com/developers/applications, and the FinnHub Key can be obtained from: https://finnhub.io/dashboard. From there the bot will start up as normal and the bot will be able to interact with users in any server it's been invited to. As you continue to use the bot, the basic logging will be saved to the `logfile`, the user stats will be saved to the `botTracking.json` and the images will be saved in the `./img` folder where the bot is running. 

## Startup Flags
- `-dev` - The dev flag will start up the bot in "dev" mode where it will use a dev discord bot token instead of the main token. Similar to the main token and the OpenAI key, it will prompt you for one if one doesn't exist in the `config.json`
- `-removeCommands` - This flag will have the bot unregister all of its slash commands on startup, allowing for existing commands to be unregistered from the Discord API. This is useful if you want to change the commands or if you want to remove them all together.
- `-disableLogging` - WIP. This stops the bot from logging all interactions to the `logfile`.
- `-disableTracking` - WIP. This will stop the bot from keeping track of the number of messages on a global and per-user basis.

## Notable Features
- Token System
  - Ability for users to get tokens to spend on Dall-E image requests
  - Daily bonus for Tokens that ranges from between .5 and 5 tokens
  - Ability to transfer tokens to other users
  - LootBox game where a user can spend 2.5 tokens and gets 1 to 250 tokens back on an RNG roll from a randomly generated seed
  - Stock market game. Users are able to purchase stocks for 1 Token = 1 USD and can have their tokens reise and fall based on the activity of the markets.

## Wishlist
- Configurable Commands - Have commands be configurable in the config.json instead of hard coded into the application
  - {command, response, token cost (see below)}
  - A way to specify return elements in the command. ie. {user has <user.messages> sent messages}
- Log File Rotation
  - The file isn't getting super huge at the moment but having the ability to generate a new file when the existing one reaches a certain length.
- Image Folder checking
  - Have a scheduled task that checks to see if the `./img` folder is getting larger than a size set in the `config.json` and trim it as needed.
- Web Manager 
  - A more friendly way of interacting with the bot configurations through a web page
- Tokens
  - Some kind of games to get more tokens
  - Logging of token balance for users to be able to see how their balance changed
- On Initialization of the bot when a new Discord bot token is added, have it log out an invite link for that bot user with the proper permissions
- Tracking of user activity (in chat channels or voice channels) to give out tokens for use for special commands or other future external applications and for use by admins to see activity in their servers.
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
- Lootbox improvements
  - Allow for users to purchase multiple loot boxes at one time
  - Allow for users to open multiple loot boxes at one time
  - Allow for users to store lootboxes for later use, but retaining the RNG seed
- More Configuration
  - Have the daily bonus be part of the config.json instead of being hard coded
  - Have the Admins be configurable in the config.json, also, instead of being hard coded
- Shop Improvements
  - Have users be able to purchase various items that can go up and down in price every x number of minutes/hours
- Better Documentation
  - Have the `/help` command offer better information about what each command is doing
  - More detailed examples to show how commands are used
  - Examples of how data is being stored in the `botTracking.json` file