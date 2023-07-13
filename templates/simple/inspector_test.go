package simple

import (
	"errors"
	"github.com/B9O2/Inspector/decorators"
	"testing"
)

func TestInspector(t *testing.T) {
	Insp.Print(Text("Hello world!"))
	Insp.Print(Text("Hello world!", decorators.Yellow))
	Insp.Print(Path("foo.txt"), Text("Hello world!"))
	Insp.Print(Text("Alice:", decorators.Green), Text("It works!"))
	Insp.Print(LEVEL_WARNING, Text("Alice:", decorators.Green), Text("OK!"))
	Insp.Print(LEVEL_ERROR, Text("Bob:", decorators.Magenta), Error(errors.New("oh! something wrong")))
}
