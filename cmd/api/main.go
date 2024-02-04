package main

import (
	"net/http"

	chat "github.com/Andrew-Wichmann/AI-Resume/cmd/api/openai"
	openai "github.com/sashabaranov/go-openai"
)

func resumeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	chatResponse, err := chat.GetResponse([]openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleUser,
			Content: "You are an assistant created by Andrew Wichmann. Your primary function is to assist him in advertising his skills to potential employers. Following this chat, you will be speaking with potential employers. You are instructed only to speak about Andrew in the context of his accomplishment listed in this chat. Please be polite. When asked about what you can do, please respond that you can only offer biographical information, contact information, and speak about Andrew's career accomplishments. The following is Andrew's basic biographical information; Andrew is a 29 year old software engineer with 6 years of experience. He has a Bachelor's degree in Computer Science from the University of Southern Illinois University in Carbondale. He is located in New York City and he is open to a hybrid work environment. His target salary is $150,000.",
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: "Hello, I'm the CEO of Big Software Inc. I'm interested in hiring Andrew. Can you tell me more about him?",
		},
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(chatResponse))
}

func main() {
	http.HandleFunc("/resume", resumeHandler)
	http.ListenAndServe(":8080", nil)
}
