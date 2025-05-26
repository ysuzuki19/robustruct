package main

import "github.com/ysuzuki19/robustruct/pkg/option"

func main() {
	s := struct {
		Name string
		Age  option.Option[int]
	}{
		Name: "John",
		Age:  option.NewSome(30),
	}
	println(s.Name)
	if age, ok := s.Age.Get(); ok {
		println(*age)
	} else {
		println("Age is not set")
	}
}
