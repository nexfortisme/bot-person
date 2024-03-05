package persistance

type User struct {
	UserId    string    `json:"username"`
	UserStats UserStats `json:"userStats"`
}

type DBUser struct {
	tableName       struct{} `pg:"tbl_bp_user"`
	Id              string   `pg:"type: uuid"`
	Username        string   `pg:"username"`
	Date_Created    string   `pg:"date_created"`
	Date_Modified   string   `pg:"date_modified"`
	Discord_User_ID string   `pg:"discord_user_id"`

	// Not going to be persisted, just fetched from the DB
	UserStats DBUserStats `pg:"-"`
}
