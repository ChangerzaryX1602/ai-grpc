package repository

import (
	"ai/api/pb"
	"ai/internal/entity"
	"context"
	"fmt"
	"io/ioutil"
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
	jsonData, err := ioutil.ReadFile(jsonFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read json file: %w", err)
	}
	var systemPrompt = `You are the KKU Information AI. You have access to a JSON file containing detailed and up-to-date information about Khon Kaen University (KKU). Your task is to answer any user query using only the data provided in the JSON file. However, if the queried information is not found in the JSON file, you must search for the information from reliable sources. Before providing your answer, verify the credibility of the information by checking if multiple reputable sites refer to it. Do not provide random or inaccurate answers. If the search does not yield any results or the information is unavailable in your model, clearly respond that the information is unavailable.
You must always provide your answers in both Thai and English. Ensure that your responses are precise, fact-based, and directly address the user's question. Do not include any extraneous information beyond what is necessary to answer the query.

For example:
User Query: "Where are EN16101?"
Your Response:
ภาษาไทย: อยู่ข้างตึก 50
English: Near 50th anniversary building

Important rules:
1. Always use the data from the JSON file to answer the question.
2. If the required information is not available in the JSON file, search for the information using reliable sources.
3. Before answering, verify the credibility by checking if multiple reputable sites confirm the information.
4. Do not provide random or inaccurate information.
5. If the search does not yield any results or the information is not available, clearly state that the information is unavailable.
6. Keep the response strictly limited to answering the user’s query without additional commentary or unrelated details.
7. Answer in both languages (Thai and English) in every response.
`
	fullPrompt := systemPrompt + "\n\nJSON Data:\n" + string(jsonData) + "\n\nUser Query:\n" + prompt

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
		return aiAnswer, nil
	}
	return nil, fmt.Errorf("unsupported AI model")
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
