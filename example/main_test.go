package main

import "testing"

func TestMain(testing *testing.T) {
	s := Human{
		Name: "John",
	}
	println(s.Name)
	println(s.Age)
}
