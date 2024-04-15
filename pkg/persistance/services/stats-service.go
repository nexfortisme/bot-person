package persistance

import (

	// "log"

	persistanceModels "main/pkg/persistance/models"

	// loggingTypes "main/pkg/logging/enums"
	// logging "main/pkg/logging/services"

	"github.com/bwmarrin/discordgo"
	// "github.com/surrealdb/surrealdb.go"
	// "github.com/go-pg/pg/v10"
)

type BPInteraction int

func GetUserStats(discord_user_id string, session *discordgo.Session) (persistanceModels.UserStats, error) {

	user, err := GetUserByDiscordUserId(discord_user_id, session)
	if err != nil {
		return persistanceModels.UserStats{}, err
	}

	return user.UserStats, nil
}

// func UpsertUserStats(stats persistanceModels.UserStats) (persistanceModels.UserStats, error) {

// 	// databaseConnection := persistance.GetDB()

// 	// _, err := databaseConnection.Model(&stats).OnConflict("(bp_user_id) DO UPDATE").Insert()
// 	// if err != nil {
// 	// 	return persistanceModels.DBUserStats{}, err
// 	// }

// 	return stats, nil
// }

// func GetGlobalStats() (persistanceModels.GlobalStats, error) {

// 	var globalStats persistanceModels.GlobalStats = persistanceModels.GlobalStats{}

// 	databaseConnection := persistance.GetDB()
// 	discordSession := state.GetDiscordSession()

// 	globalStats.TotalServers = len(discordSession.State.Guilds)

// 	// Get user by ID
// 	data, err := databaseConnection.Query("SELECT COUNT(stats.", nil)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// Unmarshal data
// 	selectedUser := new(User)
// 	err = surrealdb.Unmarshal(data, &selectedUser)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// _, err := databaseConnection.QueryOne(pg.Scan(&globalStats.ImagesRequested), "SELECT COUNT(*) FROM tbl_bp_event WHERE event_type = ?", loggingTypes.COMMAND_IMAGE)
// 	// _, err = databaseConnection.QueryOne(pg.Scan(&globalStats.LastImageRequest), "SELECT createDate FROM tbl_bp_event WHERE event_type = ? ORDER BY date_created DESC LIMIT 1", loggingTypes.COMMAND_IMAGE)
// 	// _, err = databaseConnection.QueryOne(pg.Scan(&globalStats.LongestBonusStreakRecord), "SELECT event_value FROM tbl_bp_event WHERE event_type = ? ORDER BY event_value DESC LIMIT 1", loggingTypes.USER_SET_BONUS_STREAK)
// 	// _, err = databaseConnection.QueryOne(pg.Scan(&globalStats.CurrentLongestBonusStreak), "SELECT bonus_streak FROM tbl_bp_user_stats ORDER BY bonus_streak DESC LIMIT 1")
// 	// _, err = databaseConnection.QueryOne(pg.Scan(&globalStats.LongestBonusStreakUser), "SELECT user.username FROM tbl_bp_user_stats tbus JOIN tbl_bp_user user ON tbus.bp_user_id = user.id WHERE tbus.bonus_streak = (SELECT MAX(bonus_streak) FROM tbl_bp_user_stats) LIMIT 1")
// 	// _, err = databaseConnection.QueryOne(pg.Scan(&globalStats.TotalTokensInCirculation), "SELECT COUNT(token_balance) FROM tbl_bp_user_stats")
// 	// _, err = databaseConnection.QueryOne(pg.Scan(&globalStats.TotalUsers), "SELECT COUNT(*) FROM tbl_bp_user")
// 	// _, err = databaseConnection.QueryOne(pg.Scan(&globalStats.GoodBotCount), "SELECT COUNT(good_bot_count) FROM tbl_bp_user_stats")
// 	// _, err = databaseConnection.QueryOne(pg.Scan(&globalStats.BadBotCount), "SELECT COUNT(bacd_bot_count) FROM tbl_bp_user_stats")
// 	// if err != nil {
// 	// 	log.Fatalf("Error executing count query: %v", err)

// 	// 	logging.LogEvent(loggingTypes.INTERNAL_ERROR, "Error executing count query: %v", "System", "System", nil)

// 	// 	return globalStats, err
// 	// }

// 	return globalStats, nil
// }
