package persistance

type User struct {
	ID        string    `json:"id,omitempty"`
	UserId    string    `json:"UserId"`
	UserStats UserStats `json:"UserStats"`
	Username  string    `json:"Username"`
}
