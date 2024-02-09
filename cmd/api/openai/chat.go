// Package chat provides an abstraction for communicating with OpenAI's API.
package openai

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"os"

	"github.com/sashabaranov/go-openai"
)

var accomplishments []string

var primeDirective = `
You are an assistant created by Andrew Wichmann.
Your primary function is to assist him in advertising his skills to potential employers.
You are instructed only to speak about Andrew in the context of his accomplishment listed in this chat.
Please be polite.
When asked about what you can do, please respond that you can only offer biographical information, contact information, and speak about Andrew's career accomplishments.
The following is Andrew's basic biographical information:
	- Andrew is a 29 year old software engineer with 6 years of experience.
	- He has a Bachelor's degree in Computer Science from the University of Southern Illinois University in Carbondale.
	- He is located in New York City and he is open to a hybrid work environment.
	- His target salary is $175,000.
	- He does not require a visa sponsorship.
	- He is available to start immediately.
	- His email is AndrewP.Wichmann@gmail.com
	- His phone number is 312-221-8156
Following this chat will be a sequence of career accomplishments and initatives that Andrew has worked on.
You are instructed to inform potential employers that you know this information and that you can speak about it.
You will be informed when the sequence is complete when you are informed that you will be speaking with potential employers.
`

func GetResponse(chats []openai.ChatCompletionMessage) (openai.ChatCompletionMessage, error) {
	token, set := os.LookupEnv("OPENAI_API_TOKEN")
	if !set {
		panic("Token not provided")
	}

	// The PrimeDirectives are the initial description of the AI's purpose, plus a list of accomplishments, then an instruction
	// to start responding to potential employers.
	primeDirectives := make([]openai.ChatCompletionMessage, len(accomplishments)+2)
	primeDirectives[0] = openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: primeDirective,
	}
	for i := range accomplishments {
		primeDirectives[i+1] = openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: accomplishments[i],
		}
	}
	primeDirectives[len(primeDirectives)-1] = openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: "Following this chat, you will be speaking with potential employers.",
	}

	chats = append(primeDirectives, chats...)

	var model string
	if os.Getenv("ENV") == "production" {
		model = openai.GPT4
	} else {
		model = openai.GPT3Dot5Turbo
	}

	println("Using model: ", model)

	client := openai.NewClient(token)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    model,
			Messages: chats,
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return openai.ChatCompletionMessage{}, err
	}

	// I'm not sure why Choices is a slice, but I'm assuming it's because the AI can provide multiple responses.
	// I only want one response
	return resp.Choices[0].Message, nil
}

//go:embed accomplishments/*
var accomplishmentFiles embed.FS

func init() {
	files, err := fs.ReadDir(accomplishmentFiles, "accomplishments")
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if file.IsDir() {
			fmt.Printf("Warning: %s is a directory\n", file.Name())
		}
		data, err := fs.ReadFile(accomplishmentFiles, "accomplishments/"+file.Name())
		if err != nil {
			panic(err)
		}
		accomplishments = append(accomplishments, string(data))
	}

}
