package inspect

type Condition func(v *Value) bool

func NewCondition(f func(v *Value) bool) Condition {
	return f
}
