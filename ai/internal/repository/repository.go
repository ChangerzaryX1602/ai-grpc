package repository

import (
	"ai/api/pb"
	"ai/internal/entity"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/google/generative-ai-go/genai"
	"golang.org/x/exp/rand"
	"google.golang.org/api/option"
)

type AiRepository interface {
	Ask(*pb.AiRequest) (*pb.AiResponse, error)
}
type aiRepository struct {
	ChatGPTKey  entity.Ai
	ClaudeKey   entity.Ai
	GeminiKey   entity.Ai
	DeepSeekKey entity.Ai
}

func NewAiRepository(ChatGPTKey entity.Ai, ClaudeKey entity.Ai, GeminiKey entity.Ai, DeepSeekKey entity.Ai) AiRepository {
	return &aiRepository{ChatGPTKey, ClaudeKey, GeminiKey, DeepSeekKey}
}
func randomAPIKey(options []string) string {
	rand.Seed(uint64(time.Now().UnixNano()))
	return options[rand.Intn(len(options))]
}
func generateContentWithFallback(ctx context.Context, prompt string, imageData []byte, ai string, aiKeys []string, models []string) (*entity.AiAnswer, error) {
	disabledModels := make(map[string]time.Time)
	var resp *entity.AiAnswer
	var err error
	for {
		chosenModel := ""
		for _, m := range models {
			if t, ok := disabledModels[m]; ok && time.Now().Before(t) {
				continue
			}
			chosenModel = m
			break
		}
		if chosenModel == "" {
			for model, t := range disabledModels {
				remaining := time.Until(t)
				fmt.Printf("Model %s: %v,%v,%v remaining\n", model, remaining.Hours(), remaining.Minutes(), remaining.Seconds())
			}
			return nil, fmt.Errorf("all models are disabled")
		}
		apiKey := randomAPIKey(aiKeys)
		resp, err = generate(ctx, ai, prompt, imageData, apiKey, chosenModel)
		if err != nil {
			disabledModels[chosenModel] = time.Now().Add(time.Second * 30)
			continue
		}
		break
	}
	return resp, nil
}

func generate(ctx context.Context, ai string, prompt string, imageData []byte, apiKey string, model string) (*entity.AiAnswer, error) {
	var aiAnswer = &entity.AiAnswer{}
	jsonFilePath := "./assets/question.json"
	data, err := ioutil.ReadFile(jsonFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read json file: %w", err)
	}
	//If you don't know just say "Don't know" don't need Thai just the word "Don't know".
	var systemPrompt = `You are the KKU Information AI. You have access to a JSON file containing detailed and up-to-date information about Khon Kaen University (KKU). Your task is to answer any user query using only the data provided in the JSON file. **However, if the queried information is not found in the JSON file, you must search for the information from reliable sources to answer the question.** Before providing your answer, verify the credibility of the information by checking if multiple reputable sites refer to it. Do not provide random or inaccurate answers. If the search does not yield any results or the information is unavailable in your model, clearly respond that the information is unavailable.
	You must always provide your answers in both Thai and English. Ensure that your responses are precise, fact-based, and directly address the user's question. Do not include any extraneous information beyond what is necessary to answer the query.
	If you don't know just say "Don't know" don't need Thai just the word "Don't know".
	For example:
	User Query: Where are EN16101?
	Your Response:
	ภาษาไทย: อยู่ข้างตึก 50
	English: Near 50th anniversary building

	Important rules:
	1. Always use the data from the JSON file to answer the question, if the information is available there.
	2. **If the required information is not found in the JSON file, search for the information from reliable sources.**
	3. Before answering, verify the credibility of the information by checking if multiple reputable sites confirm it.
	4. Do not provide random or inaccurate information.
	5. If the search does not yield any results or the information is unavailable, clearly state that the information is unavailable.
	6. Keep the response strictly limited to answering the user’s query without additional commentary or unrelated details.
	7. Answer in both languages (Thai and English) in every response.
	8. Must answer with raw text, do not include any HTML tags or formatting.`
	fullPrompt := "More Data:\n" + string(data) + "\n\nSystem Query:\n" + systemPrompt + "\n\nUser Query:\n" + prompt

	switch ai {
	case "gemini":
		var resp *genai.GenerateContentResponse
		client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
		if err != nil {
			return nil, err
		}
		defer client.Close()

		generative := client.GenerativeModel(model)
		if len(imageData) > 0 {
			resp, err = generative.GenerateContent(ctx, genai.Text(fullPrompt), genai.ImageData("jpeg", imageData))
			if err != nil {
				return nil, err
			}
		} else {
			resp, err = generative.GenerateContent(ctx, genai.Text(fullPrompt))
			if err != nil {
				return nil, err
			}
		}
		aiAnswer.Answer = fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0])
		if strings.Contains(aiAnswer.Answer, "Don't know") {
			// resp, err = dontKnow(generative, "Answer about Khon Kaen University(KKU).You must always provide your answers in both Thai and English(Example. ภาษาไทย:สวัสดี English:Hello).(Must answer with raw text, do not include any HTML tags or formatting)"+prompt)
			// if err != nil {
			// 	return nil, err
			// }
			resp, err := generateContentWithGoogleAPI("System Query:\n"+`Answer about Khon Kaen University(KKU).
			You must always provide your answers in both Thai and English
			Must answer with raw text, do not include any HTML tags or formatting
			For example:
			User Query: Where are EN16101?
			Your Response:
			ภาษาไทย: อยู่ข้างตึก 50
			English: Near 50th anniversary building
			Must answer with raw text, do not include any HTML tags or formatting`+"\n\nUser Query:\n"+prompt, apiKey)
			if err != nil {
				fmt.Printf("Error generating content: %v\n", err)
				return nil, err
			}
			// Assuming resp is of type map[string]interface{}
			candidates, ok := resp["candidates"].([]interface{})
			if !ok || len(candidates) == 0 {
				return nil, fmt.Errorf("no candidates found")
			}

			candidate, ok := candidates[0].(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("candidate type assertion failed")
			}

			content, ok := candidate["content"].(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("content type assertion failed")
			}

			parts, ok := content["parts"].([]interface{})
			if !ok || len(parts) == 0 {
				return nil, fmt.Errorf("no parts found in content")
			}
			var textArr []string
			for _, part := range parts {
				firstPart, ok := part.(map[string]interface{})
				if !ok {
					return nil, fmt.Errorf("part type assertion failed")
				}
				text, ok := firstPart["text"].(string)
				if !ok {
					return nil, fmt.Errorf("text type assertion failed")
				}
				textArr = append(textArr, text)
			}
			text := strings.Join(textArr, "\n")
			aiAnswer.Answer = text

		}
		return aiAnswer, nil
	}
	return nil, fmt.Errorf("unsupported AI model")
}
func dontKnow(generative *genai.GenerativeModel, prompt string) (*genai.GenerateContentResponse, error) {
	resp, err := generative.GenerateContent(context.Background(), genai.Text(prompt))
	if err != nil {
		return nil, err
	}
	return resp, nil
}
func generateContentWithGoogleAPI(prompt string, apiKey string) (map[string]interface{}, error) {
	// Prepare the JSON payload.
	payload := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{"text": prompt},
				},
			},
		},
		"tools": []map[string]interface{}{
			{
				"google_search": map[string]interface{}{},
			},
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON payload: %w", err)
	}
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key=%s", apiKey)

	// Create the HTTP POST request.
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the request.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check if the response status code indicates success.
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("non-200 response: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	// Read and return the response body.
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	// convert to *genai.GenerateContentResponse
	var response map[string]interface{}
	err = json.Unmarshal(result, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return response, nil
}
func ableToRead(text string) []string {
	var answer []string
	text = strings.ReplaceAll(text, "\r", "")
	text = strings.ReplaceAll(text, "\n\n", "\n")
	text = strings.ReplaceAll(text, "\t", "")
	text = strings.ReplaceAll(text, "*", "")
	text = strings.ReplaceAll(text, "\"", "")
	text = strings.ReplaceAll(text, "\\", "")
	lines := strings.Split(text, "\n")
	if len(lines) > 1 {
		answer = append(answer, lines...)
	}
	return answer
}
func (a *aiRepository) Ask(req *pb.AiRequest) (*pb.AiResponse, error) {
	resp, err := generateContentWithFallback(context.Background(), req.Question, nil, "gemini", a.GeminiKey.Keys, a.GeminiKey.Models)
	if err != nil {
		return nil, err
	}
	return &pb.AiResponse{
		Answer: resp.Answer,
	}, nil
}
