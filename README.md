# robustruct

Go tool for struct robustness.

## Structs

### `option.Option[T]` ([README](pkg/option/README.md))

[![Go Reference](https://pkg.go.dev/badge/github.com/ysuzuki19/robustruct/pkg/option.svg)](https://pkg.go.dev/github.com/ysuzuki19/robustruct/pkg/option)

- Generic optional value type
- Provides optional value handling safety

## Generators

### testdocgen ([README](cmd/gen/testdocgen/README.md))

[![Go Reference](https://pkg.go.dev/badge/github.com/ysuzuki19/robustruct/cmd/gen/testdocgen.svg)](https://pkg.go.dev/github.com/ysuzuki19/robustruct/cmd/gen/testdocgen)

- Generates example code from test files

## Lint tool for struct

Details in [pkg/linters/README.md](pkg/linters/README.md).

- `fields_require`
  - Check all fields are initialized
  - Suggest to use zero value
- `fields_align`
  - Check all fields are aligned ordered by struct definition
  - Suggest to align fields
