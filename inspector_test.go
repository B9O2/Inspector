package Inspect

import (
	"encoding/json"
	"fmt"
	. "github.com/B9O2/Inspector/decorators"
	colors "github.com/gookit/color"
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
	JsonEvent, _ := insp.NewType("json", func(i interface{}) string {
		marshal, err := json.MarshalIndent(i.(Event), "", "  ")
		if err != nil {
			return err.Error()
		} else {
			return string(marshal) + "\n"
		}
	})
	String, _ := insp.NewType("str", func(i interface{}) string {
		return i.(string)
	}, Red)
	insp.NewAutoType("n0p3", func() interface{} {
		return "ad6fe75dfabc"
	}, func(i interface{}) string {
		return i.(string)
	})

	MengdaColor := NewDecoration("mengda.color", func(i interface{}) interface{} {
		return colors.RGB(22, 45, 45)
	})
	insp.SetOrders("n0p3", String, EventType, JsonEvent, "_time")

	/*
		HasErr := Testing("是否包含错误", func(i interface{}) (bool, error) {
			return strings.Contains(i.(string), "err"), nil
		})
		HasOK := Testing("是否包含OK", func(i interface{}) (bool, error) {
			return strings.Contains(i.(string), "ok"), nil
		})
	*/
	Start, _ := insp.GetAutoType("_start")
	//insp.Print(EventType(e, Green), String("!!err!!!!ok!!!!", MengdaColor, HasErr, HasOK))
	insp.Print(Start(">", MengdaColor), String("hello"))
	insp.Print(String(2, Yellow), EventType(e, Green))

}
