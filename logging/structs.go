package logging


type BotTracking struct {
	MessageCount int          `json:"MessageCount"`
	GoodBotCount int          `json:"GoodBotCount"`
	BadBotCount  int          `json:"BadBotCount"`
	UserStats    []UserStruct `json:"UserTracking"`
}

type UserStruct struct {
	UserId    string          `json:"username"`
	UserStats UserStatsStruct `json:"userStats"`
}

type UserStatsStruct struct {
	MessageCount int `json:"MessageCount"`
	GoodBotCount int `json:"GoodBotCount"`
	BadBotCount  int `json:"BadBotCount"`
}
