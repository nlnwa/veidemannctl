// Copyright © 2020 National Library of Norway
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
	"context"
	"fmt"
	"github.com/nlnwa/veidemann-api/go/config/v1"
	"github.com/nlnwa/veidemannctl/src/connection"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// scriptParametersCmd represents the script-parameters command
var scriptParametersCmd = &cobra.Command{
	Use:   "script-parameters CRAWLJOB_CONFIG_ID [SEED_ID]",
	Short: "Get the active script parameters for a Crawl Job",
	Long: `Get the active script parameters for a Crawl Job

Examples:
  # See active script parameters for a Crawl Job
  veidemannctl script-parameters 5604f0cc-315d-4091-8d6e-1b17a7eb990b

  # See active script parameters for a Crawl Job and eventual overrides from Seed and Entity
  veidemannctl script-parameters 5604f0cc-315d-4091-8d6e-1b17a7eb990b 9f89ca44-afe0-4f8f-808f-9df1a0fe64c9
`,
	Aliases: []string{"params"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Missing CrawlJobConfig id")
			cmd.Usage()
		}
		if len(args) <= 2 {
			configClient, conn := connection.NewConfigClient()
			defer conn.Close()

			request := &config.GetScriptAnnotationsRequest{
				Job: &config.ConfigRef{Kind: config.Kind_crawlJob, Id: args[0]},
			}
			if len(args) == 2 {
				request.Seed = &config.ConfigRef{Kind: config.Kind_seed, Id: args[1]}
			}

			response, err := configClient.GetScriptAnnotations(context.Background(), request)
			if err != nil {
				log.Fatalf("Failed getting parameters for %v. Cause: %v", args[0], err)
			}

			for _, a := range response.GetAnnotation() {
				fmt.Printf("Param: %s = '%s'\n", a.Key, a.Value)
			}
		} else {
			fmt.Println("Too many arguments")
			cmd.Usage()
		}
	},
}

func init() {
	RootCmd.AddCommand(scriptParametersCmd)
}
