package persistance

type User struct {
	ID          string  `json:"id,omitempty"`
	UserId      string  `json:"UserId"`
	Username    string  `json:"Username"`
	ImageTokens float64 `json:"ImageTokens"`
	BonusStreak int     `json:"BonusStreak"`
	LastBonus   string  `json:"LastBonus"`
}
