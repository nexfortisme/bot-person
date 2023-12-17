package external

import (
	"fmt"
	"strings"
	"net/http"
	"io/ioutil"
)

func GetTTS(requestString string) {

	url := "https://api.elevenlabs.io/v1/text-to-speech/vRF4YMPsdFONAIl9yM3E"

	// 
	payload := strings.NewReader("{\n  \"text\": \"Some Days, I think I am a pirate.\",\n  \"model_id\": \"eleven_multilingual_v2\",\n  \"voice_settings\": {\n    \"stability\": 0.5,\n    \"similarity_boost\": 0.5\n  }\n}")

	req, _ := http.NewRequest("POST", url, payload)

	// TODO - Move API Key to config file
	req.Header.Add("xi-api-key", "a0db636977ac78d4a79a6c24198ecd06")
	req.Header.Add("Content-Type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))

}