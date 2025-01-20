package persistance

import (
	persistance "main/pkg/persistance/models"
)

func GetUser(userId string) (*persistance.User, error) {

	user := persistance.User{}

	err := RunQuery("SELECT * FROM users WHERE UserId = ?", user, userId)
	if err != nil {
		panic(err)
	}

	if user.ID == "" || user.ID == "0" {

		newUser := persistance.User{}

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
		newUser.UserStats.ImageTokens = 50

		err = RunQuery("INSERT INTO users (UserId, Username, UserStats) VALUES (?, ?, ?)", nil, userId, newUser.Username, newUser.UserStats)
		if err != nil {
			panic(err)
		}

		return &newUser, nil
	}

	return &user, nil
}

func UpdateUser(updateUser persistance.User) bool {

	err := RunQuery("UPDATE users SET Username = ?, UserStats = ? WHERE UserId = ?", nil, updateUser.Username, updateUser.UserStats, updateUser.UserId)
	if err != nil {
		panic(err)
	}

	return true
}
