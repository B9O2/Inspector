package core

import (
	"strings"

	"github.com/gookit/color"
)

type Value struct {
	Label string
	Text  string
	Color color.Style
}

func (vs Value) String() string {
	return vs.ToString(false)
}

func (vs Value) ToString(decolorate bool) string {
	if decolorate {
		return vs.Text
	}
	var lines []string
	for _, line := range strings.Split(vs.Text, "\n") {
		lines = append(lines, vs.Color.Sprint(line))
	}
	
	return strings.Join(lines, "\n")
}
