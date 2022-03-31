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
	"fmt"
	commonsV1 "github.com/nlnwa/veidemann-api/go/commons/v1"
	frontierV1 "github.com/nlnwa/veidemann-api/go/frontier/v1"
	reportV1 "github.com/nlnwa/veidemann-api/go/report/v1"
	"github.com/nlnwa/veidemannctl/src/apiutil"
	"github.com/nlnwa/veidemannctl/src/connection"
	"github.com/nlnwa/veidemannctl/src/format"
	"github.com/spf13/cobra"
	"io"
)

var crawlExecFlags struct {
	filters     []string
	pageSize    int32
	page        int32
	goTemplate  string
	format      string
	file        string
	orderByPath string
	orderDesc   bool
	watch       bool
	states      []string
}

// crawlexecutionCmd represents the crawlexecution command
var crawlexecutionCmd = &cobra.Command{
	Use:   "crawlexecution",
	Short: "Get current status for crawl executions",
	Long:  `Get current status for crawl executions.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, conn := connection.NewReportClient()
		defer conn.Close()

		var ids []string

		if len(args) > 0 {
			ids = args
		}

		request, err := CreateCrawlExecutionsListRequest(ids)
		if err != nil {
			return fmt.Errorf("failed creating request: %w", err)
		}

		cmd.SilenceUsage = true

		r, err := client.ListExecutions(context.Background(), request)
		if err != nil {
			return fmt.Errorf("error from controller: %w", err)
		}
		out, err := format.ResolveWriter(crawlExecFlags.file)
		if err != nil {
			return fmt.Errorf("could not resolve output '%s': %w", crawlExecFlags.file, err)
		}
		s, err := format.NewFormatter("CrawlExecutionStatus", out, crawlExecFlags.format, crawlExecFlags.goTemplate)
		if err != nil {
			return err
		}
		defer s.Close()

		for {
			msg, err := r.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			if err := s.WriteRecord(msg); err != nil {
				return err
			}
		}
		return nil
	},
}

func init() {
	crawlexecutionCmd.Flags().Int32VarP(&crawlExecFlags.pageSize, "pagesize", "s", 10, "Number of objects to get")
	crawlexecutionCmd.Flags().Int32VarP(&crawlExecFlags.page, "page", "p", 0, "The page number")
	crawlexecutionCmd.Flags().StringVarP(&crawlExecFlags.format, "output", "o", "table", "Output format (table|wide|json|yaml|template|template-file)")
	crawlexecutionCmd.Flags().StringVarP(&crawlExecFlags.goTemplate, "template", "t", "", "A Go template used to format the output")
	crawlexecutionCmd.Flags().StringSliceVarP(&crawlExecFlags.filters, "filter", "q", nil, "Filter objects by field (i.e. meta.description=foo")
	crawlexecutionCmd.Flags().StringSliceVar(&crawlExecFlags.states, "state", nil, "Filter objects by state(s)")
	crawlexecutionCmd.Flags().StringVarP(&crawlExecFlags.file, "filename", "f", "", "File name to write to")
	crawlexecutionCmd.Flags().StringVar(&crawlExecFlags.orderByPath, "order-by", "", "Order by path")
	crawlexecutionCmd.Flags().BoolVar(&crawlExecFlags.orderDesc, "desc", false, "Order descending")
	crawlexecutionCmd.Flags().BoolVarP(&crawlExecFlags.watch, "watch", "w", false, "Get a continous stream of changes")

	ReportCmd.AddCommand(crawlexecutionCmd)
}

func CreateCrawlExecutionsListRequest(ids []string) (*reportV1.CrawlExecutionsListRequest, error) {
	request := &reportV1.CrawlExecutionsListRequest{
		Id:       ids,
		Watch:    crawlExecFlags.watch,
		PageSize: crawlExecFlags.pageSize,
		Offset:   crawlExecFlags.page,
		OrderByPath: crawlExecFlags.orderByPath,
		OrderDescending: crawlExecFlags.orderDesc,
	}

	if crawlExecFlags.watch {
		request.PageSize = 0
	}

	if len(crawlExecFlags.states) > 0 {
		for _, state := range crawlExecFlags.states {
			if s, ok := frontierV1.CrawlExecutionStatus_State_value[state]; !ok {
				return nil, fmt.Errorf("not a crawlexecution state: %s", state)
			} else {
				request.State = append(request.State, frontierV1.CrawlExecutionStatus_State(s))
			}
		}
	}

	if len(crawlExecFlags.filters) > 0 {
		queryTemplate := new(frontierV1.CrawlExecutionStatus)
		queryMask := new(commonsV1.FieldMask)

		for _, filter := range crawlExecFlags.filters {
			err := apiutil.CreateTemplateFilter(filter, queryTemplate, queryMask)
			if err != nil {
				return nil, err
			}
		}

		request.QueryMask = queryMask
		request.QueryTemplate = queryTemplate
	}

	return request, nil
}
