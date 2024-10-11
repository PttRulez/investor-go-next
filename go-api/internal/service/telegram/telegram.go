package telegram

import (
	"context"
	"fmt"
	"log"

	"github.com/pttrulez/invest_telega/pkg/protogen"
	"github.com/pttrulez/investor-go-next/go-api/internal/infrastructure/grpc/telegramclient"
	"github.com/pttrulez/investor-go-next/go-api/internal/utils"

	"github.com/pttrulez/investor-go-next/go-api/pkg/logger"
)

func (t *Telegram) SendMsg(ctx context.Context, text string) {
	chatId := utils.GetCurrentUserTgChatID(ctx)
	fmt.Println("CHATID", chatId)
	if chatId != "" {
		err := t.grpcClient.SendMsg(ctx, &protogen.MessageInfo{
			ChatId: chatId,
			Text:   text,
		})
		if err != nil {
			t.logger.Debug(fmt.Sprintf("Telegram.SendMsg: %v", err))
		}
	}
}

func New(endpoint string, logger *logger.Logger) *Telegram {
	grpcClient, err := telegramclient.NewTelegramClient(endpoint)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to start tgClient: %w", err))
	}

	return &Telegram{
		grpcClient: grpcClient,
		logger:     logger,
	}
}

type Telegram struct {
	logger     *logger.Logger
	grpcClient *telegramclient.TelegramGRPCClient
}
