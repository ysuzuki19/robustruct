package main

import (
	"github.com/ysuzuki19/robustruct"
	"golang.org/x/tools/go/analysis"
)

func New(conf any) ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{robustruct.Analyzer}, nil
}
