package persistance

import "time"

type UserStats struct {
	ImageCount   int     `json:"ImageCount"`
	ImageTokens  float64 `json:"ImageTokens"`

	LastBonus   time.Time `json:"LastBonus"`
	BonusStreak int       `json:"BonusStreak"`
	
	Stocks []Stock `json:"Stocks"`
}