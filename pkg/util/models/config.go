package util

type Config struct {
	OpenAIKey       string   `json:"OpenAIKey"`
	DiscordToken    string   `json:"DiscordToken"`
	DevDiscordToken string   `json:"DevDiscordToken"`
	FinnHubToken    string   `json:"FinnHubToken"`
	AdminIDs        []string `json:"AdminIDs"`
}
