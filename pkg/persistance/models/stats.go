package persistance

import "time"

type UserStats struct {
	MessageCount int     `json:"MessageCount"`
	GoodBotCount int     `json:"GoodBotCount"`
	BadBotCount  int     `json:"BadBotCount"`
	ImageCount   int     `json:"imageCount"`
	ImageTokens  float64 `json:"imageTokens"`

	LastBonus   time.Time `json:"lastBonus"`
	BonusStreak int       `json:"bonusStreak"`
	
	HoldStreakTimer time.Time `json:"holdStreakTimer"`

	SaveStreakTokens int `json:"saveStreakTokens"`

	Stocks []Stock `json:"stocks"`
}