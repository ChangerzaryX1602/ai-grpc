package repository

import (
	"ai/api/pb"
	"ai/internal/entity"
	"context"
	"fmt"
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
			resp, err = generative.GenerateContent(ctx, genai.Text(prompt), genai.ImageData("jpeg", imageData))
			if err != nil {
				return nil, err
			}
		} else {
			resp, err = generative.GenerateContent(ctx, genai.Text(prompt))
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
	answer := ableToRead(resp.Answer)
	return &pb.AiResponse{
		Answer: answer,
	}, nil
}
