package persistance

type DBUserStats struct {
	ID            string
	Date_Created  string
	Date_Modified string

	BP_User_ID string

	Message_count  int
	Good_Bot_Count int
	Bad_Bot_Count  int

	Token_Balance int
	Last_Bonus    string
	Bonus_Streak  int
}
