[![Discord](https://img.shields.io/badge/Discord-7289DA?style=for-the-badge&logo=discord&logoColor=white)](https://discord.gg/MtEG5zMtUR)
## What is Bot Person
Bot Person is a discord bot originally developed by [AltarCrystal](github.com/AltarCrystal) and had continued development done on it by [myself](github.com/nexfortisme). The bot incorporates many features from [OpenAI](https://openai.com/), such as ChatGPT GPT-4, GPT-3.5 and Dall-E. 

Users are able to use the bot to interact with those APIs and get responses put into the discord chat they were asked in. To limit users and prevent them from racking up an absurd API bill, there is a token system built into the bot that users have to spend to make certain calls (mostly just Dall-E at the moment). To store and track the user data (only necessary for the use of the bot), the bot uses [SurrealDB](https://surrealdb.com/) as the main database.

## Getting Started

To get started, edit the `example.env` and fill it in with the required fields and rename it to just `.env`. You will also need an [SurrealDB](https://surrealdb.com/) instance running and have a provisioned user for the bot to use. 

### Quick Start
To be able to get an executable to run the bot, run the following command:
```
go build cmd/bot-person/main.go
```
That will generate a `main` executable that will be able to be run on whatever system was used to build it. Make sure that the proper `.env` variables are set and you're good to go.
### SurrealDB
The way that the main Bot Person instance is using the database is through Docker. The way I managed to have success running it was with the following commands, replace the items in <> with the indicated field:
```
mkdir mydata

docker run --pull always -p <port>:8000 -v $(pwd)/mydata:/mydata surrealdb/surrealdb:latest start --auth --user <root-username> --pass <root-password> file:/mydata/mydatabase.db
```

That will create the database instance, and the folder where the data will be stored. The user that will get provisioned will have the user and password of `<root-username>` and `<root-password>`. They are what you will use in [Surrealist](https://surrealdb.com/docs/surrealist) to connect to and manage the database with

### Startup Flags
The bot comes with some flags that can be useful when starting and running the bot.
- `--dev` - This will start up the bot in "dev" mode. It will have all the features of the main bot, but will use the dev discord token to run the bot instead of the main token.
- `--removeCommands` - This will have the bot find all the commands registered with discord and unregister them before re-registering them on startup.

## Notable Features
- Token System
	- Ability for users to get tokens to spend on Dall-E image requests
	- Daily bonus for Tokens that ranges from between .5 and 5 tokens
	- Ability to transfer tokens to other users
	- LootBox game where a user can spend 5 tokens and gets 3.63 to 50 tokens back on an RNG roll from a randomly generated seed