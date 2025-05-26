# robustruct / options

`option` package provides a generic optional value type `Option[T]` for Go.

it is designed to handle optional values safely, avoiding common pitfalls associated with nil pointers and uninitialized values.

## Example

```go
package main

import "github.com/ysuzuki19/robustruct/pkg/option"

type User struct {
	Name string
	Age  option.Option[int]
}

func main() {
	u := User{
		Name: "John",
		Age:  option.NewSome(30),
	}
	println(u.Name)
	if age, ok := u.Age.Get(); ok {
		println(*age)
	} else {
		println("Age is not set")
	}
}
```
