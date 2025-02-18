package usecase

import (
	"ai/api/pb"
	"ai/internal/repository"
)

type AiUsecase interface {
	Ask(*pb.AiRequest) (*pb.AiResponse, error)
}
type aiUsecase struct {
	repository repository.AiRepository
}

func NewAiService(repository repository.AiRepository) AiUsecase {
	return &aiUsecase{repository}
}
func (a *aiUsecase) Ask(req *pb.AiRequest) (*pb.AiResponse, error) {
	return a.repository.Ask(req)
}
