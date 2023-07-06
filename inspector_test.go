package Inspect

import (
	"encoding/json"
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
		Name:   "危险！",
		threat: 10,
		eTime:  time.Now(),
	}
	toki := NewInspector("TOKI", 99)
	EventType, _ := toki.NewType("event", func(i interface{}) string {
		return i.(Event).ToString()
	})
	JsonEvent, _ := toki.NewType("json", func(i interface{}) string {
		marshal, err := json.MarshalIndent(i.(Event), "", "  ")
		if err != nil {
			return err.Error()
		} else {
			return string(marshal) + "\n"
		}
	})
	String, _ := toki.NewType("str", func(i interface{}) string {
		return i.(string)
	})
	toki.NewAutoType("n0p3", func() interface{} {
		return "ad6fe75dfabc"
	}, func(i interface{}) string {
		return i.(string)
	})

	toki.SetOrders("n0p3", String, EventType, JsonEvent, "_time")
	//id := toki.Record(EventType(e))
	toki.Record(EventType(e))
	toki.Print(String("hello"))
	toki.Print(String(2))

}
