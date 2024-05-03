package tg

import (
	"context"
	"errors"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/mnogokotin/golang-packages/utils/e"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"log/slog"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
	"work-routine-bot/internal/bot"
	"work-routine-bot/internal/domain"
	"work-routine-bot/internal/fsm"
	"work-routine-bot/internal/handler/task"
)

type ProjectProvider interface {
	GetList(ctx context.Context) ([]*domain.Project, error)
}

type TaskProvider interface {
	Delete(ctx context.Context, id int) error
	GetById(ctx context.Context, id int) (*domain.Task, error)
	GetListByUserId(ctx context.Context, userId int) ([]*domain.Task, error)
	Store(ctx context.Context, task *domain.Task) (*domain.Task, error)
}

type UserProvider interface {
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
}

type Handler struct {
	log             *slog.Logger
	bot             bot.Bot
	projectProvider ProjectProvider
	taskProvider    TaskProvider
	userProvider    UserProvider
}

func New(log *slog.Logger, bot bot.Bot, projectProvider ProjectProvider, taskProvider TaskProvider, userProvider UserProvider) *Handler {
	return &Handler{
		log:             log,
		bot:             bot,
		projectProvider: projectProvider,
		taskProvider:    taskProvider,
		userProvider:    userProvider,
	}
}

type AddTaskInput struct {
	ProjectId   int
	Description string
	Duration    int
	Date        time.Time
}

type DeleteTaskInput struct {
	ID int
}

func (h *Handler) Handle() {
	h.bot.Bh.Handle(func(bot *telego.Bot, update telego.Update) {
		_, _ = bot.SendMessage(
			tu.Message(
				tu.ID(update.Message.Chat.ID),
				task.MsgMyTasks,
			).WithParseMode(telego.ModeMarkdownV2).WithReplyMarkup(
				tu.InlineKeyboard(
					tu.InlineKeyboardRow(
						tu.InlineKeyboardButton("List tasks").WithCallbackData("list"),
						tu.InlineKeyboardButton("Add task").WithCallbackData("add"),
					),
					tu.InlineKeyboardRow(
						tu.InlineKeyboardButton("Edit task").WithCallbackData("edit"),
						tu.InlineKeyboardButton("Delete task").WithCallbackData("delete"),
					),
				),
			),
		)
	}, th.CommandEqual(task.MyTasks.Command))

	h.bot.Bh.HandleCallbackQuery(func(bot *telego.Bot, query telego.CallbackQuery) {
		c := task.ListTasks.Command
		username := query.From.Username
		ctx := context.Background()

		user, err := h.userProvider.GetByUsername(ctx, username)
		if err != nil {
			h.log.Error("", "", e.Wrap(c, err).Error())
			return
		}

		tasks, err := h.taskProvider.GetListByUserId(ctx, user.ID)
		if err != nil {
			h.log.Error("", "", e.Wrap(c, err).Error())
			return
		}

		projects, err := h.projectProvider.GetList(ctx)
		if err != nil {
			h.log.Error("", "", e.Wrap(c, err).Error())
			return
		}

		_, _ = bot.SendMessage(
			tu.Message(
				tu.ID(query.Message.GetChat().ID),
				task.MsgListTasks+h.BuildListTasksMessage(tasks, projects),
			),
		)
	}, th.AnyCallbackQueryWithMessage(), th.CallbackDataEqual("list"))

	h.bot.Bh.HandleCallbackQuery(func(bot *telego.Bot, query telego.CallbackQuery) {
		c := task.AddTask.Command + " project name"
		ctx := context.Background()
		h.bot.Fsm.Fsm = fsm.NewConvFsm("add-task", AddTaskInput{})
		h.bot.Fsm.Data = AddTaskInput{}

		projects, err := h.projectProvider.GetList(ctx)
		if err != nil {
			h.log.Error("", "", e.Wrap(c, err).Error())
			return
		}

		var buttons []telego.InlineKeyboardButton

		for _, project := range projects {
			buttons = append(buttons, tu.InlineKeyboardButton(project.Name).WithCallbackData("add-task-projectid-"+strconv.Itoa(project.ID)))
		}

		var rows [][]telego.InlineKeyboardButton
		var rowSize = 3

		for i := 0; i < len(buttons); i += rowSize {
			end := i + rowSize

			if end > len(projects) {
				end = len(projects)
			}

			rows = append(rows, buttons[i:end])
		}

		_, _ = bot.SendMessage(
			tu.Message(
				tu.ID(query.Message.GetChat().ID),
				task.MsgAddTask,
			).WithParseMode(telego.ModeMarkdownV2).WithParseMode(telego.ModeMarkdownV2).WithReplyMarkup(
				tu.InlineKeyboard(rows...),
			),
		)
	}, th.AnyCallbackQueryWithMessage(), th.CallbackDataEqual("add"))

	h.bot.Bh.HandleCallbackQuery(func(bot *telego.Bot, query telego.CallbackQuery) {
		ctx := context.Background()
		fsmData := h.bot.Fsm.Data.(AddTaskInput)
		fsmData.ProjectId, _ = strconv.Atoi(strings.Replace(query.Data, "add-task-projectid-", "", -1))
		h.bot.Fsm.Data = fsmData
		_ = h.bot.Fsm.Fsm.Event(ctx, "get-projectid")

		_, _ = bot.SendMessage(
			tu.Message(
				tu.ID(query.Message.GetChat().ID),
				task.MsgAddTaskDesc,
			),
		)
	}, th.AnyCallbackQueryWithMessage(), th.CallbackDataContains("add-task-projectid"))

	h.bot.Bh.Handle(func(bot *telego.Bot, update telego.Update) {
		fsmData := h.bot.Fsm.Data.(AddTaskInput)
		fsmData.Description = update.Message.Text
		h.bot.Fsm.Data = fsmData
		_ = h.bot.Fsm.Fsm.Event(update.Context(), "get-description")

		_, _ = bot.SendMessage(
			tu.Message(
				tu.ID(update.Message.Chat.ID),
				task.MsgAddTaskDuration,
			),
		)
	}, fsm.FsmStateEqual(h.bot.Fsm, "add-task-get-description"))

	h.bot.Bh.Handle(func(bot *telego.Bot, update telego.Update) {
		c := task.AddTask.Command + " project duration"

		durationParts := strings.Split(update.Message.Text, " ")
		if len(durationParts) != 2 {
			h.log.Error("", "", e.Wrap(c, errors.New("wrong duration format")).Error())
			_, _ = bot.SendMessage(
				tu.Message(
					tu.ID(update.Message.Chat.ID),
					task.MsgAddTaskDurationFormatErr,
				),
			)
			return
		}

		durationHours, err1 := strconv.ParseInt(durationParts[0], 0, 64)
		durationMinutes, err2 := strconv.ParseInt(durationParts[1], 0, 64)
		if err1 != nil || err2 != nil || durationMinutes > 59 {
			var err error
			if err1 != nil {
				err = err1
			} else if err2 != nil {
				err = err2
			} else {
				err = errors.New("wrong minutes format")
			}

			h.log.Error("", "", e.Wrap(c, err).Error())
			_, _ = bot.SendMessage(
				tu.Message(
					tu.ID(update.Message.Chat.ID),
					task.MsgAddTaskDurationFormatErr,
				),
			)
			return
		}

		fsmData := h.bot.Fsm.Data.(AddTaskInput)
		fsmData.Duration = int((durationHours * 60) + durationMinutes)
		h.bot.Fsm.Data = fsmData
		_ = h.bot.Fsm.Fsm.Event(update.Context(), "get-duration")

		_, _ = bot.SendMessage(
			tu.Message(
				tu.ID(update.Message.Chat.ID),
				task.MsgAddTaskDate,
			),
		)
	}, fsm.FsmStateEqual(h.bot.Fsm, "add-task-get-duration"))

	h.bot.Bh.Handle(func(bot *telego.Bot, update telego.Update) {
		c := task.AddTask.Command + " date"

		dateParts := strings.Split(update.Message.Text, " ")
		if len(dateParts) != 2 {
			h.log.Error("", "", e.Wrap(c, errors.New("wrong date format")).Error())
			_, _ = bot.SendMessage(
				tu.Message(
					tu.ID(update.Message.Chat.ID),
					task.MsgAddTaskDateFormatErr,
				),
			)
			return
		}

		dateDay, err1 := strconv.ParseInt(dateParts[0], 0, 64)
		dateMonth, err2 := strconv.ParseInt(dateParts[1], 0, 64)
		if err1 != nil || err2 != nil || dateDay > 31 || dateMonth > 12 {
			var err error
			if err1 != nil {
				err = err1
			} else if err2 != nil {
				err = err2
			} else {
				err = errors.New("wrong minutes format")
			}

			h.log.Error("", "", e.Wrap(c, err).Error())
			_, _ = bot.SendMessage(
				tu.Message(
					tu.ID(update.Message.Chat.ID),
					task.MsgAddTaskDateFormatErr,
				),
			)
			return
		}

		dateString := strconv.Itoa(time.Now().Year()) + "-" + dateParts[1] + "-" + dateParts[0]
		date, err := time.Parse("2006-01-02", dateString)
		if err != nil {
			h.log.Error("", "", e.Wrap(c, err).Error())
			_, _ = bot.SendMessage(
				tu.Message(
					tu.ID(update.Message.Chat.ID),
					task.MsgAddTaskDateFormatErr,
				),
			)
			return
		}

		fsmData := h.bot.Fsm.Data.(AddTaskInput)
		fsmData.Date = date

		taskDomain := domain.Task{
			UserId:      int(update.Message.From.ID),
			ProjectId:   fsmData.ProjectId,
			Description: fsmData.Description,
			Duration:    fsmData.Duration,
			Date:        fsmData.Date,
			CreatedAt:   time.Now(),
		}

		storedTask_, err := h.taskProvider.Store(update.Context(), &taskDomain)
		if err != nil {
			h.log.Error("", "", e.Wrap(c, err).Error())
			_, _ = bot.SendMessage(
				tu.Message(
					tu.ID(update.Message.Chat.ID),
					task.MsgAddTaskCantStoreErr,
				),
			)
			return
		}

		_, _ = bot.SendMessage(
			tu.Messagef(
				tu.ID(update.Message.Chat.ID),
				task.MsgAddTaskDone,
				&storedTask_.ID,
			),
		)
	}, fsm.FsmStateEqual(h.bot.Fsm, "add-task-get-date"))

	h.bot.Bh.HandleCallbackQuery(func(bot *telego.Bot, query telego.CallbackQuery) {
		h.bot.Fsm.Fsm = fsm.NewConvFsm("delete-task", DeleteTaskInput{})

		_, _ = bot.SendMessage(
			tu.Message(
				tu.ID(query.Message.GetChat().ID),
				task.MsgDeleteTask,
			),
		)
	}, th.AnyCallbackQueryWithMessage(), th.CallbackDataEqual("delete"))

	h.bot.Bh.Handle(func(bot *telego.Bot, update telego.Update) {
		c := task.DeleteTask.Command + " task_id"
		taskId, err := strconv.Atoi(update.Message.Text)

		if err != nil {
			h.log.Error("", "", e.Wrap(c, err).Error())
			_, _ = bot.SendMessage(
				tu.Message(
					tu.ID(update.Message.Chat.ID),
					task.MsgDeleteTaskFormatErr,
				),
			)
			return
		}

		task_, err := h.taskProvider.GetById(update.Context(), taskId)
		if err != nil {
			h.log.Error("", "", e.Wrap(c, errors.New("task not found")).Error())
			_, _ = bot.SendMessage(
				tu.Message(
					tu.ID(update.Message.Chat.ID),
					task.MsgDeleteTaskNotFound,
				),
			)
			return
		}

		if task_.UserId != int(update.Message.From.ID) {
			_, _ = bot.SendMessage(
				tu.Message(
					tu.ID(update.Message.Chat.ID),
					task.MsgDeleteTaskNotFound,
				),
			)
			return
		}

		err = h.taskProvider.Delete(update.Context(), taskId)
		if err != nil {
			h.log.Error("", "", e.Wrap(c, errors.New("can't delete")).Error())
			_, _ = bot.SendMessage(
				tu.Message(
					tu.ID(update.Message.Chat.ID),
					task.MsgDeleteTaskCantStoreErr,
				),
			)
			return
		}

		_, _ = bot.SendMessage(
			tu.Messagef(
				tu.ID(update.Message.Chat.ID),
				task.MsgDeleteTaskDone,
				taskId,
			),
		)
	}, fsm.FsmStateEqual(h.bot.Fsm, "delete-task-get-id"))
}

func (h *Handler) BuildListTasksMessage(tasks []*domain.Task, projects []*domain.Project) string {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendRow(table.Row{"id", "project", "duration", "date", "desc"})

	for _, task_ := range tasks {
		idx := slices.IndexFunc(projects, func(p *domain.Project) bool { return p.ID == task_.ProjectId })
		durationText := fmt.Sprintf("%dh %dm", task_.Duration/60, task_.Duration%60)
		t.AppendRow([]interface{}{task_.ID, projects[idx].Name, durationText, task_.Date.Format("02 Jan"), task_.Description})
	}

	return t.RenderMarkdown()
}
