package persistance

import "time"

type MyStats struct {
	InteractionCount int       `json:"interactionCount"`
	GoodBotCount     int       `json:"goodBotCount"`
	BadBotCount      int       `json:"badBotCount"`
	ImageTokens      int   `json:"imageTokens"`
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
	ImageTokens int `json:"ImageTokens"`
	BonusStreak int     `json:"BonusStreak"`
	LastBonus   string  `json:"LastBonus"`
}

type UserAttribute struct {
	ID        string `json:"id"`
	UserId    string `json:"userId"`
	Attribute string `json:"attribute"`
	Value     string `json:"value"`
}
