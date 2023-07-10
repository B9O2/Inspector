package Inspect

import (
	"fmt"
	"github.com/B9O2/Inspector/decorators"
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
	insp := NewInspector("TEST", 99)
	EventType, _ := insp.NewType("event", func(i interface{}) string {
		return i.(Event).ToString()
	})
	FileName, _ := insp.NewType("file", func(i interface{}) string {
		return "<" + i.(string) + ">"
	})
	UserName, _ := insp.NewType("user", func(i interface{}) string {
		return "\"" + i.(string) + "\""
	})

	FlagStart := decorators.NewDecoration("flag.prefix", func(i interface{}) interface{} {
		return "flag{"
	})
	FlagEnd := decorators.NewDecoration("flag.suffix", func(i interface{}) interface{} {
		return "}end"
	})

	insp.SetTypeDecorations("_func", decorators.Invisible)
	insp.SetVisible(false)
	insp.Print(FileName("test.adwtxt"), EventType(e, FlagStart, FlagEnd, decorators.Invisible), UserName("root"))

}
