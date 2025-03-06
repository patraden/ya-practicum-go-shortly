package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestIdentifyExitAlgo(t *testing.T) {
	t.Parallel()

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "./testdata/main.go", nil, 0)
	require.NoError(t, err)

	assert.Equal(t, "main", f.Name.Name)

	ast.Inspect(f, func(n ast.Node) bool {
		callExpr, ok := n.(*ast.CallExpr)

		if !ok {
			return true
		}

		selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		pkgIdent, ok := selExpr.X.(*ast.Ident)
		if !ok || pkgIdent.Name != "os" || selExpr.Sel.Name != "Exit" {
			return true
		}

		pos := fset.Position(callExpr.Pos())
		t.Logf("direct call to os.Exit found in %s:%d, at pos: %d", pos.Filename, pos.Line, callExpr.Pos())

		return false
	})
}

func TestExitCheckAnalyzer(t *testing.T) {
	t.Parallel()

	analysistest.Run(t, analysistest.TestData(), ExitCheckAnalyzer, "./...")
}
