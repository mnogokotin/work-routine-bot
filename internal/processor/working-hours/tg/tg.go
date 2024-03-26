package tg

import (
	"context"
	"fmt"
	"github.com/mnogokotin/golang-packages/utils/e"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"log/slog"
	"work-routine-bot/internal/bot"
	"work-routine-bot/internal/domain"
)

type UserProvider interface {
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	Store(ctx context.Context, user *domain.User) (*domain.User, error)
}

type TaskProvider interface {
	GetListByUserId(ctx context.Context, userId int) ([]*domain.Task, error)
	Store(ctx context.Context, task *domain.Task) (*domain.Task, error)
}

type Processor struct {
	log          *slog.Logger
	bot          bot.Bot
	userProvider UserProvider
	taskProvider TaskProvider
}

func New(log *slog.Logger, bot bot.Bot, userProvider UserProvider, taskProvider TaskProvider) *Processor {
	return &Processor{
		log:          log,
		bot:          bot,
		userProvider: userProvider,
		taskProvider: taskProvider,
	}
}

func (p *Processor) GetChan() <-chan telego.Update {
	return p.bot.Channel
}

func (p *Processor) Process(ctx context.Context, u telego.Update) error {
	if err := p.executeCmd(ctx, u); err != nil {
		return e.Wrap("can't process update", err)
	}

	return nil
}

func (p *Processor) SendMessage(ctx context.Context, chatID int64, text string, keyboard *telego.ReplyKeyboardMarkup) error {
	message := tu.Message(
		tu.ID(chatID),
		text,
	).WithParseMode(telego.ModeMarkdownV2)
	if keyboard != nil {
		message = message.WithReplyMarkup(keyboard)
	}

	_, err := p.bot.Bot.SendMessage(message)
	if err != nil {
		return e.Wrap("can't send message", err)
	}

	return nil
}

func (p *Processor) SendStartMessage(ctx context.Context, chatID int64, text string) error {
	keyboard := tu.Keyboard(
		tu.KeyboardRow(
			tu.KeyboardButton("Button 1"),
			tu.KeyboardButton("Button 2"),
		),
	).WithResizeKeyboard().WithInputFieldPlaceholder("Select something")
	return p.SendMessage(ctx, chatID, text, keyboard)
}

func (p *Processor) SendListWorkingHoursMessage(ctx context.Context, chatID int64, tasks []*domain.Task) error {
	var message string

	for _, task := range tasks {
		message += fmt.Sprintf("%d %d %s\n", task.ID, task.ProjectId, task.Description)
	}

	return p.SendMessage(ctx, chatID, message, nil)
}
