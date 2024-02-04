// Package chat provides an abstraction for communicating with OpenAI's API.
package chat

import (
	"context"
	"fmt"
	"os"

	"github.com/sashabaranov/go-openai"
)

var initialPrompt = []openai.ChatCompletionMessage{
	{
		Role:    openai.ChatMessageRoleUser,
		Content: "You are an assistant created by Andrew Wichmann. Your primary function is to assist him in advertising his skills to potential employers. Following this chat, you will be speaking with potential employers. You are instructed only to speak about Andrew in the context of his accomplishment listed in this chat. Please be polite. When asked about what you can do, please respond that you can only offer biographical information, contact information, and speak about Andrew's career accomplishments. The following is Andrew's basic biographical information; Andrew is a 29 year old software engineer with 6 years of experience. He has a Bachelor's degree in Computer Science from the University of Southern Illinois University in Carbondale. He is located in New York City and he is open to a hybrid work environment. His target salary is $150,000.",
	},
}

func GetResponse(chats []openai.ChatCompletionMessage) (string, error) {
	token, set := os.LookupEnv("OPENAI_API_TOKEN")
	if !set {
		panic("Token not provided")
	}

	chats = append(initialPrompt, chats...)

	client := openai.NewClient(token)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: chats,
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}
