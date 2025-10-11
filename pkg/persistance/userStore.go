package persistance

import (
	"time"
)

func GetUser(userId string) (*User, error) {

	user := User{}

	err := RunQuery("SELECT * FROM users WHERE UserId = ?", &user, userId)
	if err != nil {
		panic(err)
	}

	// If the user is not found, create a new user
	if user.ID == "" || user.ID == "0" {

		newUser := User{}

		if userId != "SYSTEM" {
			// discordSession := stateService.GetDiscordSession()
			// discordUser, _ := discordSession.User(userId)
			// try {
			// 	newUser.Username = discordUser.
			// } catch(err Error){
			// 	newUser.Username = "SYSTEM"
			// }
		} else {
			newUser.Username = "SYSTEM"
		}

		newUser.UserId = userId
		newUser.ImageTokens = 10 // Starting token amount
		newUser.BonusStreak = 0
		newUser.LastBonus = time.Now().String()

		err = RunQuery("INSERT INTO users (UserId, ImageTokens, BonusStreak, LastBonus) VALUES (?, ?, ?, ?)", nil, userId, newUser.ImageTokens, newUser.BonusStreak, time.Now().String())
		if err != nil {
			panic(err)
		}

		return &newUser, nil
	}

	return &user, nil
}

func UpdateUser(updateUser User) bool {

	err := RunQuery("UPDATE users SET Username = ?, ImageTokens = ?, BonusStreak = ?, LastBonus = ? WHERE UserId = ?", nil, updateUser.Username, updateUser.ImageTokens, updateUser.BonusStreak, updateUser.LastBonus, updateUser.UserId)
	if err != nil {
		panic(err)
	}

	return true
}

func GetUserStatsObj(userId string) (*UserStats, error) {
	userStats := UserStats{}

	err := RunQuery("SELECT * FROM userStats WHERE UserId = ?", userStats, userId)
	if err != nil {
		panic(err)
	}

	return &userStats, nil
}

func UpdateUserStats(userStats UserStats) bool {
	err := RunQuery("UPDATE userStats SET InteractionCount = InteractionCount + 1, ChatCount = ?, GoodBotCount = ?, BadBotCount = ?, ImageCount = ?, LootBoxCount = ? WHERE UserId = ?", nil, userStats.ChatCount, userStats.GoodBotCount, userStats.BadBotCount, userStats.ImageCount, userStats.LootBoxCount, userStats.UserId)
	if err != nil {
		panic(err)
	}

	return true
}
