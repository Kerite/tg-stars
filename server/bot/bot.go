package bot

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Bot interface {
	Start(ctx context.Context) error
	GetBot() *bot.Bot
	SendMessage(chatId int64, text string)
}

type botImpl struct {
	bot *bot.Bot
}

func NewBot(botToken string) Bot {
	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
	}

	b, err := bot.New(botToken, opts...)
	if err != nil {
		panic(err)
	}

	bot_impl := botImpl{
		bot: b,
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, startHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/donate", bot.MatchTypeExact, bot_impl.donateHandler)

	return &bot_impl
}

func (b *botImpl) GetBot() *bot.Bot {
	return b.bot
}

func (b *botImpl) Start(ctx context.Context) error {
	fmt.Println("[Bot] Starting bot...")
	b.bot.Start(ctx)

	return nil
}

func (b *botImpl) SendMessage(chatId int64, text string) {
	b.bot.SendMessage(context.Background(), &bot.SendMessageParams{
		ChatID: chatId,
		Text:   text,
	})
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	fmt.Println("ChatID:", update.Message.Chat.ID)
	fmt.Println("Received Message:", update.Message.Text)
	fmt.Println(json.MarshalIndent(*update, "", "  "))
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   update.Message.Text,
	})
}

func startHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Welcome to the bot!",
		})
	}
}

func (b1 *botImpl) donateHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	fmt.Println("Stars handler triggered")
	msg, err := b.SendInvoice(ctx, &bot.SendInvoiceParams{
		ChatID:        update.Message.Chat.ID,
		Title:         "Stars",
		Description:   "Support us with stars!",
		Payload:       "Hello",
		ProviderToken: "",
		Currency:      "XTR",
		Prices: []models.LabeledPrice{
			{
				Label:  "Stars",
				Amount: 1,
			},
			{
				Label:  "Support",
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
