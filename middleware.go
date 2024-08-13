package inspect

import "github.com/B9O2/Inspector/core"

type Middleware interface {
	Run(core.Record) core.Record
}

type BaseMiddleware struct {
	f func(core.Record) core.Record
}

func (bm *BaseMiddleware) Run(r core.Record) core.Record {
	return bm.f(r)
}

func NewBaseMiddleware(f func(core.Record) core.Record) *BaseMiddleware {
	return &BaseMiddleware{
		f: f,
	}
}
