package composition

type guardable interface {
	guard()
}

type dog struct {
}

type tiger struct {
}

type cat struct {
}

func (d dog) guard() {
	println("Dog guards...")
}

func (t tiger) guard() {
	println("Tiger guards...")
}

func guard(g guardable) {
	g.guard()
}

func main() {
	println("We inherit an interface by implementing all its methods.........")
	t := tiger{}
	d := dog{}
	// c := cat{}
	guard(t)
	guard(d)
	// guard(c) //does not compile cat does not inherit guardable
}
