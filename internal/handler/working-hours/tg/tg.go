package tg

import (
	"context"
	"errors"
	"fmt"
	"github.com/mnogokotin/golang-packages/utils/e"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"log/slog"
	"work-routine-bot/internal/bot"
	"work-routine-bot/internal/domain"
	wh "work-routine-bot/internal/handler/working-hours"
	"work-routine-bot/internal/storage/users"
)

type UserProvider interface {
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	Store(ctx context.Context, user *domain.User) (*domain.User, error)
}

type TaskProvider interface {
	GetListByUserId(ctx context.Context, userId int) ([]*domain.Task, error)
	Store(ctx context.Context, task *domain.Task) (*domain.Task, error)
}

type Handler struct {
	log          *slog.Logger
	bot          bot.Bot
	userProvider UserProvider
	taskProvider TaskProvider
}

func New(log *slog.Logger, bot bot.Bot, userProvider UserProvider, taskProvider TaskProvider) *Handler {
	return &Handler{
		log:          log,
		bot:          bot,
		userProvider: userProvider,
		taskProvider: taskProvider,
	}
}

func (h *Handler) Handle() {
	h.bot.Bh.HandleMessageCtx(func(ctx context.Context, bot *telego.Bot, message telego.Message) {
		c := wh.ListWorkingHours.Command

		username := message.From.Username

		h.log.Info("got new command", username, message.Text)

		user, err := h.userProvider.GetByUsername(ctx, username)
		if err != nil {
			if errors.Is(err, users.ErrUserNotFound) {
				newUser := domain.User{Username: username}
				user_, err := h.userProvider.Store(ctx, &newUser)
				if err != nil {
					h.log.Error("", "", e.Wrap(c, err).Error())
					return
				}
				user = user_
			} else {
				h.log.Error("", "", e.Wrap(c, err).Error())
				return
			}
		}

		tasks, err := h.taskProvider.GetListByUserId(ctx, user.ID)
		if err != nil {
			h.log.Error("3", "", e.Wrap(c, err).Error())
			return
		}

		_, _ = bot.SendMessage(
			tu.Message(
				tu.ID(message.Chat.ID),
				h.BuildListWorkingHoursMessage(tasks),
			),
		)
	}, th.CommandEqual(wh.ListWorkingHours.Command))
}

func (h *Handler) BuildListWorkingHoursMessage(tasks []*domain.Task) string {
	var message string

	for _, task := range tasks {
		message += fmt.Sprintf("%d %d %s\n", task.ID, task.ProjectId, task.Description)
	}

	return wh.MsgList + message
}
