package util

import (
	"strings"
	"sync"
	"time"
)

const (
	discordMessageCharLimit = 2000
	streamEditInterval      = 750 * time.Millisecond
	streamMinCharDelta      = 24
)

var threadResponseState = struct {
	lock     sync.Mutex
	inFlight map[string]struct{}
}{
	inFlight: make(map[string]struct{}),
}

func TryStartThreadResponse(threadID string) bool {
	normalizedThreadID := strings.TrimSpace(threadID)
	if normalizedThreadID == "" {
		return false
	}

	threadResponseState.lock.Lock()
	defer threadResponseState.lock.Unlock()

	if _, exists := threadResponseState.inFlight[normalizedThreadID]; exists {
		return false
	}

	threadResponseState.inFlight[normalizedThreadID] = struct{}{}
	return true
}

func FinishThreadResponse(threadID string) {
	normalizedThreadID := strings.TrimSpace(threadID)
	if normalizedThreadID == "" {
		return
	}

	threadResponseState.lock.Lock()
	delete(threadResponseState.inFlight, normalizedThreadID)
	threadResponseState.lock.Unlock()
}

func TruncateForDiscord(content string) string {
	if len(content) <= discordMessageCharLimit {
		return content
	}

	return content[:discordMessageCharLimit-3] + "..."
}

func StreamResponseWithThrottledEdits(
	initialDisplay string,
	formatDisplay func(assistantResponse string) string,
	streamFn func(onDelta func(string)) (string, error),
	editFn func(content string) error,
) (string, error) {
	lastPublished := TruncateForDiscord(initialDisplay)
	lastPublishTime := time.Now()
	partialResponse := strings.Builder{}

	assistantResponse, streamErr := streamFn(func(delta string) {
		if delta == "" {
			return
		}

		partialResponse.WriteString(delta)
		candidate := TruncateForDiscord(formatDisplay(partialResponse.String()))
		if candidate == lastPublished {
			return
		}

		now := time.Now()
		charDelta := len(candidate) - len(lastPublished)
		if charDelta < 0 {
			charDelta = -charDelta
		}
		if now.Sub(lastPublishTime) < streamEditInterval && charDelta < streamMinCharDelta {
			return
		}

		if err := editFn(candidate); err == nil {
			lastPublished = candidate
			lastPublishTime = now
		}
	})

	if strings.TrimSpace(assistantResponse) == "" {
		assistantResponse = partialResponse.String()
	}

	finalDisplay := TruncateForDiscord(formatDisplay(assistantResponse))
	if finalDisplay != lastPublished {
		_ = editFn(finalDisplay)
	}

	return assistantResponse, streamErr
}
