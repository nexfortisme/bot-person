package persistance

import (
	"errors"
	"fmt"
	"main/pkg/util"
	"math/rand"
	"strconv"
	"strings"
	"time"

	persistance "main/pkg/persistance/eums"
	models "main/pkg/persistance/models"

	"github.com/bwmarrin/discordgo"
	"github.com/hako/durafmt"
)

func AddBotPersonTokens(tokenAmount float64, userId string) bool {

	user, _ := GetUser(userId)
	user.UserStats.ImageTokens += tokenAmount

	UpdateUser(*user)

	return true
}

func RemoveBotPersonTokens(tokenAmount float64, userId string) bool {

	user, _ := GetUser(userId)
	user.UserStats.ImageTokens -= tokenAmount

	if user.UserStats.ImageTokens < 0 {
		user.UserStats.ImageTokens = 0
	}

	UpdateUser(*user)

	return true
}

func TransferBotPersonTokens(tokenAmount float64, fromUserId string, toUserId string) bool {

	fromUser, fromErr := GetUser(fromUserId)
	toUser, toErr := GetUser(toUserId)

	if fromErr != nil || toErr != nil {
		return false
	}

	// The User exists but is trying to transfer more tokens than they have
	if fromUser.UserStats.ImageTokens < tokenAmount {
		return false
	}

	toUser.UserStats.ImageTokens += tokenAmount
	fromUser.UserStats.ImageTokens -= tokenAmount
	return UpdateUser(*toUser) && UpdateUser(*fromUser)
}

func GetUserReward(userId string) (float64, persistance.RewardStatus, error) {

	// Getting user and setting necessary variables
	user, err := GetUser(userId)
	if err != nil {
		panic(err)
	}

	modifier := 1

	// Checking to see if the user has a LastBonus time
	if (user.UserStats.LastBonus != time.Time{}) {

		// Checking diff between now and lastBonus time
		diff := time.Since(user.UserStats.LastBonus)

		// if diff is less than 1 day (86400 seconds), then throws error
		if diff.Seconds() <= 86400 {

			// Doing math to for countdown to next bonus
			nextBonus := user.UserStats.LastBonus.AddDate(0, 0, 1)
			timeToNextBonus := time.Until(nextBonus)
			formattedString := durafmt.Parse(timeToNextBonus).LimitFirstN(3)

			errString := "Please try again in: " + formattedString.String()

			return -1.0, persistance.TOO_EARLY, errors.New(errString)
		}

		timeWindow := user.UserStats.LastBonus.Add(time.Hour * 48)
		timeWindowDiff := time.Since(timeWindow)

		if timeWindowDiff > 0 {
			return -1.0, persistance.MISSED, errors.New("Streak Missed")
		}

	}

	// Get Final Bonus Reward
	finalReward := util.GetUserBonus(5, 50, modifier)

	// Updating User Record
	user.UserStats.LastBonus = time.Now()
	user.UserStats.ImageTokens += finalReward
	user.UserStats.BonusStreak++

	if !UpdateUser(*user) {
		return -1.0, persistance.AVAILABLE, errors.New("error updating user record")
	} else {
		return finalReward, persistance.AVAILABLE, nil
	}

}

func BuyLootbox(userId string) (float64, int, error) {

	user, err := GetUser(userId)

	if err != nil {
		return -1, -1, err
	}

	if user.UserStats.ImageTokens < 5 {
		return -1, -1, errors.New("you do not have the 5 tokens needed to purchase a lootbox")
	} else {
		user.UserStats.ImageTokens -= 5
	}

	random := rand.New(rand.NewSource(time.Now().UnixMilli()))
	lootboxSeed := random.Intn(9999999999) + 1000000000

	val := hashLootBoxSeed(lootboxSeed)
	reward := 0.0

	if val <= 7992 {
		reward += 3.63
	} else if val > 7992 && val <= 9590 {
		reward += 8
	} else if val > 9590 && val <= 9910 {
		reward += 15
	} else if val > 9910 && val <= 9974 {
		reward += 25
	} else if val > 9974 && val <= 10000 {
		reward += 50
	}

	user.UserStats.ImageTokens += float64(reward)

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

func APictureIsWorthAThousand(incomingMessage string, m *discordgo.MessageCreate) {

	// Looking at messages in the channel and returning WORD_COUNT / 1000 number of tokens
	// ie. A picture is worth 1000 words
	wordCount := len(strings.Fields(incomingMessage))
	tokenValue := fmt.Sprintf("%.2f", (float64(wordCount) / 1000.0))
	tokenAddAmount, _ := strconv.ParseFloat(tokenValue, 64)

	AddBotPersonTokens(tokenAddAmount, m.Author.ID)
}

func AddStock(userId string, stockTicker string, quantity float64) error {

	user, err := GetUser(userId)

	if err != nil {
		return err
	}

	userPortfolio := user.UserStats.Stocks

	for index, element := range userPortfolio {
		if element.StockTicker == stockTicker {
			element.StockCount += quantity
			userPortfolio[index] = element
			UpdateUser(*user)
			return nil
		}
	}

	newStock := models.Stock{StockTicker: stockTicker, StockCount: quantity}
	user.UserStats.Stocks = append(user.UserStats.Stocks, newStock)
	UpdateUser(*user)

	return nil
}

func RemoveStock(userId string, stockTicker string, quantity float64) error {

	user, err := GetUser(userId)

	if err != nil {
		return err
	}

	userPortfolio := user.UserStats.Stocks

	for index, element := range userPortfolio {
		if element.StockTicker == stockTicker {
			element.StockCount -= quantity
			userPortfolio[index] = element
			UpdateUser(*user)
			return nil
		}
	}

	return errors.New("stock not found")
}

func GetUserStock(userId string, stockTicker string) (models.Stock, error) {

	user, err := GetUser(userId)

	if err != nil {
		return models.Stock{}, err
	}

	userPortfolio := user.UserStats.Stocks

	for _, element := range userPortfolio {
		if element.StockTicker == stockTicker {
			return element, nil
		}
	}

	return models.Stock{}, errors.New("stock not found")

}
