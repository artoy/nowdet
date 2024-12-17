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
	"fmt"

	"golang.org/x/tools/go/ssa"

	"github.com/cokeBeer/goot/pkg/dataflow/golang/switcher"
	"github.com/cokeBeer/goot/pkg/dataflow/toolkits/graph"
	"github.com/cokeBeer/goot/pkg/dataflow/toolkits/scalar"
	"github.com/cokeBeer/goot/pkg/dataflow/util/entry"
)

const (
	none = iota
	now
	detected
)

type NowDetectorAnalysis struct {
	*scalar.BaseFlowAnalysis
	ndaSwitcher *nowDetectorSwitcher
	isDebug     bool
}

func newNowDetectorAnalysis(g *graph.UnitGraph, isDebug bool) *NowDetectorAnalysis {
	return &NowDetectorAnalysis{
		BaseFlowAnalysis: scalar.NewBase(g),
		ndaSwitcher:      newNowDetectorSwitcher(),
		isDebug:          isDebug,
	}
}

func (a *NowDetectorAnalysis) NewInitalFlow() *map[any]any {
	m := make(map[any]any)
	for _, v := range a.Graph.Func.Params {
		m[v] = none
	}

	return &m
}

func (a *NowDetectorAnalysis) FlowThrougth(inMap *map[any]any, unit ssa.Instruction, outMap *map[any]any) {
	a.Copy(inMap, outMap)
	a.ndaSwitcher.inMap = inMap
	a.ndaSwitcher.outMap = outMap
	switcher.Apply(a.ndaSwitcher, unit)
}

func (a *NowDetectorAnalysis) End(universe []*entry.Entry) {
	for _, e := range universe {
		if a.isDebug {
			fmt.Println(e.InFlow)
			fmt.Println(e.Data)
			fmt.Println(e.OutFlow)
			fmt.Println()
		}

		for _, v := range *e.OutFlow {
			if v == detected {
				fmt.Println(fmt.Sprintf("time.Now() is detected in %s", e.Data))
				return
			}
		}
	}
}
