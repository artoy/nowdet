package now_detector

import (
	"log"

	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"

	"github.com/cokeBeer/goot/pkg/dataflow/toolkits/graph"
	"github.com/cokeBeer/goot/pkg/dataflow/toolkits/solver"
)

type Runner struct {
	Src      string
	Function string
}

func NewRunner(src, function string) *Runner {
	return &Runner{src, function}
}

func (r *Runner) Run() {
	// Generate AST
	fSet := token.NewFileSet()
	f, err := parser.ParseFile(fSet, "", r.Src, parser.Mode(0))
	if err != nil {
		log.Println(err)
	}
	files := []*ast.File{f}

	// Build package
	pkg := types.NewPackage("now_detector", "")
	srcPkg, _, err := ssautil.BuildPackage(
		&types.Config{Importer: importer.Default()},
		fSet,
		pkg,
		files,
		ssa.SanityCheckFunctions,
	)
	if err != nil {
		log.Println(err)
	}

	// Build graph
	g := graph.New(srcPkg.Func(r.Function))

	// Build analysis
	analysis := newNowDetectorAnalysis(g)

	// Solve analysis
	solver.Solve(analysis, true)
}
