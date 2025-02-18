package server

import (
	"ai/api/pb"
	"ai/internal/usecase"
	"context"
)

type aiServer struct {
	usecase usecase.AiUsecase
	pb.UnimplementedAiServiceServer
}

func NewAiServer(usecase usecase.AiUsecase) pb.AiServiceServer {
	return &aiServer{usecase: usecase}
}

func (s *aiServer) Ask(ctx context.Context, req *pb.AiRequest) (*pb.AiResponse, error) {
	res, err := s.usecase.Ask(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
