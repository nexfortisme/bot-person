package persistance

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

var (
	botTracking  BotTracking
	tempTracking BotTracking
)

func ReadBotStatistics() {

	log.Println("Reading botTracking.json...")

	var trackingFile []byte
	trackingFile, err := ioutil.ReadFile("botTracking.json")
	if err != nil {

		log.Printf("Error Reading botTracking. Creating File")
		os.WriteFile("botTracking.json", []byte("{\"BadBotCount\":0,\"MessageCount\":0, \"UserStats\":[]}"), 0666)

		trackingFile, err = ioutil.ReadFile("botTracking.json")
		if err != nil {
			log.Fatalf("Could not read config file: botTracking.json")
		}

	}

	err = json.Unmarshal(trackingFile, &botTracking)
	if err != nil {
		log.Fatalf("Could not parse: botTracking.json")
	}

	log.Println("Done Reading botTracking.json")
}

func quickReadStats() {
	var trackingFile []byte
	trackingFile, err := ioutil.ReadFile("botTracking.json")
	err = json.Unmarshal(trackingFile, &tempTracking)
	if err != nil {
		log.Fatalf("Could not parse: botTracking.json")
	}
}

func SaveBotStatistics() {
	log.Println("Writing botTracking.json...")
	fle, _ := json.Marshal(botTracking)
	os.WriteFile("botTracking.json", fle, 0666)
	log.Println("Done Writing botTracking.json")
}

func GetBotTracking() BotTracking {
	return botTracking
}

func GetTempTracking() BotTracking {
	quickReadStats()
	return tempTracking
}
