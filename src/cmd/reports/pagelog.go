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

// pagelogCmd represents the pagelog command
var pagelogCmd = &cobra.Command{
	Use:   "pagelog",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, conn := connection.NewReportClient()
		defer conn.Close()

		request := api.PageLogListRequest{}
		if len(args) > 0 {
			request.WarcId = args
		}
		if flags.executionId != "" {
			request.ExecutionId = flags.executionId
		}
		request.Filter = applyFilter(flags.filter)
		request.Page = flags.page
		request.PageSize = flags.pageSize

		r, err := client.ListPageLogs(context.Background(), &request)
		if err != nil {
			log.Fatalf("could not get page log: %v", err)
		}

		ApplyTemplate(r, "pagelog.template")
	},
}

func init() {
	ReportCmd.AddCommand(pagelogCmd)
}
