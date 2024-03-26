package working_hours

import "github.com/mymmrac/telego"

var StartCmd = telego.BotCommand{
	Command:     "/start",
	Description: "get app info and instruction",
}

var HelpCmd = telego.BotCommand{
	Command:     "/help",
	Description: "get app instruction",
}

var ListWorkingHours = telego.BotCommand{
	Command:     "/listworkinghours",
	Description: "get a list of your working hours",
}
