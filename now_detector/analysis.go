package now_detector

import (
	"golang.org/x/tools/go/ssa"

	"github.com/cokeBeer/goot/pkg/dataflow/golang/switcher"
	"github.com/cokeBeer/goot/pkg/dataflow/toolkits/graph"
	"github.com/cokeBeer/goot/pkg/dataflow/toolkits/scalar"
)

const (
	detect      = "Alert"
	maybeDetect = "Warning"
	notDetect   = "Safe"
)

type NowDetectorAnalysis struct {
	*scalar.BaseFlowAnalysis
	ndaSwitcher *nowDetectorSwitcher
}

func newNowDetectorAnalysis(g *graph.UnitGraph) *NowDetectorAnalysis {
	return &NowDetectorAnalysis{
		BaseFlowAnalysis: scalar.NewBase(g),
		ndaSwitcher:      newNowDetectorSwitcher(),
	}
}

func (a *NowDetectorAnalysis) NewInitialFlow() *map[any]any {
	m := make(map[any]any)
	for _, v := range a.Graph.Func.Params {
		m[v] = notDetect
	}
	return &m
}

func (a *NowDetectorAnalysis) FlowThrough(inMap *map[any]any, unit ssa.Instruction, outMap *map[any]any) {
	a.Copy(inMap, outMap)
	a.ndaSwitcher.inMap = inMap
	a.ndaSwitcher.outMap = outMap
	switcher.Apply(a.ndaSwitcher, unit)
}
