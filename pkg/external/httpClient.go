package external

import (
	"net"
	"net/http"
	"time"
)

// llmHTTPClient is shared by the local-LLM chat completion calls. The short
// dial timeout makes an unreachable endpoint (e.g. host.docker.internal when
// Ollama is not running on the host) fail within seconds instead of hanging
// for the OS-level connect timeout of 30+ seconds. No overall request timeout
// is set because streamed responses can legitimately run for minutes.
var llmHTTPClient = &http.Client{
	Transport: &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout: 10 * time.Second,
		ForceAttemptHTTP2:   true,
		MaxIdleConns:        20,
		IdleConnTimeout:     90 * time.Second,
	},
}
