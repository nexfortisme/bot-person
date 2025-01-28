package logging

type EventType int

const (
	// System Events
	BOT_START EventType = iota
	BOT_STOP

	// General Commands
	COMMAND_ABOUT
	COMMAND_BROKEN
	COMMAND_DONATIONS
	COMMAND_HELP

	// Economy Commands
	COMMAND_BALANCE
	COMMAND_BONUS
	COMMAND_BURN
	COMMAND_LOOTBOX
	COMMAND_SEND
	COMMAND_INVITE

	// Fun Commands
	COMMAND_BOT
	COMMAND_BOT_GPT
	EXTERNAL_GPT_REQUEST
	EXTERNAL_GPT_RESPONSE

	// Image Commands
	COMMAND_IMAGE
	EXTERNAL_DALLE_REQUEST
	EXTERNAL_DALLE_RESPONSE

	// Stat Commands
	COMMAND_BOT_STATS
	COMMAND_MY_STATS

	// Economy Events
	ECONOMY_ADD_TOKENS
	ECONOMY_REMOVE_TOKENS
	ECONOMY_BURN_TOKENS
	ECONOMY_CREATE_TOKENS

	// Stock Commands
	COMMAND_PORTFOLIO
	COMMAND_STOCKS

	// External Stock Events
	EXTERNAL_STOCK_REQUEST
	EXTERNAL_STOCK_RESPONSE

	// Economy Stock Events
	ECONOMY_ADD_STOCK
	ECONOMY_REMOVE_STOCK
	ECONOMY_SELL_STOCK

	// User Events
	USER_MESSAGE
	USER_BAD_BOT
	USER_GOOD_BOT

	// Error
	EXTERNAL_API_ERROR
	INTERNAL_ERROR
	DATABASE_ERROR

	// Lenny
	LENNY

	// Logging
	USER_SET_BONUS_STREAK

	// Test
	TEST_EVENT

	HSR_CODE

	TTS_JOIN
	TTS_LEAVE
	TTS_MESSAGE
	
	ERROR

	COMMAND_SEARCH
)

func (e EventType) ToInt() int {
	switch e {
	case BOT_START:
		return 1
	case BOT_STOP:
		return 2
	case COMMAND_BOT:
		return 3
	case COMMAND_BOT_GPT:
		return 4
	case COMMAND_BOT_STATS:
		return 5
	case COMMAND_MY_STATS:
		return 6
	case ECONOMY_ADD_TOKENS:
		return 7
	case ECONOMY_REMOVE_TOKENS:
		return 8
	case ECONOMY_BURN_TOKENS:
		return 9
	case ECONOMY_CREATE_TOKENS:
		return 10
	case COMMAND_PORTFOLIO:
		return 11
	case COMMAND_STOCKS:
		return 12
	case EXTERNAL_STOCK_REQUEST:
		return 13
	case EXTERNAL_STOCK_RESPONSE:
		return 14
	case ECONOMY_ADD_STOCK:
		return 15
	case ECONOMY_REMOVE_STOCK:
		return 16
	case ECONOMY_SELL_STOCK:
		return 17
	case USER_MESSAGE:
		return 18
	case USER_BAD_BOT:
		return 19
	case USER_GOOD_BOT:
		return 20
	case EXTERNAL_API_ERROR:
		return 21
	case INTERNAL_ERROR:
		return 22
	case DATABASE_ERROR:
		return 23
	case LENNY:
		return 24
	case USER_SET_BONUS_STREAK:
		return 25
	case TEST_EVENT:
		return 26
	case HSR_CODE:
		return 27
	case TTS_JOIN:
		return 28
	case TTS_LEAVE:
		return 29
	case TTS_MESSAGE:
		return 30
	case ERROR:
		return 31
	case COMMAND_SEARCH:
		return 32
	default:
		return 0
	}
}
