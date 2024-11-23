package result

//go:generate go run /home/yuya/Github/robustruct/cmd/exp/variantgen/main.go
type result[T any] struct {
	ok  *T
	err error
}
