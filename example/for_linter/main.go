package main

type Color string

const (
	Red   Color = "red"
	Green Color = "green"
	Blue  Color = "blue"
)

func (c Color) Code() string {
	switch c {
	case Red:
		return "FF0000"
	case Green:
		return "00FF00"
	case Blue:
		return "0000FF"
	default:
		return "000000"
	}
}

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
