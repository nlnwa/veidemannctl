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

package reports

import (
	"context"
	"github.com/nlnwa/veidemannctl/src/connection"
	api "github.com/nlnwa/veidemannctl/veidemann_api"
	"github.com/spf13/cobra"
	"log"
)

// crawlexecutionCmd represents the crawlexecution command
var crawlexecutionCmd = &cobra.Command{
	Use:   "crawlexecution",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, conn := connection.NewStatusClient()
		defer conn.Close()

		request := api.ListExecutionsRequest{}
		if len(args) > 1 {
			request.Id = args[1:]
		}
		//if executionId != "" {
		//	request.ExecutionId = executionId
		//}
		//request.Filter = applyFilter(filter)
		request.Page = page
		request.PageSize = pageSize

		r, err := client.ListExecutions(context.Background(), &request)
		if err != nil {
			log.Fatalf("could not get crawl log: %v", err)
		}

		ApplyTemplate(r, "crawlexecution.template")
	},
}

func init() {
	ReportCmd.AddCommand(crawlexecutionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// crawlexecutionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// crawlexecutionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
