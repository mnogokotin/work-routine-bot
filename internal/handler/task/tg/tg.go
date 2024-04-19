package tg

import (
	"context"
	"errors"
	"fmt"
	"github.com/mnogokotin/golang-packages/utils/e"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"log/slog"
	"strconv"
	"time"
	"work-routine-bot/internal/bot"
	"work-routine-bot/internal/domain"
	"work-routine-bot/internal/fsm"
	"work-routine-bot/internal/handler/task"
	"work-routine-bot/internal/storage/users"
)

type UserProvider interface {
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	Store(ctx context.Context, user *domain.User) (*domain.User, error)
}

type TaskProvider interface {
	GetListByUserId(ctx context.Context, userId int) ([]*domain.Task, error)
	Store(ctx context.Context, task *domain.Task) (*domain.Task, error)
}

type Handler struct {
	log          *slog.Logger
	bot          bot.Bot
	userProvider UserProvider
	taskProvider TaskProvider
}

func New(log *slog.Logger, bot bot.Bot, userProvider UserProvider, taskProvider TaskProvider) *Handler {
	return &Handler{
		log:          log,
		bot:          bot,
		userProvider: userProvider,
		taskProvider: taskProvider,
	}
}

type AddTaskInput struct {
	ProjectId   int
	Description string
	Duration    float64
	Date        time.Time
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
			if errors.Is(err, users.ErrUserNotFound) {
				newUser := domain.User{Username: username}
				user_, err := h.userProvider.Store(ctx, &newUser)
				if err != nil {
					h.log.Error("", "", e.Wrap(c, err).Error())
					return
				}
				user = user_
			} else {
				h.log.Error("", "", e.Wrap(c, err).Error())
				return
			}
		}

		tasks, err := h.taskProvider.GetListByUserId(ctx, user.ID)
		if err != nil {
			h.log.Error("", "", e.Wrap(c, err).Error())
			return
		}

		_, _ = bot.SendMessage(
			tu.Message(
				tu.ID(query.Message.GetChat().ID),
				h.BuildListTasksMessage(tasks),
			),
		)
	}, th.AnyCallbackQueryWithMessage(), th.CallbackDataEqual("list"))

	h.bot.Bh.HandleCallbackQuery(func(bot *telego.Bot, query telego.CallbackQuery) {
		h.bot.Fsm.Fsm = fsm.NewConvFsm("add-task", AddTaskInput{})
		h.bot.Fsm.Data = AddTaskInput{}

		_, _ = bot.SendMessage(
			tu.Message(
				tu.ID(query.Message.GetChat().ID),
				task.MsgAddTask,
			).WithParseMode(telego.ModeMarkdownV2),
		)
	}, th.AnyCallbackQueryWithMessage(), th.CallbackDataEqual("add"))

	h.bot.Bh.Handle(func(bot *telego.Bot, update telego.Update) {
		fsmData := h.bot.Fsm.Data.(AddTaskInput)
		fsmData.ProjectId = 1
		h.bot.Fsm.Data = fsmData
		_ = h.bot.Fsm.Fsm.Event(update.Context(), "get-projectid")

		_, _ = bot.SendMessage(
			tu.Message(
				tu.ID(update.Message.Chat.ID),
				task.MsgAddTaskDesc,
			),
		)
	}, fsm.FsmStateEqual(h.bot.Fsm, "add-task-get-projectid"))

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
		fsmData := h.bot.Fsm.Data.(AddTaskInput)
		fsmData.Duration, _ = strconv.ParseFloat(update.Message.Text, 64)
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
		fsmData := h.bot.Fsm.Data.(AddTaskInput)
		fsmData.Date = time.Now()

		task1 := domain.Task{
			UserId:      1,
			ProjectId:   fsmData.ProjectId,
			Description: fsmData.Description,
			Duration:    fsmData.Duration,
			Date:        fsmData.Date,
			CreatedAt:   time.Now(),
		}

		_, _ = h.taskProvider.Store(update.Context(), &task1)

		_, _ = bot.SendMessage(
			tu.Message(
				tu.ID(update.Message.Chat.ID),
				task.MsgAddTaskDone,
			),
		)
	}, fsm.FsmStateEqual(h.bot.Fsm, "add-task-get-date"))
}

func (h *Handler) BuildListTasksMessage(tasks []*domain.Task) string {
	var message string

	for _, task_ := range tasks {
		message += fmt.Sprintf("%d %d %s\n", task_.ID, task_.ProjectId, task_.Description)
	}

	return task.MsgListTasks + message
}
