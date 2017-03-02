// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zhoumingjun/gog/fo"
)

// foCmd represents the gen command
var foCmd = &cobra.Command{
	Use:   "fo",
	Short: "generate functional options for a type",
	Long: `
refer to https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
`,
	Run: func(cmd *cobra.Command, args []string) {
		t, _ := cmd.Flags().GetString("type")
		if t == "" {
			fmt.Println("please set type first")
			return
		}

		var g fo.Generator
		g.ParsePackageDir(".")
		g.Generate(t)

	},
}

func init() {
	RootCmd.AddCommand(foCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// foCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// foCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	foCmd.Flags().StringP("type", "t", "", "type name")
}
