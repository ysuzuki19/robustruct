# robustruct

Go Lint tool for struct

This tool is module plugin for `golangci-lint`.

# Features

- `fields_require`
  - Check all fields are initialized
  - Suggest to use zero value
- `fields_align`
  - Check all fields are aligned ordered by struct definition
  - Suggest to align fields

# Disabling

You can disable a check with a comment.

ignore all features.

```go
// ignore:robustruct
```

ignore `fields_require`.

```go
// ignore:fields_require
```

ignore `fields_align`.

```go
// ignore:fields_align
```

# Sample

for the following code, robustruct will suggest the following fixes.

```go
package main

type Sample struct {
    A int
    B string
    AA bool
}
```

## fields_require

Before fix

```go
func main() {
    s := Sample{
        A: 1,
        B: "hello",
    }
}
```

After fix

```go
func main() {
    s := Sample{
        A: 1,
        B: "hello",
        AA: false,
    }
}
```

## fields_align

Before fix

```go
func main() {
    s := Sample{
        B: "hello",
        A: 1,
        AA: false,
    }
}
```

After fix

```go
func main() {
    s := Sample{
        A: 1,
        B: "hello",
        AA: false,
    }
}
```
