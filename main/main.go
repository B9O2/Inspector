package main

import (
	inspect "github.com/B9O2/Inspector"
	"github.com/B9O2/Inspector/decorators"
	"github.com/B9O2/Inspector/value_types"
)

func main() {
	// nc, err := net.Dial("tcp4", "127.0.0.1:3344")
	// if err != nil {
	// 	fmt.Println(err)
	// }

	insp := inspect.NewInspector("main", 99)
	//insp.SetWriter(nc)
	insp.Print(value_types.Text("Hello World!", decorators.Cyan))
}
