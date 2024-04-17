package commands

import (
	"fmt"
	"main/pkg/external"
	"main/pkg/persistance"

	"main/pkg/logging"
	eventType "main/pkg/logging/enums"

	"github.com/bwmarrin/discordgo"
)

func Image(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Access options in the order provided by the user.
	userImageOptions := i.ApplicationCommandData().Options

	// Or convert the slice into a map
	userImageOptionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(userImageOptions))
	for _, opt := range userImageOptions {
		userImageOptionMap[opt.Name] = opt
	}

	user, _ := persistance.GetUser(i.Interaction.Member.User.ID)

	if user.UserStats.ImageTokens < 10 {

		logging.LogEvent(eventType.COMMAND_IMAGE, i.Interaction.Member.User.ID, "NOT ENOUGH TOKENS", i.Interaction.GuildID)

		// Getting user stat data
		imageReturnString := "You don't have enough tokens (10) to generate an image."

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: imageReturnString,
			},
		})
		return
	}

	// Pulling the propt out of the optionsMap
	if option, ok := userImageOptionMap["prompt"]; ok {

		// Generating the response
		placeholder := "Prompt: " + option.StringValue()

		// Immediately responding in the 3 second window before the interaciton times out
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: placeholder,
			},
		})

		// Going out to make the OpenAI call to get the proper response
		returnFile, err := ParseDalleRequest(s, option.StringValue())

		if err != nil {

			errString := fmt.Sprintf("Something Went Wrong: %s", err.Error())

			// Not 100% sure this is the approach I want to take with handling errors from the API
			_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content: &errString,
			})

			if err != nil {
				s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
					Content: "Something went wrong. Send help.",
				})
			}

			return
		}

		user.UserStats.ImageTokens -= 10
		persistance.UpdateUser(*user)

		logging.LogEvent(eventType.COMMAND_IMAGE, i.Interaction.Member.User.ID, option.StringValue(), i.Interaction.GuildID)

		// Updating the initial message with the response from the OpenAI API
		_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Files: []*discordgo.File{&returnFile},
		})

		if err != nil {
			// Not 100% sure this is the approach I want to take with handling errors from the API
			s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
				Content: "Something went oopsie with sending the file.",
			})
			return
		}
	}
}

func ParseDalleRequest(s *discordgo.Session, prompt string) (discordgo.File, error) {
	dalleResponse, err := external.GetDalleResponse(prompt)

	if err != nil {
		return discordgo.File{}, err
	}

	return dalleResponse, nil
}