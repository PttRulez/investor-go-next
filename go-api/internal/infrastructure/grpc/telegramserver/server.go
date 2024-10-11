package server

import (
	"context"

	"github.com/pttrulez/investor-go-next/go-api/internal/domain"
	"github.com/pttrulez/investor-go-next/go-api/pkg/protogen"
)

func (s *GRPCTelegaServer) GetPortfolioList(ctx context.Context,
	req *protogen.PortfolioListRequest) (*protogen.PortfolioListResponse, error) {
	chatId := req.GetChatId()
	ps, err := s.portfolioService.GetPortfolioListByChatID(chatId)
	if err != nil {
		return nil, err
	}

	psProto := make([]*protogen.Portfolio, 0, len(ps))
	for i, p := range ps {
		psProto[i].Id = int32(p.ID)
		psProto[i].Name = p.Name
	}
	return &protogen.PortfolioListResponse{Portfolios: psProto}, nil
}

func (s *GRPCTelegaServer) GetPortfolioSummaryMessage(ctx context.Context,
	req *protogen.PortfolioRequest) (*protogen.PortfolioSummaryResponse, error) {

	msg, err := s.portfolioService.GetPortfolioSummary(int(req.GetId()), req.GetChatId())
	if err != nil {
		return nil, err
	}

	return &protogen.PortfolioSummaryResponse{
		Text: msg,
	}, nil
}

type GRPCTelegaServer struct {
	protogen.UnimplementedTelegaServer
	portfolioService PortfolioService
}

type PortfolioService interface {
	GetPortfolioListByChatID(chatId string) ([]domain.Portfolio, error)
	GetPortfolioSummary(portfolioId int, chatId string) (string, error)
}

func NewGRPCTelegaServer(portfolioService PortfolioService) *GRPCTelegaServer {
	return &GRPCTelegaServer{
		portfolioService: portfolioService,
	}
}
