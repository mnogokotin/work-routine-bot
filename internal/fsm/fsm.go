package fsm

import (
	"github.com/looplab/fsm"
	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
	"reflect"
	"strings"
	"work-routine-bot/internal/bot"
)

func NewConvFsm(name string, abc interface{}) *fsm.FSM {
	var events fsm.Events
	var nextFieldName string
	t := reflect.TypeOf(abc)
	for i := 0; i < t.NumField(); i++ {
		fieldName := strings.ToLower(t.Field(i).Name)
		if i+1 < t.NumField() {
			nextFieldName = strings.ToLower(t.Field(i + 1).Name)
		} else {
			nextFieldName = "end"
		}
		event := fsm.EventDesc{
			Name: "get-" + fieldName,
			Src:  []string{name + "-get-" + fieldName},
			Dst:  name + "-get-" + nextFieldName,
		}
		events = append(events, event)
	}

	return fsm.NewFSM(
		name+"-get-"+strings.ToLower(t.Field(0).Name),
		events,
		fsm.Callbacks{},
	)

}

func FsmStateEqual(fsm *bot.Fsm, fsmState string) telegohandler.Predicate {
	return func(_ telego.Update) bool {
		if fsm.Fsm == nil {
			return false
		}

		return fsm.Fsm.Current() == fsmState
	}
}
