package telegram

import (
	"errors"
	"log"
	"net/url"
	"strings"
	"telegram-bot/clients/telegram"
	"telegram-bot/lib/e"
	"telegram-bot/storage"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) doCmd(text string, chatId int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s'", text, username)

	if isAddCmd(text) {
		return p.savePage(chatId, text, username)
	}

	switch text {
	case RndCmd:
		return p.sendRandom(chatId, username)
	case HelpCmd:
		return p.sendHelp(chatId)
	case StartCmd:
		return p.sendHello(chatId)
	default:
		return p.tg.SendMessage(chatId, msgUnknownCommand)
	}
}

func (p *Processor) savePage(chatID int, pageUrl string, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: save page", err) }()

	page := &storage.Page{
		URL:      pageUrl,
		UserName: username,
	}

	sender := newMessageSender(chatID, p.tg)

	isExists, err := p.storage.IsExists(page)
	if err != nil {
		return err
	}
	if isExists {
		return sender(msgAlreadyExists)
	}

	if err := p.storage.Save(page); err != nil {
		return err
	}

	if err := sender(msgSaved); err != nil {
		return err
	}

	return nil
}

func (p *Processor) sendRandom(chatID int, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: can't send random", err) }()

	sender := newMessageSender(chatID, p.tg)
	page, err := p.storage.PickRandom(username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}

	if errors.Is(err, storage.ErrNoSavedPages) {
		return sender(msgNoSavedPages)
	}

	if err := sender(page.URL); err != nil {
		return err
	}

	return nil
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}

func newMessageSender(chatID int, tg *telegram.Client) func(string) error {
	return func(msg string) error {
		return tg.SendMessage(chatID, msg)
	}
}

func isAddCmd(text string) bool {
	u, err := url.Parse(text)
	return err == nil && u.Host != ""
}
