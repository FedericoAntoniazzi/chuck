package connector

import (
	"context"
	"fmt"

	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/telegram"
)

type ConsoleConnector struct{}

func (cc ConsoleConnector) Send(title, message string) error {
	fmt.Printf("%s\n\n%s", title, message)
	return nil
}

type TelegramConnector struct {
	notify  *notify.Notify
	service *telegram.Telegram
}

func NewTelegramConnector(botToken string, chatID int64) (*TelegramConnector, error) {
	telegramService, err := telegram.New(botToken)
	if err != nil {
		return nil, err
	}

	telegramService.AddReceivers(chatID)

	notifier := notify.New()

	notifier.UseServices(telegramService)

	return &TelegramConnector{
		service: telegramService,
		notify:  notifier,
	}, nil
}

func (tc TelegramConnector) Send(title, message string) error {
	ctx := context.Background()
	return tc.notify.Send(ctx, title, message)
}
