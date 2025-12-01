package persistance

import (
	"errors"
	"fmt"
	"main/pkg/util"
	"math/rand"
	"time"

	persistance "main/pkg/persistance/eums"

	"github.com/hako/durafmt"
)

func AddBotPersonTokens(tokenAmount int, userId string) bool {

	user, err := GetUser(userId)
	if err != nil {
		return false
	}

	user.ImageTokens += tokenAmount
	UpdateUser(*user)

	return true
}

func RemoveBotPersonTokens(tokenAmount int, userId string) bool {

	user, _ := GetUser(userId)
	user.ImageTokens -= tokenAmount

	if user.ImageTokens < 0 {
		user.ImageTokens = 0
	}

	UpdateUser(*user)

	return true
}

func TransferBotPersonTokens(tokenAmount int, fromUserId string, toUserId string) bool {

	fromUser, fromErr := GetUser(fromUserId)
	toUser, toErr := GetUser(toUserId)

	if fromErr != nil || toErr != nil {
		return false
	}

	// The User exists but is trying to transfer more tokens than they have
	if fromUser.ImageTokens < tokenAmount {
		return false
	}

	toUser.ImageTokens += tokenAmount
	fromUser.ImageTokens -= tokenAmount
	return UpdateUser(*toUser) && UpdateUser(*fromUser)
}

func GetUserReward(userId string) (int, persistance.RewardStatus, error) {

	// Getting user and setting necessary variables
	user, err := GetUser(userId)
	if err != nil {
		panic(err)
	}

	modifier := 1

	fmt.Printf("User: %+v\n", user)

	// Checking to see if the user has a LastBonus time
	if user.LastBonus != "" {

		// Parse LastBonus string into time and check diff between now and lastBonus time
		lastBonusTime, parseErr := time.Parse(time.RFC3339, user.LastBonus)
		if parseErr == nil {
			diff := time.Since(lastBonusTime)

			fmt.Println("Time Diff: " + diff.String() + " Seconds: " + fmt.Sprintf("%.0f", diff.Seconds()))

			// if diff is less than 1 day (86400 seconds), then throws error
			if diff.Seconds() <= 86400 {
				// Doing math for countdown to next bonus
				nextBonus := lastBonusTime.AddDate(0, 0, 1)
				timeToNextBonus := time.Until(nextBonus)
				formattedString := durafmt.Parse(timeToNextBonus).LimitFirstN(3)

				errString := "Please try again in: " + formattedString.String()

				return -1.0, persistance.TOO_EARLY, errors.New(errString)
			}

			timeWindow := lastBonusTime.Add(time.Hour * 48)
			timeWindowDiff := time.Since(timeWindow)

			if timeWindowDiff > 0 {
				return -1.0, persistance.MISSED, errors.New("streak missed")
			}
		}

	}

	// Get Final Bonus Reward
	finalReward := util.GetUserBonus(1, 4, modifier)

	// Updating User Record
	user.LastBonus = time.Now().Format(time.RFC3339)
	user.ImageTokens += finalReward
	user.BonusStreak++

	if !UpdateUser(*user) {
		return -1.0, persistance.AVAILABLE, errors.New("error updating user record")
	} else {
		return finalReward, persistance.AVAILABLE, nil
	}

}

func BuyLootbox(userId string) (int, int, error) {

	user, err := GetUser(userId)

	if err != nil {
		return -1, -1, err
	}

	if user.ImageTokens < 5 {
		return -1, -1, errors.New("you do not have the 5 tokens needed to purchase a lootbox")
	} else {
		user.ImageTokens -= 5
	}

	random := rand.New(rand.NewSource(time.Now().UnixMilli()))
	lootboxSeed := random.Intn(9999999999) + 1000000000

	val := hashLootBoxSeed(lootboxSeed)
	reward := 0

	if val <= 7992 {
		reward += 3
	} else if val > 7992 && val <= 9590 {
		reward += 8
	} else if val > 9590 && val <= 9910 {
		reward += 15
	} else if val > 9910 && val <= 9974 {
		reward += 25
	} else if val > 9974 && val <= 10000 {
		reward += 50
	}

	user.ImageTokens += reward

	if !UpdateUser(*user) {
		return -1, -1, errors.New("error updating user record")
	} else {
		return reward, lootboxSeed, nil
	}

}

func hashLootBoxSeed(bar int) int {

	/*
	* Blue – 100 items – 79.92%
	* Purple – 20 items – 15.98%
	* Pink – 4 items – 3.2%
	* Red – 0.8 items – 0.64%
	* Yellow – 0.32 items – 0.26%
	 */

	modifier := 1000000000

	total := 1

	for i := 1; i < 11; i++ {
		num := (bar / modifier) % 10
		modifier /= 10

		if i <= 8 {
			if num != 0 {
				total *= num
			} else {
				total += 10
			}
		} else {
			if num != 0 {
				total *= num
			} else {
				total += 10
			}
		}

	}

	return total % 10000
}

// func APictureIsWorthAThousand(incomingMessage string, m *discordgo.MessageCreate) {

// 	// Looking at messages in the channel and returning WORD_COUNT / 1000 number of tokens
// 	// ie. A picture is worth 1000 words
// 	wordCount := len(strings.Fields(incomingMessage))
// 	tokenValue := fmt.Sprintf("%.2f", (float64(wordCount) / 1000.0))
// 	tokenAddAmount, _ := strconv.ParseInt(tokenValue, 10, 64)

// 	AddBotPersonTokens(tokenAddAmount, m.Author.ID)
// }
