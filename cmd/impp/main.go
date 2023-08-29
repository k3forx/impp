package main

import (
	"github.com/k3forx/impp"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(impp.Analyzer) }
