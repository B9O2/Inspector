package inspect

type Middleware interface {
	Run(Record) Record
}

type BaseMiddleware struct {
	f func(Record) Record
}

func (bm *BaseMiddleware) Run(r Record) Record {
	return bm.f(r)
}

func NewBaseMiddleware(f func(Record) Record) *BaseMiddleware {
	return &BaseMiddleware{
		f: f,
	}
}
