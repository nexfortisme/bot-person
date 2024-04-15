package persistance

type UserStats struct {
	UserId string `json:"userId"`

	Date_Created  string `json:"dateCreated"`
	Date_Modified string `json:"dateModified"`

	// BP_User_ID string `pg:"bp_user_id"`

	Message_count  int `json:"messageCount"`
	Good_Bot_Count int `json:"goodBotCount"`
	Bad_Bot_Count  int `json:"badBotCount"`

	Token_Balance float64 `json:"tokenBalance"`
	Last_Bonus    string  `json:"lastBonus"`
	Bonus_Streak  int     `json:"bonusStreak"`
}
