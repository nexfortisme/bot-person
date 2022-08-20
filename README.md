# Bot Person
Bot Person is a farily simple Discord bot written in Go. It started off as a fun side project by AltarCrystal (citation needed) and is now being maintained by Nex. The main feature that comes with Bot Person is integration with the OpenAI davinci chat bot. It allows users to send prompts to the API to see what kind of response it returns

## Getting Started
If an existing `config.json` doesn't exist, the command line will prompt you for an OpenAI key which can be optained from https://openai.com/api/ and for a Discord bot token which can be gotten from the discord developer page. From there the bot will start up as normal and the bot will be able to interact with users in any server its been invited to. The bot will also generate `botTracking.json` which keeps a count of how many interactions the bot has has, how many times its been called good, and how many times its been called bad. Along with the general statistics, it also keeps track of those 3 tallies on a per user basis.

## Flags
- `-dev` - The dev flag will start up the bot in "dev" mode where it will use a dev discord bot token instead of the main token. Similar to the main token and the OpenAI key, it will prompt you for one if one doesn't exist in the `config.json`
- `-removeCommands` - This flag will have the bot unregister all of its slash commands on shut down, making the commands unable to be access by users in servers the bot is in.
- `-disableLogging` - WIP. This stops the bot from logging all interactions to the `logfile`.
- `-disableTracking` - WIP. This will stop the bot from keeping track of the number of messages on a global and per-user basis.

## Wishlist
- Configurable Commands - Have commands be configurable in the config.json instead of hard coded into the application
- Log File Rotation - The file isn't getting super huge at the moment but having the ability to generate a new file when the existing one reaches a certain length.
- Web Manager - A more friendly way of interacting with the bot configurations through a web page
