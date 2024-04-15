package persistance

import "time"

type GlobalStats struct {
	TotalServers              int
	ImagesRequested           int
	LastImageRequest          time.Time
	LongestBonusStreakRecord  int
	CurrentLongestBonusStreak int
	LongestBonusStreakUser    string
	TotalTokensInCirculation  float64
	TotalUsers                int
	GoodBotCount              int
	BadBotCount               int
}
