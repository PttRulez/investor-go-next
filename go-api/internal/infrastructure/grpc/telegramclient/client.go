package telegramclient

import (
	"context"

	"github.com/pttrulez/invest_telega/pkg/protogen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TelegramGRPCClient struct {
	protogen.TelegaClient
	grpcClient protogen.TelegaClient
}

func NewTelegramClient(endpoint string) (*TelegramGRPCClient, error) {
	conn, err := grpc.NewClient(endpoint, grpc.WithTransportCredentials(
		insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	grpcClient := protogen.NewTelegaClient(conn)

	return &TelegramGRPCClient{grpcClient: grpcClient}, nil
}

func (c *TelegramGRPCClient) SendMsg(ctx context.Context, msgInfo *protogen.MessageInfo) error {
	_, err := c.grpcClient.SendMsg(ctx, msgInfo)
	return err
}
