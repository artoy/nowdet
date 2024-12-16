/*
Copyright © 2024 Ryota Kobayashi

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package now_detector

import (
	"os"

	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa/ssautil"

	"github.com/cokeBeer/goot/pkg/dataflow/toolkits/graph"
	"github.com/cokeBeer/goot/pkg/dataflow/toolkits/solver"
)

type Runner struct {
	pkgPath  string
	funcName string
	isDebug  bool
}

func NewRunner(pkg string, fn string, isDebug bool) *Runner {
	return &Runner{
		pkgPath:  pkg,
		funcName: fn,
		isDebug:  isDebug,
	}
}

func (r *Runner) Run() error {
	cfg := packages.Config{Mode: packages.LoadAllSyntax}
	initial, err := packages.Load(&cfg, r.pkgPath)
	if err != nil {
		return err
	}

	prog, pkgs := ssautil.AllPackages(initial, 0)
	prog.Build()

	for _, p := range pkgs {
		fn := p.Func(r.funcName)

		if r.isDebug {
			_, err = fn.WriteTo(os.Stdout)
			if err != nil {
				return err
			}
		}

		g := graph.New(fn)
		analysis := newNowDetectorAnalysis(g, r.isDebug)
		solver.Solve(analysis, true)
	}

	return nil
}
