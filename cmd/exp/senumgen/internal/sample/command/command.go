package command

import (
	"github.com/ysuzuki19/robustruct/cmd/exp/senumgen/command/sub"
	"github.com/ysuzuki19/robustruct/cmd/exp/senumgen/types"
)

//go:generate go run /home/yuya/Github/robustruct/cmd/exp/senumgen/main.go
type command struct {
	help  types.NonVar
	run   *string
	sub   *sub.SubCommand
	local *Local
}

type Local struct {
	Dummy int
}
