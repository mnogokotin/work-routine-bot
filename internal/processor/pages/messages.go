package pages

const MsgHelp = `I can store and keep you pages. Also I can offer you them to read.

In order to store the page, just send me a link to it.

In order to get a random page from your list, send me command /random.
Caution! After that, this page will be removed from your list!`

const (
	MsgStart          = "Hi there! 👾\n\n" + MsgHelp
	MsgUnknownCommand = "Unknown command 🤔"
	MsgNoStoredPages  = "You have no stored pages 🙊"
	MsgStored         = "Stored! 👌"
	MsgAlreadyExists  = "You have already have this page in your list 🤗"
)
