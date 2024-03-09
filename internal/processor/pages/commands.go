package pages

import "github.com/mymmrac/telego"

var HelpCmd = telego.BotCommand{
	Command:     "/help",
	Description: "get app info",
}

var RandomCmd = telego.BotCommand{
	Command:     "/random",
	Description: "get app instruction",
}

var StartCmd = telego.BotCommand{
	Command:     "/start",
	Description: "get random page from your list",
}

var Cmds = []telego.BotCommand{
	StartCmd,
	HelpCmd,
	RandomCmd,
}
