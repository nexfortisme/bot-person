package enums

type Attribute string

const (
	BOT_PREPROMPT Attribute = "Bot Pre-Prompt"
)

func (a Attribute) String() string {
	return string(a)
}

func (a Attribute) GetAttribute() Attribute {
	return a
}
