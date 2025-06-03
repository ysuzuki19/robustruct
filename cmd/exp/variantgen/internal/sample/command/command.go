package command

import (
	"github.com/ysuzuki19/robustruct/cmd/exp/variantgen/internal/sample/command/sub"
	"github.com/ysuzuki19/robustruct/cmd/exp/variantgen/types"
)

//go:generate go run ../../../main.go
type command struct {
	help  types.NonVar
	run   *string
	sub   *sub.SubCommand
	local *Local
}

type Local struct {
	Dummy int
}
