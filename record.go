package Inspect

import (
	"fmt"
	"strings"
)

type Value struct {
	typeLabel string
	formatter func(interface{}) string
	data      interface{}
}

func (v Value) Label() string {
	return v.typeLabel
}

func (v Value) Data() interface{} {
	return v.data
}

type Record []*Value

func (r Record) String() string {
	var parts []string
	for _, v := range r {
		if v == nil {
			continue
		}
		if v.formatter != nil {
			parts = append(parts, v.formatter(v.data))
		} else {
			parts = append(parts, fmt.Sprintf("%v", v.data))
		}
	}
	return strings.Join(parts, " ")
}

func (r Record) ToString(sep string) string {
	var parts []string
	for _, v := range r {
		if v == nil {
			continue
		}
		if v.formatter != nil {
			parts = append(parts, v.formatter(v.data))
		} else {
			parts = append(parts, fmt.Sprintf("%v", v.data))
		}
	}
	return strings.Join(parts, sep)
}

func (r Record) StringWithVType(sep string) string {
	var parts []string
	part := ""
	for _, v := range r {
		if v == nil {
			continue
		}
		if v.formatter != nil {
			part = fmt.Sprintf("%v", v.data)
		} else {
			part = v.formatter(v.data)
		}
		part += "<" + v.typeLabel + ">"
		parts = append(parts, part)
	}
	return strings.Join(parts, sep)
}
