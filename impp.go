package impp

import (
	"go/ast"
	"log"
	"os"
	"strings"

	"slices"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"gopkg.in/yaml.v3"
)

const doc = "impp check imported packages"

// Analyzer is ...
var Analyzer = &analysis.Analyzer{
	Name: "impp",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

type CheckConfig struct {
	PackageNames []string `yaml:"packages"`
}

func run(pass *analysis.Pass) (any, error) {
	fileName := "impp.yaml"
	fb, err := os.ReadFile("impp.yaml")
	if err != nil {
		log.Fatalf("failed to open file: %s, err: %+v\n", fileName, err)
		return nil, err
	}

	var cfg CheckConfig
	if err := yaml.Unmarshal(fb, &cfg); err != nil {
		log.Fatalf("failed to read yaml file.\nerr:%+v\n", err)
		return nil, err
	}

	if len(cfg.PackageNames) == 0 {
		return nil, nil
	}

	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{
		(*ast.File)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		f, _ := n.(*ast.File)
		if len(f.Imports) == 0 {
			return
		}

		for _, imspec := range f.Imports {
			if slices.Contains(cfg.PackageNames, strings.Trim(imspec.Path.Value, "\"")) {
				pass.Reportf(imspec.Pos(), "%s is not allowed to be imported", imspec.Path.Value)
			}
		}
	})

	return nil, nil
}
