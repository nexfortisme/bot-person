package persistance

import "time"

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
	MessageCount int     `json:"MessageCount"`
	GoodBotCount int     `json:"GoodBotCount"`
	BadBotCount  int     `json:"BadBotCount"`
	ImageCount   int     `json:"imageCount"`
	ImageTokens  float64 `json:"imageTokens"`

	LastBonus   time.Time `json:"lastBonus"`
	BonusStreak int       `json:"bonusStreak"`
	
	HoldStreakTimer time.Time `json:"holdStreakTimer"`

	SaveStreakTokens int `json:"saveStreakTokens"`

	Stocks []UserStock `json:"stocks"`
}

type UserStock struct {
	StockTicker string  `json:"stockTicker"`
	StockCount  float64 `json:"stockCount"`
}
