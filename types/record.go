package types

import (
	"strings"
)

type Record []*Value

func (r Record) String() string {
	return r.ToString(" ", false, false)
}

func (r Record) ToString(sep string, showLabel, decolorate bool) string {
	var parts []string
	for _, v := range r {
		part := v.ToString(decolorate)
		if showLabel {
			part += v.Label
		}
		parts = append(parts, part)
	}
	return strings.Join(parts, sep)
}
