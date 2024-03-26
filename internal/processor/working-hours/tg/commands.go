package tg

import (
	"context"
	"errors"
	"github.com/mnogokotin/golang-packages/utils/e"
	"github.com/mymmrac/telego"
	"work-routine-bot/internal/domain"
	"work-routine-bot/internal/processor/working-hours"
	"work-routine-bot/internal/storage/users"
)

func (p *Processor) executeCmd(ctx context.Context, u telego.Update) error {
	chatID := u.Message.Chat.ID
	username := u.Message.From.Username
	text := u.Message.Text

	p.log.Info("got new command", username, text)

	switch text {
	case working_hours.HelpCmd.Command:
		return p.handleStartCmd(ctx, chatID)
	case working_hours.StartCmd.Command:
		return p.handleStartCmd(ctx, chatID)
	case working_hours.ListWorkingHours.Command:
		return p.handleListWorkingHoursCmd(ctx, chatID, username)
	default:
		return p.handleUnknownCommandMsg(ctx, chatID)
	}
}

func (p *Processor) handleStartCmd(ctx context.Context, chatID int64) error {
	return p.SendStartMessage(ctx, chatID, working_hours.MsgStart)
}

func (p *Processor) handleUnknownCommandMsg(ctx context.Context, chatID int64) error {
	return p.SendMessage(ctx, chatID, working_hours.MsgUnknownCommand, nil)
}

func (p *Processor) handleListWorkingHoursCmd(ctx context.Context, chatID int64, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't handle command 'ListWorkingHours'", err) }()
	user, err := p.userProvider.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, users.ErrUserNotFound) {
			newUser := domain.User{Username: username}
			user_, err := p.userProvider.Store(ctx, &newUser)
			if err != nil {
				return err
			}
			user = user_
		} else {
			return err
		}
	}

	tasks, err := p.taskProvider.GetListByUserId(ctx, user.ID)
	if err != nil {
		return err
	}

	return p.SendListWorkingHoursMessage(ctx, chatID, tasks)
}
