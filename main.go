package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/nsf/termbox-go"
)

const SPACE = " "
const BACKSPACE = "\x08"
const OPENAI_MODEL = "gpt-3.5-turbo-0125"
const OPENAI_URL = "https://api.openai.com/v1/chat/completions"
const OPENAI_MODEL_TEMP = 0.5

func main() {
	godotenv.Load()

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("no OPENAI_API_KEY environment variable found")
	}

	client := &http.Client{}

	err := termbox.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer termbox.Close()

	fmt.Printf("using %s model at %s\n", OPENAI_MODEL, OPENAI_URL)

	var input strings.Builder

	fmt.Print(">")

main_loop:
	for {
		switch keyPress := termbox.PollEvent(); keyPress.Type {
		case termbox.EventKey:
			switch keyPress.Key {
			case termbox.KeyEsc:
				break main_loop
			case termbox.KeyEnter:
				message := strconv.Quote(input.String())
				resp, err := postGptRequest(client, message, apiKey)
				if err != nil {
					log.Printf("\nError making request to OpenAI servers: %s", err.Error())
				}
				fmt.Printf("\n\tchat gipity: %v\n>", resp)
				input.Reset()
				break
			case termbox.KeyBackspace, termbox.KeyBackspace2:
				input.WriteString(BACKSPACE)
				fmt.Print(BACKSPACE)
			case termbox.KeySpace:
				input.WriteString(SPACE)
				fmt.Print(SPACE)
			default:
				if keyPress.Ch != 0 {
					char := fmt.Sprintf("%c", keyPress.Ch)
					input.WriteString(char)
					fmt.Print(char)
				}
			}
		}
	}
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIRequest struct {
	Model       string    `json:"model"`
	Messages    []message `json:"messages"`
	Temperature float32   `json:"temperature"`
}

type openAIResponse struct {
	ID                string `json:"id"`
	Object            string `json:"object"`
	Created           int32  `json:"created"`
	Model             string `json:"model"`
	SystemFingerprint string `json:"system_fingerprint"`
	Choices           []struct {
		Index        int32   `json:"index"`
		Message      message `json:"message"`
		FinishReason string  `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int32 `json:"prompt_tokens"`
		CompletionTokens int32 `json:"completion_tokens"`
		TotalTokens      int32 `json:"total_tokens"`
	}
}

func postGptRequest(client *http.Client, input string, apiKey string) (*openAIResponse, error) {
	openAIRequest := openAIRequest{
		Model:       OPENAI_MODEL,
		Messages:    []message{{Role: "", Content: input}},
		Temperature: OPENAI_MODEL_TEMP,
	}
	payload, err := json.Marshal(openAIRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal json payload %s", err.Error())
	}

	request, err := http.NewRequest("POST", OPENAI_URL, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create http request %s", err.Error())
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to send http request %s", err.Error())
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status code from server %d", response.StatusCode)
	}

	var body openAIResponse
	err = json.NewDecoder(response.Body).Decode(&body)
	if err != nil {
		return nil, fmt.Errorf("failed to decode the http response %s", err.Error())
	}

	response.Body.Close()

	return &body, nil

}
