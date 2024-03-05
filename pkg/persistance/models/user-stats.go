package persistance

type DBUserStats struct {
	tableName     struct{} `pg:"tbl_bp_user_stats"`
	ID            string   `pg:"type: uuid"`
	Date_Created  string   `pg:"date_created"`
	Date_Modified string   `pg:"date_modified"`

	BP_User_ID string `pg:"bp_user_id"`

	Message_count  int `pg:"message_count"`
	Good_Bot_Count int `pg:"good_bot_count"`
	Bad_Bot_Count  int `pg:"bad_bot_count"`

	Token_Balance float64    `pg:"token_balance"`
	Last_Bonus    string `pg:"last_bonus"`
	Bonus_Streak  int    `pg:"bonus_streak"`
}
