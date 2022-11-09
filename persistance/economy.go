package persistance

func AddImageTokens(tokenAmount float64, userId string) bool {

	user, err := getUser(userId)

	if err != nil {
		botTracking.UserStats = append(botTracking.UserStats, UserStruct{userId, UserStatsStruct{0, 0, 0, 0, tokenAmount}})
		return true
	} else {
		user.UserStats.ImageTokens += tokenAmount
		return updateUser(user)
	}

}

func TransferrImageTokens(tokenAmount float64, fromUserId string, toUserId string) bool {

	fromUser, fromErr := getUser(fromUserId)
	toUser, toErr := getUser(toUserId)

	// Checking to see if user exists
	if fromErr != nil {
		return false
	} else {

		// The User exists but is trying to transferr more tokens then they have
		if fromUser.UserStats.ImageTokens < tokenAmount {
			return false
		}

		// Checking to see if the toUser exists
		if toErr != nil {

			// toUser doesn't exist
			// Creates user and assigns them the number of tokens that is being transferred
			botTracking.UserStats = append(botTracking.UserStats, UserStruct{toUserId, UserStatsStruct{0, 0, 0, 0, tokenAmount}})
			fromUser.UserStats.ImageTokens -= tokenAmount
			return updateUser(fromUser)
		} else {
			toUser.UserStats.ImageTokens += tokenAmount
			fromUser.UserStats.ImageTokens -= tokenAmount
			return updateUser(toUser) && updateUser(fromUser)
		}
	}

}

func UseImageToken(userId string) bool {

	user, err := getUser(userId)

	if err != nil {
		return false
	} else {
		if user.UserStats.ImageTokens < 1 {
			return false
		} else {
			user.UserStats.ImageTokens--
			return updateUser(user)
		}
	}

}

// Is this needed, can be canned
func UserHasTokens(userId string) bool {
	user, err := getUser(userId)

	if err != nil {
		return false
	} else {
		if user.UserStats.ImageTokens <= 0 {
			return false
		} else {
			return true
		}
	}

}

func GetUserTokenCount(userId string) float64 {
	user, err := getUser(userId)

	if err != nil {
		return 0
	} else {
		return user.UserStats.ImageTokens
	}
}

func SetUserTokenCount(userId string, tokenAmount float64) bool {
	user, err := getUser(userId)

	if err != nil {
		botTracking.UserStats = append(botTracking.UserStats, UserStruct{userId, UserStatsStruct{0, 0, 0, 0, tokenAmount}})
		return true
	} else {
		user.UserStats.ImageTokens = tokenAmount
		return updateUser(user)
	}
}

func RemoveUserTokens(userId string, tokenAmount float64) bool {
	user, err := getUser(userId)

	if err != nil {
		botTracking.UserStats = append(botTracking.UserStats, UserStruct{userId, UserStatsStruct{0, 0, 0, 0, 0}})
		return true
	} else {
		if (user.UserStats.ImageTokens - tokenAmount) <= 0 {
			user.UserStats.ImageTokens = 0
			return updateUser(user)
		} else {
			user.UserStats.ImageTokens -= tokenAmount
			return updateUser(user)
		}
	}
}
