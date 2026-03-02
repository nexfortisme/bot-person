package persistance

import (
	"time"
)

type BPInteraction int

const (
	BPChatInteraction       BPInteraction = iota // anything with the divinci chat model
	BPImageRequest                               // any image request
	BPBadBotInteraction                          // bad bot
	BPGoodBotInteraction                         // good bot
	BPBasicInteraction                           // any basic bot interaction
	BPSystemInteraction                          // TODO - Figure out what this means
	BPCreateUserInteraction                      // used by system for creating users
	BPLennyFaceInteracton
)

func GetUserStats(userId string) MyStats {

	user, err := GetUser(userId)
	if err != nil {
		return MyStats{}
	}

	var interactionCountData int64
	var goodBotCountData int64
	var badBotCountData int64
	var lootBoxCountData int64
	var imageCountData int64
	var chatCountData int64

	// Interaction Count
	err = RunQuery("SELECT count(*) AS count FROM events WHERE EventUser = ? AND EventType NOT IN (12, 13)", &interactionCountData, userId)
	if err != nil {
		return MyStats{}
	}

	// Good Bot Count
	err = RunQuery("SELECT count(*) AS count FROM events WHERE EventUser = ? AND EventType IN (34)", &goodBotCountData, userId)
	if err != nil {
		return MyStats{}
	}

	// Bad Bot Count
	err = RunQuery("SELECT count(*) AS count FROM events WHERE EventUser = ? AND EventType IN (33)", &badBotCountData, userId)
	if err != nil {
		return MyStats{}
	}

	// Loot Box Count
	err = RunQuery("SELECT count(*) AS count FROM events WHERE EventUser = ? AND EventType IN (9)", &lootBoxCountData, userId)
	if err != nil {
		return MyStats{}
	}

	// Image Count
	err = RunQuery("SELECT count(*) AS count FROM events WHERE EventUser = ? AND EventType IN (16)", &imageCountData, userId)
	if err != nil {
		return MyStats{}
	}

	// Chat Count
	err = RunQuery("SELECT count(*) AS count FROM events WHERE EventUser = ? AND EventType IN (12, 13)", &chatCountData, userId)
	if err != nil {
		return MyStats{}
	}

	var myStats MyStats

	myStats.InteractionCount = int(interactionCountData)
	myStats.GoodBotCount = int(goodBotCountData)
	myStats.BadBotCount = int(badBotCountData)
	myStats.LootBoxCount = int(lootBoxCountData)
	myStats.ImageCount = int(imageCountData)
	myStats.ChatCount = int(chatCountData)

	myStats.ImageTokens = user.ImageTokens
	myStats.BonusStreak = user.BonusStreak

	// Convert the string timestamp to a time.Time object
	var lastBonus time.Time
	if user.LastBonus != "" {
		lastBonus, err = time.Parse(time.RFC3339, user.LastBonus)
		if err != nil {
			lastBonus = time.Time{}
		}
	}

	// Set the LastBonus field to the parsed time
	myStats.LastBonus = lastBonus

	return myStats
}
