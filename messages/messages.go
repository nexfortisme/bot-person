package messages

import (

	// "encoding/base64"
	"bytes"
	"image"
	"image/png"

	// "encoding/base64"
	"encoding/base64"
	"encoding/json"
	"fmt"

	// "image/png"
	"os"
	"strconv"

	// "image/png"

	// "io"
	"io/ioutil"
	"log"
	"main/logging"
	"main/util"
	"math/rand"
	"net/http"

	// "os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func ParseMessage(s *discordgo.Session, m *discordgo.MessageCreate, openAIKey string) {

	// Ignoring messages from self
	if m.Author.ID == s.State.User.ID {
		return
	}

	var incomingMessage string
	badBotResponses := make([]string, 0)
	badBotResponses = append(badBotResponses, "I'm sorry")
	badBotResponses = append(badBotResponses, "It won't happen again")
	badBotResponses = append(badBotResponses, "Eat Shit")
	badBotResponses = append(badBotResponses, "Ok.")
	badBotResponses = append(badBotResponses, "Sure Thing.")
	badBotResponses = append(badBotResponses, "Like you are the most perfect being in existance. Pound sand pal.")

	if !strings.HasPrefix(m.Message.Content, "!") {
		incomingMessage = strings.ToLower(m.Message.Content)
	} else {
		incomingMessage = m.Message.Content
	}

	// TODO - Handle this better. I don't like this and I feel bad about it
	if strings.HasPrefix(incomingMessage, "bad bot") {
		logging.LogIncomingMessage(s, m)

		logging.IncrementTracker(2, m.Author.ID, m.Author.Username)
		badBotRetort := badBotResponses[rand.Intn(len(badBotResponses))]
		// TODO - Here Too
		log.Println("Bot Person > " + badBotRetort)
		_, err := s.ChannelMessageSend(m.ChannelID, badBotRetort)
		util.HandleErrors(err)
	} else if strings.HasPrefix(incomingMessage, "good bot") {
		logging.LogIncomingMessage(s, m)

		logging.IncrementTracker(1, m.Author.ID, m.Author.Username)
		log.Println("Bot Person > Thank You!")
		_, err := s.ChannelMessageSend(m.ChannelID, "Thank You!")
		util.HandleErrors(err)
	} else if strings.HasPrefix(incomingMessage, "!botStats") {
		logging.LogIncomingMessage(s, m)

		logging.GetBotStats(s, m)
		logging.IncrementTracker(0, m.Author.ID, m.Author.Username)
	} else if strings.HasPrefix(incomingMessage, "!myStats") {
		logging.LogIncomingMessage(s, m)

		logging.GetUserStats(s, m)
		logging.IncrementTracker(0, m.Author.ID, m.Author.Username)
	}

	// Commands to add
	// about - list who made it and maybe a link to the git repo
	// invite - Generates an invite link to be able to invite the bot to differnet servers
	// stopTracking - Allows uers to opt out of data collection

	// Only process messages that mention the bot
	id := s.State.User.ID
	if !mentionsBot(m.Mentions, id) && !mentionsKeyphrase(m) {
		return
	}

	msg := util.ReplaceIDsWithNames(m, s)

	logging.LogIncomingMessage(s, m)

	logging.IncrementTracker(0, m.Author.ID, m.Author.Username)
	rsp := getSDResponse(msg)
	if rsp == "" {
		return
	}

	fmt.Println(rsp)

	outputName := msg + ".png"

	tmp := rsp[22:]
	fmt.Println("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	fmt.Println(tmp)

	unbased, err := base64.StdEncoding.DecodeString(tmp)
	if err != nil {
		panic("Cannot decode b64")
	}

	r := bytes.NewReader(unbased)
	im, err := png.Decode(r)
	if err != nil {
		panic("Bad png")
	}

	f, err := os.OpenFile("example.png", os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		panic("Cannot open file")
	}

	png.Encode(f, im)

	// f, err := os.Create(outputName);

	// if err != nil {
	// 	logging.LogError("Unable to Create File");
	// }

	// defer f.Close()
	// f.Write([]byte(rsp));

	// bar := []byte(rsp)

	exp, err := os.OpenFile(outputName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		logging.LogError("Unable to Create File")
	}
	baz, _, err := image.Decode(strings.NewReader(rsp))
	if err != nil {
		fmt.Print(err)
		logging.LogError("Unable to Create File")
	}
	png.Encode(exp, baz)

	fle, err := ioutil.ReadFile(outputName)

	s.ChannelFileSend(m.ChannelID, outputName, bytes.NewReader(fle))

	// respTxt := getOpenAIResponse(msg, openAIKey)

	// // TODO - Here as well
	// log.Printf("Bot Person > %s \n", respTxt)
	// if mentionsKeyphrase(m) {
	// 	s.ChannelMessageSend(m.ChannelID, "!bot is deprecated. Please at the bot or use /bot for further interactions")
	// }
	// _, err := s.ChannelMessageSend(m.ChannelID, respTxt)
	// util.HandleErrors(err)

}

// TODO - Make the response that is being logged by the bot include the bot user's actual username instead of "Bot Person"
func ParseSlashCommand(s *discordgo.Session, prompt string, openAIKey string) string {
	respTxt := getOpenAIResponse(prompt, openAIKey)
	respTxt = "Request: " + prompt + " " + respTxt
	log.Printf("Bot Person > %s \n", respTxt)
	return respTxt
}

func getOpenAIResponse(prompt string, openAIKey string) string {
	client := &http.Client{}

	dataTemplate := `{
		"model": "text-davinci-002",
		"prompt": "%s",
		"temperature": 0.7,
		"max_tokens": 256,
		"top_p": 1,
		"frequency_penalty": 0,
		"presence_penalty": 0
	  }`
	data := fmt.Sprintf(dataTemplate, prompt)

	req, err := http.NewRequest(http.MethodPost, "https://api.openai.com/v1/completions", strings.NewReader(data))
	if err != nil {
		logging.LogError("Error creating POST request")
		// log.Fatalf("Error creating POST request")
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+openAIKey)

	resp, _ := client.Do(req)

	if resp == nil {
		return "Error Contacting OpenAI API. Please Try Again Later."
	}

	buf, _ := ioutil.ReadAll(resp.Body)
	var rspOAI OpenAIResponse
	// TODO: This could contain an error from OpenAI (rate limit, server issue, etc)
	// need to add proper error handling
	err = json.Unmarshal([]byte(string(buf)), &rspOAI)
	if err != nil {
		return ""
	}

	// It's possible that OpenAI returns no response, so
	// fallback to a default one
	if len(rspOAI.Choices) == 0 {
		return "I'm sorry, I don't understand?"
	} else {
		return rspOAI.Choices[0].Text
	}
}

func getSDResponse(prompt string) string {
	seed := rand.Int31n(10000000)
	fmt.Println("Prompt: " + prompt + ". Seed: " + strconv.Itoa(int(seed)))

	client := &http.Client{}

	dataTemplate := `{
		"guidance_scale": "7.5",
		"height": "512",
		"negative_prompt": "",
		"num_inference_steps": 50,
		"num_outputs": "1",
		"prompt": "A Bear",
		"sampler": "plms",
		"seed": %d,
		"session_id": %d,
		"show_only_filtered_image": true,
		"stream_image_progress": false,
		"stream_progress_updates": false,
		"turbo": true,
		"use_cpu": false,
		"use_full_precision": true,
		"width": "512"
	}`

	data := fmt.Sprintf(dataTemplate, int(seed), time.Now().Unix())

	req, err := http.NewRequest(http.MethodPost, "http://localhost:9000/image", strings.NewReader(data))
	if err != nil {
		logging.LogError("Error Creating POST Request")
	}

	req.Header.Add("Content-Type", "application/json")

	resp, _ := client.Do(req)
	if resp == nil {
		fmt.Println("Null Response From SD Backend")
		return ""
	}

	buf, _ := ioutil.ReadAll(resp.Body)
	sdRsp := SDResponse{}

	err = json.Unmarshal([]byte(string(buf)), &sdRsp)
	if err != nil {
		logging.LogError("Error Unmarshalling response from SD Backend")
	}

	return sdRsp.Output[0].Data
}

func mentionsKeyphrase(m *discordgo.MessageCreate) bool {
	return strings.HasPrefix(m.Content, "!bot") && m.Content != "!botStats"
}

// Determine if the bot's ID is in the list of users mentioned
func mentionsBot(mentions []*discordgo.User, id string) bool {
	for _, user := range mentions {
		if user.ID == id {
			return true
		}
	}
	return false
}
