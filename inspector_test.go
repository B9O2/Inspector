package Inspect

import (
	"fmt"
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

	insp.Print(FileName("test.txt"), EventType(e), UserName("root"))

}
