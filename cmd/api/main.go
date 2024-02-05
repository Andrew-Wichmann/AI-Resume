package main

import (
	"embed"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	chat "github.com/Andrew-Wichmann/AI-Resume/cmd/api/openai"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	openai "github.com/sashabaranov/go-openai"
)

func resumeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	var messages []openai.ChatCompletionMessage

	err := json.NewDecoder(r.Body).Decode(&messages)
	if err != nil {
		println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	chatResponse, err := chat.GetResponse(messages)
	if err != nil {
		println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(chatResponse))
}

func lambdaHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if request.Path == "/resume" {

	} else if request.Path == "/health" {
	} else if request.Path == "/" {
		webPage, err := webPage.ReadFile("web/static/index.html")
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       err.Error(),
			}, nil
		}
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       string(webPage),
			Headers: map[string]string{
				"Content-Type": "text/html",
			},
		}, nil
	}
	var messages []openai.ChatCompletionMessage
	err := json.NewDecoder(strings.NewReader(request.Body)).Decode(&messages)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}
	chatResponse, err := chat.GetResponse(messages)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       chatResponse,
	}, nil
}

//go:embed web/static/index.html
var webPage embed.FS

func main() {
	println("Starting server")
	mode, ok := os.LookupEnv("SERVER_MODE")
	if !ok {
		panic("SERVER_MODE not set")
	}

	if mode == "HTTP_SERVER" {
		http.Handle("/", http.FileServer(http.FS(webPage)))
		http.HandleFunc("/resume", resumeHandler)
		http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})
		http.ListenAndServe(":8080", nil)

	} else if mode == "LAMBDA" {
		lambda.Start(lambdaHandler)
	} else {
		panic("Invalid SERVER_MODE")
	}

}
