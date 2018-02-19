// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
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

	"github.com/nlnwa/veidemannctl/bindata"
	"github.com/spf13/cobra"
	"runtime"
	"strings"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Get version information",
	Long:  `Get version information.`,
	Run: func(cmd *cobra.Command, args []string) {
		Version, err := bindata.Asset("version")
		if err != nil {
			panic(err)
		}

		fmt.Printf("Version: %s, Go version: %s, Go OS/ARCH: %s %s\n",
			strings.Trim(string(Version), "\n\r "),
			runtime.Version(),
			runtime.GOOS,
			runtime.GOARCH)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
