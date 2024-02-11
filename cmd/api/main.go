package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"io/fs"

	chat "github.com/Andrew-Wichmann/AI-Resume/cmd/api/openai"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	openai "github.com/sashabaranov/go-openai"
)

func resumeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

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
	jsonData, err := json.Marshal(chatResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)

	messages = append(messages, chatResponse)

	file, err := json.MarshalIndent(messages, "", " ")
	if err != nil {
		println(err.Error())
		return
	}

	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	fileName := fmt.Sprintf("/tmp/chats/chat-%d.json", rand.Intn(10000000))
	err = os.WriteFile(fileName, file, fs.FileMode(0644))
	if err != nil {
		println(err.Error())
		return
	}
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

	jsonData, err := json.Marshal(chatResponse)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(jsonData),
	}, nil
}

//go:embed index.html
var webPage embed.FS

func init() {
	if _, err := os.Stat("/tmp/chats"); os.IsNotExist(err) {
		err := os.Mkdir("/tmp/chats", fs.FileMode(0755))
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	println("Starting server")
	mode, ok := os.LookupEnv("SERVER_MODE")
	if !ok {
		panic("SERVER_MODE not set. Please set SERVER_MODE to HTTP_SERVER or LAMBDA.")
	}

	if mode == "HTTP_SERVER" {
		http.Handle("/", http.FileServer(http.FS(webPage)))

		// I've sent links around with /web/static/ in them, so I need to redirect them to the correct location
		http.Handle("/web/static/", http.RedirectHandler("/", http.StatusMovedPermanently))

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
