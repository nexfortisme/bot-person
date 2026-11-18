package external

import "strings"

// BotRoute is a named chat model exposed as a slash command and an @mention.
type BotRoute struct {
	Command     string
	Mention     string
	Model       string
	ChoiceName  string
	Label       string
	Description string
	Help        string
}

func BotRoutes() []BotRoute {
	return []BotRoute{
		{
			Command:     "bot-gemini",
			Mention:     "gemini",
			Model:       "google/gemini-4-carbon",
			ChoiceName:  "Bot Gemini",
			Label:       "Gemini",
			Description: "Ask Gemini through OpenRouter.",
			Help:        "The `/bot-gemini` command asks Gemini (`google/gemini-4-carbon`) through OpenRouter and streams the reply. Reply to that message to continue the conversation. You can also mention `@gemini`.",
		},
		{
			Command:     "bot-muse",
			Mention:     "muse",
			Model:       "meta/muse-ember",
			ChoiceName:  "Bot Muse",
			Label:       "Muse",
			Description: "Ask Muse through OpenRouter.",
			Help:        "The `/bot-muse` command asks Muse (`meta/muse-ember`) through OpenRouter and streams the reply. Reply to that message to continue the conversation. You can also mention `@muse`.",
		},
		{
			Command:     "bot-claude",
			Mention:     "claude",
			Model:       "anthropic/claude-sonnet-5.5",
			ChoiceName:  "Bot Claude",
			Label:       "Claude",
			Description: "Ask Claude through OpenRouter.",
			Help:        "The `/bot-claude` command asks Claude (`anthropic/claude-sonnet-5.5`) through OpenRouter and streams the reply. Reply to that message to continue the conversation. You can also mention `@claude`.",
		},
		{
			Command:     "bot-grok",
			Mention:     "grok",
			Model:       "x-ai/grok-4.7",
			ChoiceName:  "Bot Grok",
			Label:       "Grok",
			Description: "Ask Grok through OpenRouter.",
			Help:        "The `/bot-grok` command asks Grok (`x-ai/grok-4.7`) through OpenRouter and streams the reply. Reply to that message to continue the conversation. You can also mention `@grok`.",
		},
		{
			Command:     "bot-bob",
			Mention:     "bob",
			Model:       "ibm-granite/granite-4.2-8b",
			ChoiceName:  "Bot Bob",
			Label:       "Bob",
			Description: "Ask Bob, using IBM Granite, through OpenRouter.",
			Help:        "The `/bot-bob` command asks Bob, who uses IBM Granite (`ibm-granite/granite-4.2-8b`), through OpenRouter and streams the reply. Reply to that message to continue the conversation. You can also mention `@bob`.",
		},
	}
}

func RouteByCommand(command string) (BotRoute, bool) {
	canonical := canonicalCommand(command)
	for _, route := range BotRoutes() {
		if route.Command == canonical {
			return route, true
		}
	}
	return BotRoute{}, false
}

func RouteByMention(mention string) (BotRoute, bool) {
	mention = strings.ToLower(strings.TrimSpace(mention))
	for _, route := range BotRoutes() {
		if route.Mention == mention {
			return route, true
		}
	}
	return BotRoute{}, false
}

func canonicalCommand(command string) string {
	return strings.ReplaceAll(strings.ToLower(strings.TrimSpace(command)), "_", "-")
}
