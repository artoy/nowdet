package now_detector

import (
	"github.com/cokeBeer/goot/pkg/dataflow/golang/switcher"
	"golang.org/x/tools/go/ssa"
)

type nowDetectorSwitcher struct {
	switcher.BaseSwitcher
	inMap  *map[any]any
	outMap *map[any]any
}

func newNowDetectorSwitcher() *nowDetectorSwitcher {
	return &nowDetectorSwitcher{}
}

func (s *nowDetectorSwitcher) CaseCall(inst *ssa.Call) {
	if inst.Name() == "Now" {
		(*s.outMap)[inst] = detect
	}
}
