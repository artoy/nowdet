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
	"testing"

	"golang.org/x/tools/go/ssa"

	"github.com/cokeBeer/goot/pkg/dataflow/golang/switcher"
	"github.com/google/go-cmp/cmp"

	"github.com/artoy/nowdet/now_detector/testutil"
)

func Test_nowDetectorSwitcher_CaseCall(t *testing.T) {
	inMap := make(map[any]any)
	outMap := make(map[any]any)

	type fields struct {
		BaseSwitcher switcher.BaseSwitcher
		inMap        *map[any]any
		outMap       *map[any]any
	}
	type args struct {
		inst *ssa.Call
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Detect",
			fields: fields{
				BaseSwitcher: switcher.BaseSwitcher{},
				inMap:        &inMap,
				outMap:       &outMap,
			},
			args: args{
				inst: testutil.GetCall(t),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &nowDetectorSwitcher{
				BaseSwitcher: tt.fields.BaseSwitcher,
				inMap:        tt.fields.inMap,
				outMap:       tt.fields.outMap,
			}
			s.CaseCall(tt.args.inst)
			expected := now
			actual := (*s.outMap)[tt.args.inst.Name()]
			if diff := cmp.Diff(expected, actual); diff != "" {
				t.Errorf("(-want, +got)\n%s", diff)
			}
		})
	}
}

func Test_nowDetectorSwitcher_CasePhi(t *testing.T) {
	inMap := make(map[any]any)
	inMap["t0"] = now

	outMap := make(map[any]any)
	outMap["t0"] = now

	type fields struct {
		BaseSwitcher switcher.BaseSwitcher
		inMap        *map[any]any
		outMap       *map[any]any
	}
	type args struct {
		inst *ssa.Phi
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Detect",
			fields: fields{
				BaseSwitcher: switcher.BaseSwitcher{},
				inMap:        &inMap,
				outMap:       &outMap,
			},
			args: args{
				inst: testutil.GetPhi(t),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &nowDetectorSwitcher{
				BaseSwitcher: tt.fields.BaseSwitcher,
				inMap:        tt.fields.inMap,
				outMap:       tt.fields.outMap,
			}
			s.CasePhi(tt.args.inst)
			expected := now
			actual := (*s.outMap)[tt.args.inst.Name()]
			if diff := cmp.Diff(expected, actual); diff != "" {
				t.Errorf("(-want, +got)\n%s", diff)
			}
		})
	}
}

func Test_nowDetectorSwitcher_CaseMakeInterface(t *testing.T) {
	inMap := make(map[any]any)
	inMap["t4"] = now

	outMap := make(map[any]any)
	outMap["t4"] = now

	type fields struct {
		BaseSwitcher switcher.BaseSwitcher
		inMap        *map[any]any
		outMap       *map[any]any
	}
	type args struct {
		inst *ssa.MakeInterface
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Detect",
			fields: fields{
				BaseSwitcher: switcher.BaseSwitcher{},
				inMap:        &inMap,
				outMap:       &outMap,
			},
			args: args{
				inst: testutil.GetMakeInterface(t),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &nowDetectorSwitcher{
				BaseSwitcher: tt.fields.BaseSwitcher,
				inMap:        tt.fields.inMap,
				outMap:       tt.fields.outMap,
			}
			s.CaseMakeInterface(tt.args.inst)
			expected := now
			actual := (*s.outMap)[tt.args.inst.Name()]
			if diff := cmp.Diff(expected, actual); diff != "" {
				t.Errorf("(-want, +got)\n%s", diff)
			}
		})
	}
}

func Test_nowDetectorSwitcher_CaseIndexAddr(t *testing.T) {
	inMap := make(map[any]any)
	outMap := make(map[any]any)

	type fields struct {
		BaseSwitcher switcher.BaseSwitcher
		inMap        *map[any]any
		outMap       *map[any]any
	}
	type args struct {
		inst *ssa.IndexAddr
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Detect",
			fields: fields{
				BaseSwitcher: switcher.BaseSwitcher{},
				inMap:        &inMap,
				outMap:       &outMap,
			},
			args: args{
				inst: testutil.GetIndexAddr(t),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &nowDetectorSwitcher{
				BaseSwitcher: tt.fields.BaseSwitcher,
				inMap:        tt.fields.inMap,
				outMap:       tt.fields.outMap,
			}
			s.CaseIndexAddr(tt.args.inst)
			expected := tt.args.inst.Name()
			actual := (*s.outMap)["t0"]
			if diff := cmp.Diff(expected, actual); diff != "" {
				t.Errorf("(-want, +got)\n%s", diff)
			}
		})
	}
}

func Test_nowDetectorSwitcher_CaseSlice(t *testing.T) {
	inMap := make(map[any]any)
	inMap["t0"] = now
	outMap := make(map[any]any)
	outMap["t0"] = now

	type fields struct {
		BaseSwitcher switcher.BaseSwitcher
		inMap        *map[any]any
		outMap       *map[any]any
	}
	type args struct {
		inst *ssa.Slice
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Detect",
			fields: fields{
				BaseSwitcher: switcher.BaseSwitcher{},
				inMap:        &inMap,
				outMap:       &outMap,
			},
			args: args{
				inst: testutil.GetSlice(t),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &nowDetectorSwitcher{
				BaseSwitcher: tt.fields.BaseSwitcher,
				inMap:        tt.fields.inMap,
				outMap:       tt.fields.outMap,
			}
			s.CaseSlice(tt.args.inst)
			expected := now
			actual := (*s.outMap)[tt.args.inst.Name()]
			if diff := cmp.Diff(expected, actual); diff != "" {
				t.Errorf("(-want, +got)\n%s", diff)
			}
		})
	}
}

func Test_nowDetectorSwitcher_CaseStore(t *testing.T) {
	inMap := make(map[any]any)
	outMap := make(map[any]any)

	type fields struct {
		BaseSwitcher switcher.BaseSwitcher
		inMap        *map[any]any
		outMap       *map[any]any
	}
	type args struct {
		inst *ssa.Store
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Detect",
			fields: fields{
				BaseSwitcher: switcher.BaseSwitcher{},
				inMap:        &inMap,
				outMap:       &outMap,
			},
			args: args{
				inst: testutil.GetStore(t),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &nowDetectorSwitcher{
				BaseSwitcher: tt.fields.BaseSwitcher,
				inMap:        tt.fields.inMap,
				outMap:       tt.fields.outMap,
			}
			s.CaseStore(tt.args.inst)
			expected := "t2"
			actual := (*s.outMap)[tt.args.inst.Addr.Name()]
			if diff := cmp.Diff(expected, actual); diff != "" {
				t.Errorf("(-want, +got)\n%s", diff)
			}
		})
	}
}
