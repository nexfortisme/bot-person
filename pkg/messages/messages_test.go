package messages

import "testing"

func TestParseAddTokensCommand(t *testing.T) {
	tests := []struct {
		name       string
		command    string
		wantUserID string
		wantAmount int
		wantError  bool
	}{
		{name: "standard mention", command: "!addTokens <@12345> 10", wantUserID: "12345", wantAmount: 10},
		{name: "nickname mention", command: "!addTokens <@!12345> 7", wantUserID: "12345", wantAmount: 7},
		{name: "extra whitespace", command: "!addTokens   <@12345>   2", wantUserID: "12345", wantAmount: 2},
		{name: "not a mention", command: "!addTokens person 2", wantError: true},
		{name: "malformed mention", command: "!addTokens <@> 2", wantError: true},
		{name: "non-numeric user", command: "!addTokens <@abc> 2", wantError: true},
		{name: "invalid amount", command: "!addTokens <@12345> nope", wantError: true},
		{name: "negative amount", command: "!addTokens <@12345> -2", wantError: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			userID, amount, err := parseAddTokensCommand(test.command)
			if test.wantError {
				if err == nil {
					t.Fatalf("parseAddTokensCommand(%q) unexpectedly succeeded", test.command)
				}
				return
			}
			if err != nil {
				t.Fatalf("parseAddTokensCommand(%q): %v", test.command, err)
			}
			if userID != test.wantUserID || amount != test.wantAmount {
				t.Fatalf("parseAddTokensCommand(%q) = (%q, %d), want (%q, %d)", test.command, userID, amount, test.wantUserID, test.wantAmount)
			}
		})
	}
}
