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
package cmd

import (
	"nowdet/now_detector"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "nowdet",
	Short: "nowdet is a tool to detect insertion of time.Now() in arguments of a function related to Spanner",
	Long:  `nowdet is a tool to detect insertion of time.Now() in arguments of a function related to Spanner`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if o.pkgName == "" {
			err := cmd.Help()
			if err != nil {
				os.Exit(1)
			}
			os.Exit(1)
		}

		runner := now_detector.NewRunner(o.pkgName, o.funcName)
		return runner.Run()
	},
}

type option struct {
	pkgName  string
	funcName string
}

var o option

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVar(&o.pkgName, "pkg", "", "target package name")
	rootCmd.Flags().StringVar(&o.funcName, "func", "", "target function name")
}