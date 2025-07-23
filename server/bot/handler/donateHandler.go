package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (b1 *BotHandler) DonateHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	fmt.Println("Stars handler triggered")
	msg, err := b.SendInvoice(ctx, &bot.SendInvoiceParams{
		ChatID:        update.Message.Chat.ID,
		Title:         "Stars",
		Description:   "Support us with stars!",
		Payload:       "Hello, " + update.Message.From.Username + "!",
		ProviderToken: "",
		Currency:      "XTR",
		Prices: []models.LabeledPrice{
			{
				Label:  "Stars",
				Amount: 1,
			},
		},
	})

	if err != nil {
		fmt.Println("Error sending invoice:", err)
		return
	}
	jb, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("Error marshalling invoice:", err)
		return
	}
	fmt.Println("Invoice sent successfully:", string(jb))
}
