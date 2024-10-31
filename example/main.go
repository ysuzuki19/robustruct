package main

type Sample struct {
	Name string
	Age  int
}

func main() {
	s := Sample{
		Name: "John",
		Age:  30,
	}
	println(s.Name)
	println(s.Age)

	s2 := Sample{
		Name: "Alice",
	}
	println(s2.Name)
	println(s2.Age)
}
