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
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"path/filepath"
)

// useContextCmd represents the use-context command
var useContextCmd = &cobra.Command{
	Use:   "use-context CONTEXT_NAME",
	Short: "Sets the current-context",
	Long: `Sets the current-context

Examples:
  # Use the context for the prod cluster
  veidemannctl config use-context prod
`,
	Aliases: []string{"use"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			ok, err := configutil.ContextExists(args[0])
			if err != nil {
				log.Fatalf("Failed switching context to %v. Cause: %v", args[0], err)
			}
			if !ok {
				fmt.Printf("Non existing context '%v'\n", args[0])
				return
			}
			if err := configutil.SetCurrentContext(args[0]); err != nil {
				log.Fatalf("Failed switching context to %v. Cause: %v", args[0], err)
			}
			fmt.Printf("Switched to context: '%v'\n", args[0])
		} else {
			fmt.Println("Missing context name")
			cmd.Usage()
		}
	},
}

// createContextCmd represents the create-context command
var createContextCmd = &cobra.Command{
	Use:   "create-context CONTEXT_NAME",
	Short: "Creates a new context",
	Long: `Creates a new context

Examples:
  # Create context for the prod cluster
  veidemannctl config create-context prod
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			ok, err := configutil.ContextExists(args[0])
			if err != nil {
				log.Fatalf("Failed creating context %v. Cause: %v", args[0], err)
			}
			if ok {
				fmt.Printf("Context '%v' already exists\n", args[0])
				return
			}

			contextDir := configutil.GetConfigDir("contexts")
			viper.SetConfigFile(filepath.Join(contextDir, args[0]+".yaml"))
			configutil.WriteConfig()

			if err := configutil.SetCurrentContext(args[0]); err != nil {
				log.Fatalf("Failed creating context %v. Cause: %v", args[0], err)
			}
			fmt.Printf("Created context: '%v'\n", args[0])
		} else {
			fmt.Println("Missing context name")
			cmd.Usage()
		}
	},
}

// currentContextCmd represents the current-context command
var currentContextCmd = &cobra.Command{
	Use:   "current-context",
	Short: "Displays the current-context",
	Long: `Displays the current-context

Examples:
  # Display the current context
  veidemannctl config current-context
`,
	Aliases: []string{"context"},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Using context:", configutil.GlobalFlags.Context)
	},
}

// listContextsCmd represents the list-contexts command
var listContextsCmd = &cobra.Command{
	Use:   "list-contexts",
	Short: "Displays the known contexts",
	Long: `Displays the known contexts.

Examples:
  # Get a list of known contexts
  veidemannctl config list-contexts
`,
	Aliases: []string{"contexts"},
	Run: func(cmd *cobra.Command, args []string) {
		cs, err := configutil.ListContexts()
		if err != nil {
			log.Fatalf("Failed listing contexts to %v. Cause: %v", args[0], err)
		}
		fmt.Println("Known contexts:")
		for _, c := range cs {
			fmt.Println(c)
		}
	},
}

func init() {
	ConfigCmd.AddCommand(useContextCmd)
	ConfigCmd.AddCommand(createContextCmd)
	ConfigCmd.AddCommand(currentContextCmd)
	ConfigCmd.AddCommand(listContextsCmd)
}
