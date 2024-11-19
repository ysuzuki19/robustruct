package result

//go:generate go run /home/yuya/Github/robustruct/cmd/senumgen/main.go
type result[T any] struct {
	ok  *T
	err error
}
