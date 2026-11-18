package systemprompts

import (
	"strings"

	_ "embed"
)

//go:embed chatgpt_system_prompt.md
var ChatGPT string

//go:embed gemini_system_prompt.md
var Gemini string

//go:embed claude_system_prompt.md
var Claude string

//go:embed grok_system_prompt.md
var Grok string

//go:embed muse_system_prompt.md
var Muse string

//go:embed bob_system_prompt.md
var Bob string

// ByMention returns the embedded system prompt for a routed bot mention.
// Anything else gets the ChatGPT prompt, which is what /bot-gpt uses.
func ByMention(mention string) string {
	switch strings.ToLower(strings.TrimSpace(mention)) {
	case "gemini":
		return strings.TrimSpace(Gemini)
	case "claude":
		return strings.TrimSpace(Claude)
	case "grok":
		return strings.TrimSpace(Grok)
	case "muse":
		return strings.TrimSpace(Muse)
	case "bob":
		return strings.TrimSpace(Bob)
	default:
		return strings.TrimSpace(ChatGPT)
	}
}
