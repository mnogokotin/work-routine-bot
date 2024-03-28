package working_hours

import "github.com/mymmrac/telego"

var (
	ListWorkingHours = telego.BotCommand{
		Command:     "listworkinghours",
		Description: "get a list of your working hours",
	}
)
