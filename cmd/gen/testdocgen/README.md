# cmd/gen/testdocgen

testdocgen generates example code from test files.

## Usage

add following line to the target Go source file.

```go
//go:generate go run github.com/ysuzuki19/robustruct/cmd/gen/testdocgen -file=$GOFILE
package lib

type Struct struct {
}

func Function() {
}

func (s Struct) Method() {
}
```

add `<target_file_name>_test.go` file to the same directory (must to be same file name as the target Go source file).

```go
package lib_test

func TestFunction(t *testing.T) {
  // 1.1
  // testdoc begin Function
  // 1.2
  // testdoc end
  // 1.3
}

func TestStruct_Method(t *testing.T) {
  // 2.1
  // testdoc begin Struct.Method
  // 2.2
  // testdoc end
  // 2.3
}
```

then run `go generate` command.

`testdocgen` will generate (update) example code in the target Go source file.

```go
//go:generate go run github.com/ysuzuki19/robustruct/cmd/gen/testdocgen -file=$GOFILE
package lib

type Struct struct {
}

// Example:
//
// 1.2
func Function() {
}

// Example:
//
// 2.2
func (s Struct) Method() {
}
```
