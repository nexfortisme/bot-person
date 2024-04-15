package persistance

type User struct {
	DiscordUserId string    `json:"username"`
	UserStats     UserStats `json:"userStats"`
}
