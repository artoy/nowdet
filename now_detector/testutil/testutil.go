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
package testutil

import (
	"fmt"
	"testing"

	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

const program = `
package main

import (
	"time"
	"fmt"
)

func main() {
	var x time.Time
	if true {
		x = time.Now()
	}
	fmt.Println(x)
	return
}`

const program2 = `
package main

import "fmt"

func main() {
	s := fmt.Sprintf("test: %s", "test")
	i := []interface{}{"test", s}
	fmt.Println(i)
}
`

func processProg(t *testing.T, p string) *ssa.Package {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "main.go", p, parser.Mode(0))
	if err != nil {
		t.Fatal(err)
	}
	files := []*ast.File{f}

	pkg := types.NewPackage("test", "")

	ssaPkg, _, err := ssautil.BuildPackage(
		&types.Config{Importer: importer.Default()},
		fset,
		pkg,
		files,
		ssa.SanityCheckFunctions,
	)
	if err != nil {
		t.Fatal(err)
	}

	return ssaPkg
}

func GetCall(t *testing.T) *ssa.Call {
	pkg := processProg(t, program)
	fn := pkg.Func("main")
	call, ok := fn.Blocks[1].Instrs[0].(*ssa.Call)
	if !ok {
		t.Fatal(fmt.Sprintf("%t is not a call", fn.Blocks[1].Instrs[0].String()))
	}
	return call
}

func GetPhi(t *testing.T) *ssa.Phi {
	pkg := processProg(t, program)
	fn := pkg.Func("main")
	phi, ok := fn.Blocks[2].Instrs[0].(*ssa.Phi)
	if !ok {
		t.Fatal(fmt.Sprintf("%t is not a phi", fn.Blocks[2].Instrs[0].String()))
	}
	return phi
}

func GetMakeInterface(t *testing.T) *ssa.MakeInterface {
	pkg := processProg(t, program2)
	fn := pkg.Func("main")
	makeInterface, ok := fn.Blocks[0].Instrs[11].(*ssa.MakeInterface)
	if !ok {
		t.Fatal(fmt.Sprintf("%t is not a makeInterface", fn.Blocks[0].Instrs[2].String()))
	}
	return makeInterface
}

func GetIndexAddr(t *testing.T) *ssa.IndexAddr {
	pkg := processProg(t, program2)
	fn := pkg.Func("main")
	indexAddr, ok := fn.Blocks[0].Instrs[1].(*ssa.IndexAddr)
	if !ok {
		t.Fatal(fmt.Sprintf("%t is not a indexAddr", fn.Blocks[0].Instrs[2].String()))
	}
	return indexAddr
}

func GetSlice(t *testing.T) *ssa.Slice {
	pkg := processProg(t, program2)
	fn := pkg.Func("main")
	slice, ok := fn.Blocks[0].Instrs[4].(*ssa.Slice)
	if !ok {
		t.Fatal(fmt.Sprintf("%t is not a slice", fn.Blocks[0].Instrs[2].String()))
	}
	return slice
}

func GetStore(t *testing.T) *ssa.Store {
	pkg := processProg(t, program2)
	fn := pkg.Func("main")
	store, ok := fn.Blocks[0].Instrs[3].(*ssa.Store)
	if !ok {
		t.Fatal(fmt.Sprintf("%t is not a store", fn.Blocks[0].Instrs[2].String()))
	}
	return store
}
