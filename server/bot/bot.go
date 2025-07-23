package bot

import (
	"context"
	"fmt"

	botHandler "tg-stars/bot/handler"

	"github.com/go-telegram/bot"
	"gorm.io/gorm"
)

type Bot interface {
	Start(ctx context.Context) error
	GetBot() *bot.Bot
	SendMessage(chatId int64, text string)
}

type botImpl struct {
	bot *bot.Bot
	db  *gorm.DB
}

func NewBot(botToken string, db *gorm.DB) Bot {
	bot_handler := botHandler.NewBotHandler(db)
	opts := []bot.Option{
		bot.WithDefaultHandler(bot_handler.DefaultHandler),
	}

	b, err := bot.New(botToken, opts...)
	if err != nil {
		panic(err)
	}

	bot_impl := botImpl{
		bot: b,
		db:  db,
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "start", bot.MatchTypeCommand, bot_handler.StartHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "donate", bot.MatchTypeCommand, bot_handler.DonateHandler)

	b.RegisterHandler(bot.HandlerTypeMessageText, "memories", bot.MatchTypeCommandStartOnly, bot_handler.MemoriesHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "clear", bot.MatchTypeCommand, bot_handler.ClearHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "reset", bot.MatchTypeCommand, bot_handler.ResetHandler)

	b.RegisterHandler(bot.HandlerTypeMessageText, "share", bot.MatchTypeCommand, bot_handler.ShareHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "import", bot.MatchTypeCommand, bot_handler.ImportHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "export", bot.MatchTypeCommand, bot_handler.ExportHandler)

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
