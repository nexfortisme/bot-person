package commands

import (
	"context"
	"fmt"
	"main/pkg/external"
	"main/pkg/persistance"
	"main/pkg/util"
	"os"
	"path/filepath"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Slop struct{}

func (b *Slop) ApplicationCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "slop",
		Description: "Generate a 4 second Sora 2 Video.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "prompt",
				Description: "The script for the slop.",
				Required:    true,
			},
		},
	}
}

func (b *Slop) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) {

	// Access options in the order provided by the user.
	options := i.ApplicationCommandData().Options

	// Or convert the slice into a map
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	user, _ := persistance.GetUser(i.Interaction.Member.User.ID)

	if user.ImageTokens < b.CommandCost() {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You don't have enough tokens to generate a slop.",
			},
		})
		return
	}

	// Pulling the propt out of the optionsMap
	if option, ok := optionMap["prompt"]; ok {

		// Generating the response
		placeholderBotResponse := "Pooping out slop with prompt: " + option.StringValue()

		// Immediately responding in the 3 second window before the interaciton times out
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: placeholderBotResponse,
			},
		})

		user.ImageTokens -= b.CommandCost()
		persistance.UpdateUser(*user)

		// Going out to make the OpenAI call to get the proper response
		// botResponseString = ParseGPTSlashCommand(s, option.StringValue(
		// Check if the response will be too long and truncate if necessary
		prompt := option.StringValue()
		go handleAsyncSlop(prompt, i, s)
	}
}

func (b *Slop) HelpString() string {
	return "The `/slop` command allows you to generate a 4 second Sora 2 Video."
}

func (b *Slop) CommandCost() int {
	return 8 // 5 Cents a Token
}

func handleAsyncSlop(prompt string, i *discordgo.InteractionCreate, s *discordgo.Session) {
	videoResponse, err := external.GetSoraRespone(prompt)
	if err != nil {
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "Error getting Sora response: " + err.Error(),
		})
		return
	}

	// Polling for the video status
	ctx := context.Background()
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Context done")
			return // Its Done I guess
		case <-ticker.C:

			videoResponse, err = external.GetSoraJobStatus(videoResponse.ID)

			fmt.Printf("Video Response: %+v\n", videoResponse)

			if err != nil {
				s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
					Content: "Error getting Sora job status: " + err.Error(),
				})
				return
			}

			if videoResponse.Status == "completed" {

				url := "https://api.openai.com/v1/videos/" + videoResponse.ID + "/content"
				headers := map[string]string{"Authorization": "Bearer " + util.GetOpenAIKey()}

				videoName := fmt.Sprintf("%s.mp4", videoResponse.ID)

				err = util.CreateDirectoryIfNotExists("videos")
				if err != nil {
					s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
						Content: "Error creating videos directory: " + err.Error(),
					})
					return
				}
				path := filepath.Join("videos", videoName)

				err = external.SaveURLToFile(ctx, url, path, headers)
				if err != nil {
					s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
						Content: "Error saving video to file: " + err.Error(),
					})
					return
				}

				reader, err := os.Open(filepath.Join("videos", fmt.Sprintf("%s.mp4", videoResponse.ID)))
				if err != nil {
					s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
						Content: "Error opening video file: " + err.Error(),
					})
					return
				}

				fileInfo, err := reader.Stat()
				if err != nil {
					s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
						Content: "Error getting video file info: " + err.Error(),
					})
					return
				}

				fileObj := &discordgo.File{
					Name:        fileInfo.Name(),
					ContentType: "video/mp4",
					Reader:      reader,
				}

				slopsReadyString := fmt.Sprintf("Here is your %s. Slop is ready! %s", prompt, videoResponse.ID)

				_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
					Content: &slopsReadyString,
					Files:   []*discordgo.File{fileObj},
				})
				if err != nil {
					// Not 100% sure this is the approach I want to take with handling errors from the API
					s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
						Content: "Something went oopsie with sending the file.",
					})
					return
				}
				return
			} else if videoResponse.Status == "failed" {

				failureString := fmt.Sprintf("Video Generation Failed: %s", videoResponse.Error.Message)

				s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
					Content: &failureString,
				})

				user, _ := persistance.GetUser(i.Interaction.Member.User.ID)
				user.ImageTokens += 8
				persistance.UpdateUser(*user)

				return
			} else {
				currentStatusString := fmt.Sprintf("Video Status: %s. Progress: %d%%", videoResponse.Status, videoResponse.Progress)
				s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
					Content: &currentStatusString,
				})
			}
		}
	}
}
