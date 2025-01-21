package persistance

type UserStats struct {
	UserId string `json:"UserId"`

	InteractionCount int `json:"InteractionCount"`
	ChatCount        int `json:"ChatCount"`
	GoodBotCount     int `json:"GoodBotCount"`
	BadBotCount      int `json:"BadBotCount"`

	ImageCount int `json:"ImageCount"`

	LootBoxCount int `json:"LootBoxCount"`
}
