package main

type Human struct {
	Name string
	Age  int
}

func main() {
	s := Human{
		Name: "John",
		Age:  30,
	}
	println(s.Name)
	println(s.Age)
}
