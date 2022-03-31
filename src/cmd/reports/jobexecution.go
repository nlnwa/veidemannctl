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
	"github.com/spf13/viper"
	"io"
	"time"
)

type jobExecFlags struct {
	filters     []string
	states      []string
	pageSize    int32
	page        int32
	orderByPath string
	orderDesc   bool
	to          *time.Time
	from        *time.Time
	goTemplate  string
	format      string
	file        string
	watch       bool
}

var jobExecConf = jobExecFlags{}

func newJobExecutionCmd() *cobra.Command {
	// jobexecutionCmd represents the jobexecution command
	var cmd = &cobra.Command{
		Use:   "jobexecution",
		Short: "Get current status for job executions",
		Long:  `Get current status for job executions.`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			v := viper.New()

			if err := v.BindPFlag("to", cmd.Flag("to")); err != nil {
				return fmt.Errorf("failed to bind flag: %w", err)
			}
			if v.IsSet("to") {
				to := v.GetTime("to")
				crawlExecFlags.to = &to
			}
			if err := v.BindPFlag("from", cmd.Flag("from")); err != nil {
				return fmt.Errorf("failed to bind flag: %w", err)
			}
			if v.IsSet("from") {
				from := v.GetTime("from")
				crawlExecFlags.from = &from
			}
			return nil
		},
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

			out, err := format.ResolveWriter(jobExecConf.file)
			if err != nil {
				return fmt.Errorf("could not resolve output '%v': %v", jobExecConf.file, err)
			}
			s, err := format.NewFormatter("JobExecutionStatus", out, jobExecConf.format, jobExecConf.goTemplate)
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

	cmd.Flags().Int32VarP(&jobExecConf.pageSize, "pagesize", "s", 10, "Number of objects to get")
	cmd.Flags().Int32VarP(&jobExecConf.page, "page", "p", 0, "The page number")
	cmd.Flags().StringVarP(&jobExecConf.format, "output", "o", "table", "Output format (table|wide|json|yaml|template|template-file)")
	cmd.Flags().StringVarP(&jobExecConf.goTemplate, "template", "t", "", "A Go template used to format the output")
	cmd.Flags().StringSliceVarP(&jobExecConf.filters, "filter", "q", nil, "Filter objects by field (i.e. meta.description=foo")
	cmd.Flags().StringSliceVar(&jobExecConf.states, "state", nil, "Filter objects by state(s)")
	cmd.Flags().StringVarP(&jobExecConf.file, "filename", "f", "", "File name to write to")
	cmd.Flags().StringVar(&jobExecConf.orderByPath, "order-by", "", "Order by path")
	cmd.Flags().BoolVar(&jobExecConf.orderDesc, "desc", false, "Order descending")
	cmd.Flags().String("to", "", "To start time")
	cmd.Flags().String("from", "", "From start time")
	cmd.Flags().BoolVarP(&jobExecConf.watch, "watch", "w", false, "Get a continous stream of changes")

	return cmd
}

func createJobExecutionsListRequest(ids []string) (*reportV1.JobExecutionsListRequest, error) {
	request := &reportV1.JobExecutionsListRequest{
		Id:              ids,
		Watch:           jobExecConf.watch,
		PageSize:        jobExecConf.pageSize,
		Offset:          jobExecConf.page,
		OrderByPath:     jobExecConf.orderByPath,
		OrderDescending: jobExecConf.orderDesc,
	}
	if jobExecConf.watch {
		request.PageSize = 0
	}

	if len(jobExecConf.states) > 0 {
		for _, state := range jobExecConf.states {
			if s, ok := frontierV1.JobExecutionStatus_State_value[state]; !ok {
				return nil, fmt.Errorf("not a jobexecution state: %s", state)
			} else {
				request.State = append(request.State, frontierV1.JobExecutionStatus_State(s))
			}
		}
	}

	if len(jobExecConf.filters) > 0 {
		queryMask := new(commonsV1.FieldMask)
		queryTemplate := new(frontierV1.JobExecutionStatus)

		for _, filter := range jobExecConf.filters {
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
