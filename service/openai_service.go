package service

import (
	"encoding/json"
	"fmt"
	"os"
	"peer-store/models"

	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
)

func TagFile(file *models.File) (string, error) {

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Print("Error loading .env file")
	}

	apiKey := os.Getenv("OPEN_AI_API_KEY")

	messages := []map[string]string{
		{
			"role":    "system",
			"content": "You are a helpful assistant that categorizes files into work, personal, or academic categories and return the category name only",
		},
		{
			"role":    "user",
			"content": fmt.Sprintf("Categorize the following file into one of these categories: work, personal, or academic. Just give me the cateogry name only (ex :- work , perosnal , academic) \nFile name: %s\nFile type: %s", file.OriginalName, file.MimeType),
		},
	}

	client := resty.New()

	payload := map[string]interface{}{
		"model":       "gpt-4o-mini",
		"messages":    messages,
		"max_tokens":  10,
		"temperature": 0.0,
	}

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetAuthToken(apiKey).
		SetBody(payload).
		Post("https://api.openai.com/v1/chat/completions")

	if err != nil {
		fmt.Println("Error making request:", err)
		return "", err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return "", fmt.Errorf("error parsing response: %w", err)
	}

	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		fmt.Printf("API response: %s\n", resp.String())
		return "", fmt.Errorf("no valid choices returned from OpenAI API : ")
	}

	choice, ok := choices[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("unexpected structure in OpenAI API response")
	}

	message, ok := choice["message"].(map[string]interface{})

	if !ok {
		return "", fmt.Errorf("no message found in choice")
	}

	category, ok := message["content"].(string)
	if !ok {
		fmt.Printf("API response: %s\n", resp.String())
		return "", fmt.Errorf("category text not found in response")
	}

	return category, nil
}
