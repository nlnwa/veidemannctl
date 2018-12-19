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
	api "github.com/nlnwa/veidemann-api-go/veidemann_api"
	"github.com/spf13/cobra"
	"log"
)

var (
	jobExecutionId string
	jobId string
	seedId string
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
		if len(args) > 0 {
			request.Id = args
		}
		if jobExecutionId != "" {
			request.JobExecutionId = jobExecutionId
		}
		if jobId != "" {
			request.JobId = jobId
		}
		if seedId != "" {
			request.SeedId = seedId
		}

		request.Page = flags.page
		request.PageSize = flags.pageSize

		r, err := client.ListExecutions(context.Background(), &request)
		if err != nil {
			log.Fatalf("could not get crawl log: %v", err)
		}

		ApplyTemplate(r, "crawlexecution.template")
	},
}

func init() {
	ReportCmd.AddCommand(crawlexecutionCmd)

	ReportCmd.PersistentFlags().StringVarP(&jobExecutionId, "jobexecution", "j", "", "All executions for a Job Execution ID")
	ReportCmd.PersistentFlags().StringVarP(&jobId, "job", "", "", "All executions for a Job ID")
	ReportCmd.PersistentFlags().StringVarP(&seedId, "seed", "", "", "All executions for a Seed ID")
}
