package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/pttrulez/investor-go-next/go-api/internal/domain"
	"github.com/pttrulez/investor-go-next/go-api/pkg/protogen"

	"google.golang.org/grpc"
)

func (s *GRPCServer) GetPortfolioList(ctx context.Context,
	req *protogen.PortfolioListRequest) (*protogen.PortfolioListResponse, error) {
	chatId := req.GetChatId()
	ps, err := s.portfolioService.GetPortfolioListByChatID(ctx, chatId)
	if err != nil {
		return nil, err
	}
	fmt.Println("GRPCServer GetPortfolioList chatId", chatId)

	psProto := make([]*protogen.Portfolio, len(ps))

	for i, p := range ps {
		var pr protogen.Portfolio
		pr.Id = int32(p.ID)
		pr.Name = p.Name
		psProto[i] = &pr
	}
	return &protogen.PortfolioListResponse{Portfolios: psProto}, nil
}

func (s *GRPCServer) GetPortfolioSummaryMessage(ctx context.Context,
	req *protogen.PortfolioRequest) (*protogen.PortfolioSummaryResponse, error) {

	msg, err := s.portfolioService.GetPortfolioSummary(ctx, int(req.GetId()))
	if err != nil {
		return nil, err
	}

	return &protogen.PortfolioSummaryResponse{
		Text: msg,
	}, nil
}

func (s *GRPCServer) MustStart(port string) {
	// Make a TCP Listener
	ln, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	// Make a new GRPC native server with (options)
	grpcServer := grpc.NewServer([]grpc.ServerOption{}...)

	// Register (OUR) GRPC server implementation to the GRPC package.
	protogen.RegisterInvestorServer(grpcServer, s)
	fmt.Println("GRPC Server is running on port", port)

	err = grpcServer.Serve(ln)
	if err != nil {
		log.Fatal(err)
	}
}

type GRPCServer struct {
	protogen.UnimplementedInvestorServer
	portfolioService PortfolioService
}

type PortfolioService interface {
	GetPortfolioListByChatID(ctx context.Context, chatId string) ([]domain.Portfolio, error)
	GetPortfolioSummary(ctx context.Context, portfolioId int) (string, error)
}

func NewGRPCServer(portfolioService PortfolioService) *GRPCServer {
	return &GRPCServer{
		portfolioService: portfolioService,
	}
}
