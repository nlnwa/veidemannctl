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

// setServerNameOverrideCmd represents the set-server-name-override command
var setServerNameOverrideCmd = &cobra.Command{
	Use:   "set-server-name-override HOST",
	Short: "Sets the server name override",
	Long: `Sets the server name override.

Use this when there is a mismatch between exposed server name for the cluster and the certificate. The use is a security
risk and is only recommended for testing.

Examples:
  # Sets the server name override to test.local
  veidemannctl config set-server-name-override test.local
`,
	Aliases: []string{"set-servername"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			viper.Set("serverNameOverride", args[0])
			configutil.WriteConfig()
		} else {
			cmd.Usage()
		}
	},
}

// getServerNameOverrideCmd represents the get-server-name-override command
var getServerNameOverrideCmd = &cobra.Command{
	Use:   "get-server-name-override HOST",
	Short: "Displays the server name override",
	Long: `Displays the server name override

Examples:
  # Display the server name override
  veidemannctl config get-server-name-override
`,
	Aliases: []string{"get-servername"},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Server name override: %s\n", viper.Get("serverNameOverride"))
	},
}

func init() {
	ConfigCmd.AddCommand(setServerNameOverrideCmd)
	ConfigCmd.AddCommand(getServerNameOverrideCmd)
}
