/*
 * Copyright 2019 National Library of Norway.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *       http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package config

import (
	"fmt"
	"github.com/nlnwa/veidemannctl/src/configutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// setApiKeyCmd represents the set-apikey command
var setApiKeyCmd = &cobra.Command{
	Use:   "set-apikey API_KEY",
	Short: "Sets the api-key to use for authentication",
	Long: `Sets the api-key to use for authentication

Examples:
  # Set the api-key to use for authentication to Veidemann controller service to myKey
  veidemannctl config set-apikey myKey
`,
	Aliases: []string{"set-apikey", "set-api-key"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			viper.Set("apiKey", args[0])
			configutil.WriteConfig()
		} else {
			cmd.Usage()
		}
	},
}

// getApiKeyCmd represents the get-apikey command
var getApiKeyCmd = &cobra.Command{
	Use:   "get-apikey",
	Short: "Displays Veidemann authentication api-key",
	Long: `Displays Veidemann authentication api-key

Examples:
  # Display Veidemann authentication api-key
  veidemannctl config get-apikey
`,
	Aliases: []string{"get-apikey", "apikey", "get-api-key", "api-key"},
	Run: func(cmd *cobra.Command, args []string) {
		apiKey := viper.Get("apiKey")
		if apiKey == nil {
			fmt.Printf("No api-key configured\n")
		} else {
			fmt.Printf("api-key: %s\n", apiKey)
		}
	},
}

func init() {
	ConfigCmd.AddCommand(setApiKeyCmd)
	ConfigCmd.AddCommand(getApiKeyCmd)
}
