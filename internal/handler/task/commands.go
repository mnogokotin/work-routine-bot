package task

import "github.com/mymmrac/telego"

var (
	MyTasks = telego.BotCommand{
		Command:     "mytasks",
		Description: "edit your tasks",
	}
	ListTasks = telego.BotCommand{
		Command:     "listtasks",
		Description: "get a list of your tasks",
	}
	AddTask = telego.BotCommand{
		Command:     "addtask",
		Description: "add new task",
	}
	DeleteTask = telego.BotCommand{
		Command:     "deletetask",
		Description: "delete task",
	}
)
