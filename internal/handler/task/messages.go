package task

const MsgMyTasks = `Choose a command from the list below:
`
const MsgListTasks = `List of tasks:
`
const MsgAddTask = `Add task\. Choose project\:`
const MsgAddTaskDesc = `Send me description:`
const MsgAddTaskDuration = `Send me duration in a "hours minutes" format (example: "00 30" "01 45"):`
const MsgAddTaskDurationFormatErr = `Wrong duration "hours minutes" format (example: "10 01" "20 02"), please, send it again':`
const MsgAddTaskDate = `Send me date in a "day month" format (example: "10 01" "20 02"):`
const MsgAddTaskDateFormatErr = `Wrong date "day month" format (example: "10 01" "20 02"), please, send it again':`
const MsgAddTaskCantStoreErr = `Can't store task, try again later`
const MsgAddTaskDone = `Task %d added`
const MsgDeleteTask = `Send me integer task's id for deletion:`
const MsgDeleteTaskFormatErr = `Wrong task's id value, it must be integer. Please, send me integer task's id for deletion again:':`
const MsgDeleteTaskNotFound = `Task not found`
const MsgDeleteTaskCantStoreErr = `Can't delete task, try again later`
const MsgDeleteTaskDone = `Task %d deleted`
