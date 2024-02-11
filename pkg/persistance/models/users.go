package persistance

type User struct {
	UserId    string `json:"username"`
	UserStats UserStats `json:"userStats"`
}