package Inspect

import (
	"fmt"
	"github.com/B9O2/Inspector/decorators"
	"github.com/B9O2/Inspector/inspect"
	"testing"
	"time"
)

type Event struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	threat int
	eTime  time.Time
}

func (e Event) ToString() string {
	return fmt.Sprintf("%s-%d-T%d %s", e.Name, e.Id, e.threat, e.eTime.Format(time.Kitchen))
}

func TestNewInspector(t *testing.T) {
	e := Event{
		Id:     1,
		Name:   "Danger！",
		threat: 10,
		eTime:  time.Now(),
	}
	alpha := inspect.NewInspector("alpha", 99)
	//beta := inspect.NewInspector("beta", 99)
	EventType, _ := alpha.NewType("event", func(i interface{}) string {
		return i.(Event).ToString()
	})
	FileName, _ := alpha.NewType("file", func(i interface{}) string {
		return "<" + i.(string) + ">"
	})
	UserName, _ := alpha.NewType("user", func(i interface{}) string {
		return "\"" + i.(string) + "\""
	})
	//alpha.SetVisibleConditions(conditions.IsGreen)
	alpha.SetSeparator("|")
	alpha.SetTypeDecorations("_func", decorators.Invisible)
	alpha.JustPrint()
	alpha.Print(UserName("1", decorators.Yellow, decorators.Green))
	alpha.Print(UserName("2", decorators.Red))
	alpha.Print(FileName("test.adwtxt"), EventType(e, decorators.Invisible), UserName("root"))
	alpha.PrintAndRecord(UserName("admin", decorators.Red))
	/*
		alpha.Range(func(record inspect.Record) bool {
			fmt.Print(record.String())
			return true
		}, conditions.IsRed, conditions.IsYellow)
	*/
}
