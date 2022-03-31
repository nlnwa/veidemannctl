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

type jobExecConf struct {
	filters     []string
	states      []string
	pageSize    int32
	page        int32
	orderByPath string
	orderDesc   bool
	goTemplate  string
	format      string
	file        string
	watch       bool
}

var jobExecFlags = &jobExecConf{}

// jobexecutionCmd represents the jobexecution command
var jobexecutionCmd = &cobra.Command{
	Use:   "jobexecution",
	Short: "Get current status for job executions",
	Long:  `Get current status for job executions.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, conn := connection.NewReportClient()
		defer conn.Close()

		var ids []string

		if len(args) > 0 {
			ids = args
		}

		request, err := createJobExecutionsListRequest(ids)
		if err != nil {
			return fmt.Errorf("error creating request: %w", err)
		}

		cmd.SilenceUsage = true

		r, err := client.ListJobExecutions(context.Background(), request)
		if err != nil {
			return fmt.Errorf("error from controller: %v", err)
		}

		out, err := format.ResolveWriter(jobExecFlags.file)
		if err != nil {
			return fmt.Errorf("could not resolve output '%v': %v", jobExecFlags.file, err)
		}
		s, err := format.NewFormatter("JobExecutionStatus", out, jobExecFlags.format, jobExecFlags.goTemplate)
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
	jobexecutionCmd.Flags().Int32VarP(&jobExecFlags.pageSize, "pagesize", "s", 10, "Number of objects to get")
	jobexecutionCmd.Flags().Int32VarP(&jobExecFlags.page, "page", "p", 0, "The page number")
	jobexecutionCmd.Flags().StringVarP(&jobExecFlags.format, "output", "o", "table", "Output format (table|wide|json|yaml|template|template-file)")
	jobexecutionCmd.Flags().StringVarP(&jobExecFlags.goTemplate, "template", "t", "", "A Go template used to format the output")
	jobexecutionCmd.Flags().StringSliceVarP(&jobExecFlags.filters, "filter", "q", nil, "Filter objects by field (i.e. meta.description=foo")
	jobexecutionCmd.Flags().StringSliceVar(&jobExecFlags.states, "state", nil, "Filter objects by state(s)")
	jobexecutionCmd.Flags().StringVarP(&jobExecFlags.file, "filename", "f", "", "File name to write to")
	jobexecutionCmd.Flags().StringVar(&jobExecFlags.orderByPath, "order-by", "", "Order by path")
	jobexecutionCmd.Flags().BoolVar(&jobExecFlags.orderDesc, "desc", false, "Order descending")
	jobexecutionCmd.Flags().BoolVarP(&jobExecFlags.watch, "watch", "w", false, "Get a continous stream of changes")

	ReportCmd.AddCommand(jobexecutionCmd)
}

func createJobExecutionsListRequest(ids []string) (*reportV1.JobExecutionsListRequest, error) {
	request := &reportV1.JobExecutionsListRequest{
		Id:              ids,
		Watch:           jobExecFlags.watch,
		PageSize:        jobExecFlags.pageSize,
		Offset:          jobExecFlags.page,
		OrderByPath:     jobExecFlags.orderByPath,
		OrderDescending: jobExecFlags.orderDesc,
	}
	if jobExecFlags.watch {
		request.PageSize = 0
	}

	if len(jobExecFlags.states) > 0 {
		for _, state := range jobExecFlags.states {
			if s, ok := frontierV1.JobExecutionStatus_State_value[state]; !ok {
				return nil, fmt.Errorf("not a jobexecution state: %s", state)
			} else {
				request.State = append(request.State, frontierV1.JobExecutionStatus_State(s))
			}
		}
	}

	if len(jobExecFlags.filters) > 0 {
		queryMask := new(commonsV1.FieldMask)
		queryTemplate := new(frontierV1.JobExecutionStatus)

		for _, filter := range jobExecFlags.filters {
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
