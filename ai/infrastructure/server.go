package infrastructure

import (
	"ai/api/pb"
	"ai/api/server"
	"ai/internal/entity"
	"ai/internal/repository"
	"ai/internal/usecase"
	"fmt"
	"log"

	"net"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type Resources struct {
	*gorm.DB
	gemini   entity.Ai
	chatgpt  entity.Ai
	claude   entity.Ai
	deepSeek entity.Ai
}

func NewServer(version, buildTag, runEnv string) (servers *Resources, err error) {
	var resources Resources
	if geminiKey := getKeyOrModel("gemini", "keys"); len(geminiKey) != 0 {
		if geminiModel := getKeyOrModel("gemini", "models"); len(geminiModel) != 0 {
			resources.gemini.Keys = geminiKey
			resources.gemini.Models = geminiModel
		}
	}
	// if chatGPTKey := getKeyOrModel("chatgpt", "keys"); len(chatGPTKey) != 0 {
	// 	if chatGPTModel := getKeyOrModel("chatgpt", "models"); len(chatGPTModel) != 0 {
	// 		resources.chatgpt.Keys = chatGPTKey
	// 		resources.chatgpt.Models = chatGPTModel
	// 	}
	// }
	// if claudeKey := getKeyOrModel("claude", "keys"); len(claudeKey) != 0 {
	// 	if claudeModel := getKeyOrModel("claude", "models"); len(claudeModel) != 0 {
	// 		resources.claude.Keys = claudeKey
	// 		resources.claude.Models = claudeModel
	// 	}
	// }
	// if deepSeekKey := getKeyOrModel("deepseek", "keys"); len(deepSeekKey) != 0 {
	// 	if deepSeekModel := getKeyOrModel("deepseek", "models"); len(deepSeekModel) != 0 {
	// 		resources.deepSeek.Keys = deepSeekKey
	// 		resources.deepSeek.Models = deepSeekModel
	// 	}
	// }
	return &resources, nil
}
func getKeyOrModel(whichKeyOrModel string, keyOrModel string) []string {
	var keys []string
	i := 1
	for {
		k := viper.GetString(fmt.Sprintf("%s.%s.%d", whichKeyOrModel, keyOrModel, i))
		i++
		if k == "" {
			break
		}
		keys = append(keys, k)
	}
	return keys
}

func (s *Resources) Run() {
	AutoMigrate(s.DB)
	// Start GRPC Server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", viper.GetInt("grpc.port")))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	aiRepository := repository.NewAiRepository(s.chatgpt, s.claude, s.gemini, s.deepSeek)
	aiUsecase := usecase.NewAiService(aiRepository)
	aiServer := server.NewAiServer(aiUsecase)
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(LogResponsesInterceptor))
	pb.RegisterAiServiceServer(grpcServer, aiServer)
	log.Println("Server started on port :", viper.GetInt("grpc.port"))
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
