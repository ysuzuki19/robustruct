//go:generate go run github.com/ysuzuki19/robustruct/cmd/gen/testdocgen -file=$GOFILE
package lib

type Struct struct {
}

// Example:
//
//	// 1.2
func Function() {
}

// Example:
//
//	// 2.2
func (s Struct) Method() {
}
