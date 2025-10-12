package external

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"main/pkg/util"
	"net/http"
	"os"
	"strings"
	"time"
)

func GetSoraRespone(prompt string, duration int) (VideoJob, error) {

	httpClient := &http.Client{}

	requestDataTemplate := `{
		"prompt": "%s",
		"seconds": "%d"
	}`
	requestData := fmt.Sprintf(requestDataTemplate, prompt, duration)

	postRequest, err := http.NewRequest(http.MethodPost, "https://api.openai.com/v1/videos", strings.NewReader(requestData))
	if err != nil {
		return VideoJob{}, errors.New("POST Request Error: " + err.Error())
	}

	postRequest.Header.Add("Content-Type", "application/json")
	postRequest.Header.Add("Authorization", "Bearer "+util.GetOpenAIKey())

	httpResponse, _ := httpClient.Do(postRequest)
	if httpResponse == nil {
		return VideoJob{}, errors.New("API Error: " + httpResponse.Status)
	}

	responseBuffer, _ := io.ReadAll(httpResponse.Body)
	err = util.SaveResponseToFile(responseBuffer, "sora_response.txt")
	if err != nil {
		fmt.Println("Error saving response to file:", err)
	}

	var videoResponse VideoJob

	err = json.Unmarshal([]byte(string(responseBuffer)), &videoResponse)
	if err != nil {
		return VideoJob{}, errors.New("error parsing response")
	}

	if videoResponse.Error != nil {
		return VideoJob{}, errors.New("error: " + videoResponse.Error.Message)
	}


	return videoResponse, nil
}

func GetSoraJobStatus(videoId string) (VideoJob, error) {

	httpClient := &http.Client{}

	getRequest, err := http.NewRequest(http.MethodGet, "https://api.openai.com/v1/videos/"+videoId, nil)
	if err != nil {
		return VideoJob{}, errors.New("GET Request Error: " + err.Error())
	}

	getRequest.Header.Add("Content-Type", "application/json")
	getRequest.Header.Add("Authorization", "Bearer "+util.GetOpenAIKey())

	httpResponse, _ := httpClient.Do(getRequest)
	if httpResponse == nil {
		return VideoJob{}, errors.New("API Error: " + httpResponse.Status)
	}

	responseBuffer, _ := io.ReadAll(httpResponse.Body)
	currentTime := time.Now().Format("2006-01-02-15-04-05")
	err = util.SaveResponseToFile(responseBuffer, fmt.Sprintf("video_status_response_%s.txt", currentTime))
	if err != nil {
		return VideoJob{}, errors.New("Error saving response to file: " + err.Error())
	}

	var videoResponse VideoJob
	err = json.Unmarshal([]byte(string(responseBuffer)), &videoResponse)
	if err != nil {
		return VideoJob{}, errors.New("error parsing response")
	}

	return videoResponse, nil
}

func SaveURLToFile(ctx context.Context, url, outputPath string, headers map[string]string) error {
	// Reasonable HTTP timeouts
	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	// Build request with context
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// Do request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("http get: %w", err)
	}
	defer resp.Body.Close()

	// Check status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("bad status %d: %s", resp.StatusCode, string(b))
	}

	// Optional: sanity check content type (some servers omit or vary it)
	ct := strings.ToLower(resp.Header.Get("Content-Type"))
	if ct != "" && !strings.HasPrefix(ct, "video/") && !strings.Contains(ct, "mp4") {
		// Not fatal for all servers, but helpful to catch obvious mistakes.
		// Remove this block if your endpoint doesn't set Content-Type.
		return fmt.Errorf("unexpected Content-Type: %q", ct)
	}

	// Write atomically via temp file, then rename
	tmpPath := outputPath + ".part"
	out, err := os.Create(tmpPath)
	if err != nil {
		return fmt.Errorf("create tmp file: %w", err)
	}

	// Ensure file closed before rename
	writeErr := func() error {
		defer out.Close()
		// Stream copy (handles chunked & large files)
		if _, err := io.Copy(out, resp.Body); err != nil {
			return fmt.Errorf("write file: %w", err)
		}
		return nil
	}()
	if writeErr != nil {
		_ = os.Remove(tmpPath)
		return writeErr
	}

	// Finalize
	if err := os.Rename(tmpPath, outputPath); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("rename: %w", err)
	}

	return nil
}