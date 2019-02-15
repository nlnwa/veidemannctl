// Copyright © 2017 National Library of Norway
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

package config

import (
	"fmt"
	"github.com/nlnwa/veidemannctl/src/configutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// setAddrCmd represents the report command
var setAddrCmd = &cobra.Command{
	Use:   "set-address HOST:PORT",
	Short: "Sets the address to Veidemann controller service",
	Long: `Sets the address to Veidemann controller service

Examples:
  # Sets the address to Veidemann controller service to localhost:50051
  veidemannctl config set-address localhost:50051
`,
	Aliases: []string{"set-addr", "set-controller"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			viper.Set("controllerAddress", args[0])
			configutil.WriteConfig()
		} else {
			cmd.Usage()
		}
	},
}

// useContextCmd represents the report command
var getAddrCmd = &cobra.Command{
	Use:   "get-address",
	Short: "Displays Veidemann controller service address",
	Long: `Displays Veidemann controller service address

Examples:
  # Display Veidemann controller service address
  veidemannctl config get-address
`,
	Aliases: []string{"get-addr", "address", "controller"},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Controller address: %s\n", viper.Get("controllerAddress"))
	},
}

func init() {
	ConfigCmd.AddCommand(setAddrCmd)
	ConfigCmd.AddCommand(getAddrCmd)
}
