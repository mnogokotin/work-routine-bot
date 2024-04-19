package task

import "github.com/mymmrac/telego"

var (
	MyTasks = telego.BotCommand{
		Command:     "mytasks",
		Description: "edit your tasks",
	}
	AddTask = telego.BotCommand{
		Command:     "newtask",
		Description: "add task record",
	}
	ListTasks = telego.BotCommand{
		Command:     "listtasks",
		Description: "get a list of your tasks",
	}
)
