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
	frontierV1 "github.com/nlnwa/veidemann-api-go/frontier/v1"
	reportV1 "github.com/nlnwa/veidemann-api-go/report/v1"
	"github.com/nlnwa/veidemannctl/src/apiutil"
	"github.com/nlnwa/veidemannctl/src/connection"
	"github.com/nlnwa/veidemannctl/src/format"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
	"os"
)

var crawlExecFlags struct {
	filter     string
	pageSize   int32
	page       int32
	goTemplate string
	format     string
	file       string
	watch      bool
}

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
		client, conn := connection.NewReportClient()
		defer conn.Close()

		var ids []string

		if len(args) > 0 {
			ids = args
		}

		request, err := CreateCrawlExecutionsListRequest(ids)
		if err != nil {
			log.Fatalf("Error creating request: %v", err)
		}

		r, err := client.ListExecutions(context.Background(), request)
		if err != nil {
			log.Fatalf("Error from controller: %v", err)
		}

		out, err := format.ResolveWriter(crawlExecFlags.file)
		if err != nil {
			log.Fatalf("Could not resolve output '%v': %v", crawlExecFlags.file, err)
		}
		s, err := format.NewFormatter("CrawlExecutionStatus", out, crawlExecFlags.format, crawlExecFlags.goTemplate)
		if err != nil {
			log.Fatal(err)
		}
		defer s.Close()

		for {
			msg, err := r.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error getting object: %v", err)
			}
			log.Debugf("Outputting page log record with id '%s'", msg.Id)
			if s.WriteRecord(msg) != nil {
				os.Exit(1)
			}
		}
	},
}

func init() {
	crawlexecutionCmd.Flags().Int32VarP(&crawlExecFlags.pageSize, "pagesize", "s", 10, "Number of objects to get")
	crawlexecutionCmd.Flags().Int32VarP(&crawlExecFlags.page, "page", "p", 0, "The page number")
	crawlexecutionCmd.Flags().StringVarP(&crawlExecFlags.format, "output", "o", "table", "Output format (table|json|yaml|template|template-file)")
	crawlexecutionCmd.Flags().StringVarP(&crawlExecFlags.goTemplate, "template", "t", "", "A Go template used to format the output")
	crawlexecutionCmd.Flags().StringVarP(&crawlExecFlags.filter, "filter", "q", "", "Filter objects by field (i.e. meta.description=foo")
	crawlexecutionCmd.Flags().StringVarP(&crawlExecFlags.file, "filename", "f", "", "File name to write to")
	crawlexecutionCmd.Flags().BoolVarP(&crawlExecFlags.watch, "watch", "w", false, "Get a continous stream of changes")

	ReportCmd.AddCommand(crawlexecutionCmd)
}

func CreateCrawlExecutionsListRequest(ids []string) (*reportV1.CrawlExecutionsListRequest, error) {
	request := &reportV1.CrawlExecutionsListRequest{}
	request.Id = ids
	request.Watch = crawlExecFlags.watch
	if crawlExecFlags.watch {
		crawlExecFlags.pageSize = 0
	}

	request.Offset = crawlExecFlags.page
	request.PageSize = crawlExecFlags.pageSize

	if crawlExecFlags.filter != "" {
		m, o, err := apiutil.CreateTemplateFilter(crawlExecFlags.filter, &frontierV1.CrawlExecutionStatus{})
		if err != nil {
			return nil, err
		}
		request.QueryMask = m
		request.QueryTemplate = o.(*frontierV1.CrawlExecutionStatus)
	}

	return request, nil
}
