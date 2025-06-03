package result

//go:generate go run ../../../main.go
type result[T any] struct {
	ok  *T
	err error
}
