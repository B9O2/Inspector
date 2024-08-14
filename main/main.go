package main

import (
	inspect "github.com/B9O2/Inspector"
	"github.com/B9O2/Inspector/decorators"
	"github.com/B9O2/Inspector/useful"
)

func main() {
	// nc, err := net.Dial("tcp4", "127.0.0.1:3344")
	// if err != nil {
	// 	fmt.Println(err)
	// }

	insp := inspect.NewInspector("main", 99)
	//insp.SetWriter(nc)
	insp.SetRecordMiddleware(useful.LevelFilter(useful.DEBUG))
	insp.SetEnable(true)
	insp.JustPrint(useful.ERROR, useful.Text("\n"))
	insp.Print(useful.INFO, useful.Text("Hello World!", decorators.Cyan))
	insp.Print(useful.WARN, useful.Text("Hello World!", decorators.Cyan))
	insp.Print(useful.ERROR, useful.Text("Hello World!", decorators.Cyan))
	insp.Print(useful.DEBUG, useful.Text("Hello World!", decorators.Cyan))
}
