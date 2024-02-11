package persistance

type User struct {
	UserId    string    `json:"username"`
	UserStats UserStats `json:"userStats"`
}

type DBUser struct {
	Id              string `pg:"type: uuid"`
	Username        string
	Date_Created    string
	Date_Modified   string
	Discord_User_ID string
}
