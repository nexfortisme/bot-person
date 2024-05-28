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
