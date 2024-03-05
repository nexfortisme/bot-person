package persistance

import (
	"fmt"

	persistance "main/pkg/persistance"
	persistanceModels "main/pkg/persistance/models"

	logging "main/pkg/logging/services"
	loggingEvents "main/pkg/logging/enums"

	"time"

	"github.com/bwmarrin/discordgo"
)

func GetUser(discord_user_id string, s *discordgo.Session) (persistanceModels.DBUser, error) {

	var dbUser persistanceModels.DBUser

	databaseConnection := persistance.GetDB()

	err := databaseConnection.Model(&dbUser).Where("discord_user_id = ?", discord_user_id).Select()
	if err != nil {
		
		discordUser, _ := s.User(discord_user_id)

		dbUser = persistanceModels.DBUser{
			Discord_User_ID: discord_user_id,
			Username:        discordUser.Username,
			Date_Created:  	time.Now().String(),
		}
		
		UpsertUser(dbUser);

		return dbUser, err
	}

	return dbUser, nil
}

func UpsertUser(dbUser persistanceModels.DBUser) (persistanceModels.DBUser, error) {
	
	databaseConnection := persistance.GetDB()

	_, err := databaseConnection.Model(&dbUser).OnConflict("(discord_user_id) DO UPDATE").Insert()
	if err != nil {
		logging.LogEvent(loggingEvents.DATABASE_ERROR, "Error on UpsertUser: " + err.Error(), dbUser.Discord_User_ID, "System", nil);
		fmt.Println("Error on UpsertUser: ", err)
		return persistanceModels.DBUser{}, err
	}

	return dbUser, nil
}
