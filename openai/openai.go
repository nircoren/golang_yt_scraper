package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

const apiURL = "https://api.openai.com/v1/chat/completions"

type Exercise struct {
	ExerciseName        ExerciseName `json:"exerciseName"`
	ExerciseDescription string       `json:"exerciseDescription"`
	YoutubeVideoId      string       `json:"youtubeVideoId"`
}

type ExerciseName struct {
	En string `json:"en"`
	He string `json:"he"`
}

type ExerciseProgram struct {
	Exercises []Exercise `json:"exercises"`
	Injury    string     `json:"injury"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type RequestBody struct {
	Messages         []Message `json:"messages"`
	Temperature      float64   `json:"temperature"`
	MaxTokens        int       `json:"max_tokens"`
	TopP             float64   `json:"top_p"`
	FrequencyPenalty float64   `json:"frequency_penalty"`
	PresencePenalty  float64   `json:"presence_penalty"`
	Model            string    `json:"model"`
}

type Choice struct {
	Message Message `json:"message"`
}

type ResponseBody struct {
	Choices []Choice `json:"choices"`
}

func GenerateExercises(system string, user string) (*ExerciseProgram, error) {
	err := godotenv.Load()
	if err != nil {

	}
	apiKey := os.Getenv("OPENAI_API_KEY1")
	if apiKey == "" {
		return nil, fmt.Errorf("failed to load env key", err)

	}

	messages := []Message{
		{Role: "system", Content: system},
		{Role: "user", Content: user},
	}

	requestBody := RequestBody{
		Messages:         messages,
		Temperature:      0,
		MaxTokens:        700,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
		Model:            "gpt-4o",
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create new request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response status: %v", resp.Status)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var responseBody ResponseBody
	err = json.Unmarshal(respBody, &responseBody)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %v", err)
	}

	if len(responseBody.Choices) > 0 {
		content := responseBody.Choices[0].Message.Content
		cleanedJSON := CleanJSON(content)
		result, err := UnmarshalExercises(cleanedJSON)
		if err != nil {
			fmt.Println("Error:", err)
			return nil, fmt.Errorf("failed to unmarshal exercises: %v", err)
		}
		// Assuming the content is a JSON array of exercise names
		return result, nil
	}

	return nil, fmt.Errorf("no choices returned in response")
}

func CleanJSON(input string) string {
	// Remove \n, ```, and "json"
	cleaned := strings.ReplaceAll(input, "\n", "")
	cleaned = strings.ReplaceAll(cleaned, "```", "")
	cleaned = strings.ReplaceAll(cleaned, "json", "")

	// Trim any leading or trailing whitespace
	cleaned = strings.TrimSpace(cleaned)

	return cleaned
}

func UnmarshalExercises(jsonString string) (*ExerciseProgram, error) {
	var program ExerciseProgram
	err := json.Unmarshal([]byte(jsonString), &program)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal exercises: %v", err)
	}
	return &program, nil
}

func UnmarshalExercises2(jsonString string) {

	var rawMap map[string]interface{}
	err := json.Unmarshal([]byte(jsonString), &rawMap)
	if err != nil {
		fmt.Printf("Error unmarshaling to map: %v\n", err)
	} else {
		fmt.Printf("Raw structure: %+v\n", rawMap)
	}
}
