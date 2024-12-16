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
	"log"
	"slices"

	"golang.org/x/tools/go/ssa"

	"github.com/cokeBeer/goot/pkg/dataflow/golang/switcher"
)

const maxRetry = 30

var targetSpannerFuncs = []string{"Insert", "Update"}

type nowDetectorSwitcher struct {
	switcher.BaseSwitcher
	inMap  *map[any]any
	outMap *map[any]any
}

func newNowDetectorSwitcher() *nowDetectorSwitcher {
	return &nowDetectorSwitcher{}
}

func (s *nowDetectorSwitcher) CaseCall(inst *ssa.Call) {
	fn := inst.Common().StaticCallee()

	if fn.Pkg.Pkg.Path() == "time" && fn.Name() == "Now" {
		(*s.outMap)[inst.Name()] = now
		return
	}

	if fn.Pkg.Pkg.Path() == "cloud.google.com/go/spanner" && (slices.Contains(targetSpannerFuncs, fn.Name())) {
		for _, arg := range inst.Common().Args {
			if propagate(s.inMap, arg.Name()) == now {
				(*s.outMap)[inst.Name()] = detected
				return
			}
		}
	}
}

func (s *nowDetectorSwitcher) CasePhi(inst *ssa.Phi) {
	for _, v := range inst.Edges {
		// TODO: We have to merge pointer information we represent as a graph from previous blocks.
		//       Moreover, we have to merge them in a way that there is a fix point.
		if propagate(s.inMap, v.Name()) == now {
			(*s.outMap)[inst.Name()] = now
			return
		}
	}
}

func (s *nowDetectorSwitcher) CaseMakeInterface(inst *ssa.MakeInterface) {
	if propagate(s.inMap, inst.X.Name()) == now {
		(*s.outMap)[inst.Name()] = now
		return
	}
}

func (s *nowDetectorSwitcher) CaseIndexAddr(inst *ssa.IndexAddr) {
	// make edge from an element to a collection to propagate the information
	(*s.outMap)[inst.X.Name()] = inst.Name()
}

func (s *nowDetectorSwitcher) CaseSlice(inst *ssa.Slice) {
	if propagate(s.inMap, inst.X.Name()) == now {
		(*s.outMap)[inst.Name()] = now
		return
	}
}

func (s *nowDetectorSwitcher) CaseStore(inst *ssa.Store) {
	(*s.outMap)[inst.Addr.Name()] = inst.Val.Name()
}

// TODO: use memo instead of max
func propagate(m *map[any]any, from string) int {
	for i := 0; i < maxRetry; i++ {
		switch t := (*m)[from].(type) {
		case int:
			return t
		case string:
			from = t
		case nil:
			return none
		default:
			panic(fmt.Sprintf("unexpected type in data flow information: %v", t))
		}
	}

	log.Print("propagation achieved max retry")
	return none
}
