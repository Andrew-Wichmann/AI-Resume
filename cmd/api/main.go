package main

import (
	"encoding/json"
	"net/http"

	chat "github.com/Andrew-Wichmann/AI-Resume/cmd/api/openai"
	openai "github.com/sashabaranov/go-openai"
)

func resumeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	var messages []openai.ChatCompletionMessage

	err := json.NewDecoder(r.Body).Decode(&messages)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	chatResponse, err := chat.GetResponse(messages)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(chatResponse))
}

func main() {
	http.HandleFunc("/resume", resumeHandler)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	http.ListenAndServe(":8080", nil)
}
