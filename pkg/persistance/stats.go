package persistance

import (
	persistance "main/pkg/persistance/models"
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

func GetUserStats(userId string) persistance.MyStats {

	user, err := GetUser(userId)
	if err != nil {
		panic(err)
	}

	var interactionCountData int64
	var goodBotCountData int64
	var badBotCountData int64
	var lootBoxCountData int64
	var imageCountData int64
	var chatCountData int64


	// Interaction Count
	err = RunQuery("SELECT count() AS count FROM events WHERE eventUser = $userId AND eventId NOT IN [12, 13] GROUP ALL", &interactionCountData, userId);
	if err != nil {
		panic(err)
	}

	// Good Bot Count
	err = RunQuery("SELECT count() AS count FROM events WHERE eventUser = $userId AND eventId IN [34] GROUP ALL", &goodBotCountData, userId);
	if err != nil {
		panic(err)
	}

	// Bad Bot Count
	err = RunQuery("SELECT count() AS count FROM events WHERE eventUser = $userId AND eventId IN [33] GROUP ALL", &badBotCountData, userId);
	if err != nil {
		panic(err)
	}

	// Loot Box Count
	err = RunQuery("SELECT count() AS count FROM events WHERE eventUser = $userId AND eventId IN [9] GROUP ALL", &lootBoxCountData, userId);
	if err != nil {
		panic(err)
	}

	// Image Count
	err = RunQuery("SELECT count() AS count FROM events WHERE eventUser = $userId AND eventId IN [16] GROUP ALL", &imageCountData, userId);
	if err != nil {
		panic(err)
	}

	// Chat Count
	err = RunQuery("SELECT count() AS count FROM events WHERE eventUser = $userId AND eventId IN [12, 13] GROUP ALL", &chatCountData, userId);
	if err != nil {
		panic(err)
	}

	var myStats persistance.MyStats

	myStats.InteractionCount = int(interactionCountData)
	myStats.GoodBotCount = int(goodBotCountData)
	myStats.BadBotCount = int(badBotCountData)
	myStats.LootBoxCount = int(lootBoxCountData)
	myStats.ImageCount = int(imageCountData)
	myStats.ChatCount = int(chatCountData)

	myStats.ImageTokens = user.UserStats.ImageTokens
	myStats.BonusStreak = user.UserStats.BonusStreak
	myStats.LastBonus = user.UserStats.LastBonus

	return myStats
}
