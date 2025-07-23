package handler

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (bh *BotHandler) ImportHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	fmt.Println("ImportHandler called with update:", update.Message.Text)
	bh.pendingUsers[update.Message.From.ID] = PendingOperationImport
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text: "Please upload the memory file you want to import.",
	})
}
