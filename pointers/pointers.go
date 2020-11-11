package pointers

import (
	"fmt"
)

func TestPointer() {
	fmt.Println("Pointer test.......................")
	var t *int
	i := 100

	t = &i
	fmt.Print("t = ")
	fmt.Print(t)
	fmt.Println("\nChaning i via t")
	*t = 200
	fmt.Println(*t)
	fmt.Println(".......................")
}
