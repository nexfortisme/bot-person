package persistance

import (
	"errors"

	persistance "main/pkg/persistance"
	persistanceModels "main/pkg/persistance/models"

	"github.com/bwmarrin/discordgo"
)

type BPInteraction int

func GetUserStats(discord_user_id string, session *discordgo.Session) (persistanceModels.DBUserStats, error) {

	dbUser, err := GetUser(discord_user_id, session)
	if err != nil {
		return persistanceModels.DBUserStats{}, errors.New("error fetching user")
	}

	var userDBStats persistanceModels.DBUserStats

	databaseConnection := persistance.GetDB()
	defer databaseConnection.Close()

	err = databaseConnection.Model(&userDBStats).Where("bd_user_id = ?", dbUser.Id).Select()
	if err != nil {

		userDBStats = persistanceModels.DBUserStats{
			BP_User_ID:    dbUser.Id,
			Token_Balance: 25, // Default token balance
		}

		_, err := UpsertUserStats(userDBStats)

		return userDBStats, err
	}

	return userDBStats, nil
}

func UpsertUserStats(stats persistanceModels.DBUserStats) (persistanceModels.DBUserStats, error) {

	databaseConnection := persistance.GetDB()

	_, err := databaseConnection.Model(&stats).OnConflict("(bp_user_id) DO UPDATE").Insert()
	if err != nil {
		return persistanceModels.DBUserStats{}, err
	}

	return stats, nil
}
