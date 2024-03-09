package tg

import (
	"context"
	"errors"
	"github.com/mnogokotin/golang-packages/utils/e"
	"github.com/mymmrac/telego"
	"net/url"
	"work-routine-bot/internal/domain"
	"work-routine-bot/internal/processor/pages"
	spages "work-routine-bot/internal/storage/pages"
)

func (p *Processor) executeCmd(ctx context.Context, u telego.Update) error {
	chatID := u.Message.Chat.ID
	username := u.Message.From.Username
	text := u.Message.Text

	p.log.Info("got new command", username, text)

	if AddNewPageCmd(text) {
		return p.handleAddNewPageCmd(ctx, chatID, text, username)
	}

	switch text {
	case pages.HelpCmd.Command:
		return p.handleHelpCmd(ctx, chatID)
	case pages.StartCmd.Command:
		return p.handleStartCmd(ctx, chatID)
	case pages.RandomCmd.Command:
		return p.handleRandomCmd(ctx, chatID, username)
	default:
		return p.sendUnknownCommandMsg(ctx, chatID)
	}
}

func (p *Processor) handleAddNewPageCmd(ctx context.Context, chatID int64, pageUrl string, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't handle command 'add new page'", err) }()

	page := &domain.Page{
		Url:      pageUrl,
		Username: username,
	}

	isExists, err := p.storage.IsExists(ctx, page)
	if err != nil {
		return err
	}

	if isExists {
		return p.SendMessage(ctx, chatID, pages.MsgAlreadyExists)
	}

	if err := p.storage.Store(ctx, page); err != nil {
		return err
	}

	if err := p.SendMessage(ctx, chatID, pages.MsgStored); err != nil {
		return err
	}

	return nil
}

func (p *Processor) handleStartCmd(ctx context.Context, chatID int64) error {
	return p.SendMessage(ctx, chatID, pages.MsgStart)
}

func (p *Processor) handleHelpCmd(ctx context.Context, chatID int64) error {
	return p.SendMessage(ctx, chatID, pages.MsgHelp)
}

func (p *Processor) sendUnknownCommandMsg(ctx context.Context, chatID int64) error {
	return p.SendMessage(ctx, chatID, pages.MsgUnknownCommand)
}

func (p *Processor) handleRandomCmd(ctx context.Context, chatID int64, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't handle command 'random'", err) }()

	page, err := p.storage.GetRandom(ctx, username)

	if err != nil {
		if errors.Is(err, spages.ErrNoStoredPages) {
			return p.SendMessage(ctx, chatID, pages.MsgNoStoredPages)
		}

		return err
	}

	if err := p.SendMessage(ctx, chatID, page.Url); err != nil {
		return err
	}

	return p.storage.Remove(ctx, page)
}

func AddNewPageCmd(text string) bool {
	return isUrl(text)
}

func isUrl(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}
