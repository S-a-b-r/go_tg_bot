package telegram

const msgHelp = `I can save and keep you pages. Also I can offer you them to read.

In order to save the page, just send me all link to it.

In order to get a random page form your list, send me command /rnd.`

const msgHello = "Hi there! \n\n" + msgHelp

const (
	msgUnknownCommand = "Sorry, it's unknown command"
	msgSaved          = "Success saved!"
	msgAlreadyExists  = "You have already have this page in your list"
	msgNoSavedPages   = "Already not have saved links"
)
