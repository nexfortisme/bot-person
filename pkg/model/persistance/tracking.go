package persistance

type BotTracking struct {
	MessageCount int    `json:"MessageCount"`
	GoodBotCount int    `json:"GoodBotCount"`
	BadBotCount  int    `json:"BadBotCount"`
	UserStats    []User `json:"UserTracking"`
}
