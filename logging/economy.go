package logging

func AddImageTokens(tokenAmount float64, userId string) bool {

	user, _ := getUser(userId)

	if user.MessageCount == -1 {
		botTracking.UserStats = append(botTracking.UserStats, UserStruct{userId, UserStatsStruct{0, 0, 0, 0, tokenAmount}})
		return true
	} else {
		user.ImageTokens += tokenAmount
		return updateUser(user, userId)
	}

}

func TransferrImageTokens(tokenAmount float64, fromUserId string, toUserId string) bool {

	fromUser, _ := getUser(fromUserId)
	toUser, _ := getUser(toUserId)

	// Checking to see if user exists
	if fromUser.MessageCount == -1 {
		return false
	} else {

		// The User exists but is trying to transferr more tokens then they have
		if fromUser.ImageTokens < tokenAmount {
			return false
		}

		// Checking to see if the toUser exists
		if toUser.MessageCount == -1 {

			// toUser doesn't exist
			// Creates user and assigns them the number of tokens that is being transferred
			botTracking.UserStats = append(botTracking.UserStats, UserStruct{toUserId, UserStatsStruct{0, 0, 0, 0, tokenAmount}})
			fromUser.ImageTokens -= tokenAmount
			return updateUser(fromUser, fromUserId)
		} else {
			toUser.ImageTokens += tokenAmount
			fromUser.ImageTokens -= tokenAmount
			return updateUser(toUser, toUserId) && updateUser(fromUser, fromUserId)
		}
	}

}

func UseImageToken(userId string) bool {

	user, _ := getUser(userId)

	if user.MessageCount == -1 {
		return false
	} else {
		if user.ImageTokens < 1 {
			return false
		} else {
			user.ImageTokens--
			return updateUser(user, userId);
		}
	}

}

func UserHasTokens(userId string) bool {
	user, _ := getUser(userId)

	if user.MessageCount == -1 {
		return false
	} else {
		if user.ImageTokens <= 0 {
			return false
		} else {
			return true
		}
	}

}

func GetUserTokenCount(userId string) float64 {
	user, _ := getUser(userId)

	if user.MessageCount == -1 {
		return 0
	} else {
		return user.ImageTokens
	}
}

func SetUserTokenCount(userId string, tokenAmount float64) bool {
	user, _ := getUser(userId)

	if user.MessageCount == -1 {
		botTracking.UserStats = append(botTracking.UserStats, UserStruct{userId, UserStatsStruct{0, 0, 0, 0, tokenAmount}})
		return true;
	} else {
		user.ImageTokens = tokenAmount;
		return updateUser(user, userId);
	}
}

func RemoveUserTokens(userId string, tokenAmount float64) bool {
	user, _ := getUser(userId)

	if user.MessageCount == -1 {
		botTracking.UserStats = append(botTracking.UserStats, UserStruct{userId, UserStatsStruct{0, 0, 0, 0, 0}})
		return true;
	} else {
		if (user.ImageTokens - tokenAmount) <= 0 {
			user.ImageTokens = 0;
			return updateUser(user, userId);
		} else {
			user.ImageTokens -= tokenAmount;
			return updateUser(user, userId);
		}
	}
}
