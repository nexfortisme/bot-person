package persistance

import "time"

type MyStats struct {
	InteractionCount int       `json:"interactionCount"`
	GoodBotCount     int       `json:"goodBotCount"`
	BadBotCount      int       `json:"badBotCount"`
	ImageTokens      float64   `json:"imageTokens"`
	BonusStreak      int       `json:"bonusStreak"`
	LastBonus        time.Time `json:"lastBonus"`
	LootBoxCount     int       `json:"lootBoxCount"`
	ImageCount       int       `json:"imageCount"`
	ChatCount        int       `json:"chatCount"`
}

type UserStats struct {
	UserId string `json:"UserId"`

	InteractionCount int `json:"InteractionCount"`
	ChatCount        int `json:"ChatCount"`
	GoodBotCount     int `json:"GoodBotCount"`
	BadBotCount      int `json:"BadBotCount"`

	ImageCount int `json:"ImageCount"`

	LootBoxCount int `json:"LootBoxCount"`
}

type Stock struct {
	StockTicker string  `json:"stockTicker"`
	StockCount  float64 `json:"stockCount"`
}

type BotTracking struct {
	MessageCount int    `json:"MessageCount"`
	GoodBotCount int    `json:"GoodBotCount"`
	BadBotCount  int    `json:"BadBotCount"`
	UserStats    []User `json:"UserTracking"`
}

type UserEventCount struct {
	Count int `json:"count"`
}

type User struct {
	ID          string  `json:"id,omitempty"`
	UserId      string  `json:"UserId"`
	Username    string  `json:"Username"`
	ImageTokens float64 `json:"ImageTokens"`
	BonusStreak int     `json:"BonusStreak"`
	LastBonus   string  `json:"LastBonus"`
}