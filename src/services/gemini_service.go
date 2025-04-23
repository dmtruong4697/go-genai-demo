package services

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/go-resty/resty/v2"
)

type GeminiService struct {
	APIKey string
}

func NewGeminiService() *GeminiService {
	return &GeminiService{
		APIKey: os.Getenv("GEMINI_API_KEY"),
	}
}

func (s *GeminiService) Chat(prompt string) (string, error) {
	client := resty.New()

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemma-3-27b-it:generateContent?key=%s", s.APIKey)

	body := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"role": "user",
				"parts": []map[string]string{
					{"text": prompt},
				},
			},
		},
	}

	var result map[string]interface{}

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		SetResult(&result).
		Post(url)

	if err != nil {
		return "", err
	}

	if resp.IsError() {
		return "", fmt.Errorf("API error: %s", resp.String())
	}

	candidates := result["candidates"].([]interface{})
	firstCandidate := candidates[0].(map[string]interface{})
	content := firstCandidate["content"].(map[string]interface{})
	parts := content["parts"].([]interface{})
	text := parts[0].(map[string]interface{})["text"].(string)

	return text, nil
}

func (s *GeminiService) StreamChat(prompt string, writer http.ResponseWriter) error {
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1/models/gemma-3-27b-it:streamGenerateContent?key=%s", s.APIKey)

	body := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"role": "user",
				"parts": []map[string]string{
					{"text": prompt},
				},
			},
		},
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	writer.Header().Set("Content-Type", "text/event-stream")
	writer.Header().Set("Cache-Control", "no-cache")
	writer.Header().Set("Connection", "keep-alive")

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")
			fmt.Fprintf(writer, "data: %s\n\n", data)
			writer.(http.Flusher).Flush()
		}
	}
	return nil
}

func (s *GeminiService) ListModels() (interface{}, error) {
	client := resty.New()

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models?key=%s", s.APIKey)

	var result map[string]interface{}

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetResult(&result).
		Get(url)

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.String())
	}

	return result["models"], nil
}