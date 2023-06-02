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

package report

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	commonsV1 "github.com/nlnwa/veidemann-api/go/commons/v1"
	frontierV1 "github.com/nlnwa/veidemann-api/go/frontier/v1"
	reportV1 "github.com/nlnwa/veidemann-api/go/report/v1"
	"github.com/nlnwa/veidemannctl/apiutil"
	"github.com/nlnwa/veidemannctl/connection"
	"github.com/nlnwa/veidemannctl/format"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type jobExecutionCmdOptions struct {
	ids         []string
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

func (o *jobExecutionCmdOptions) complete(cmd *cobra.Command, args []string) error {
	o.ids = args

	v := viper.New()

	if err := v.BindPFlag("to", cmd.Flag("to")); err != nil {
		return fmt.Errorf("failed to bind flag: %w", err)
	}
	if v.IsSet("to") {
		to := v.GetTime("to")
		o.to = &to
	}
	if err := v.BindPFlag("from", cmd.Flag("from")); err != nil {
		return fmt.Errorf("failed to bind flag: %w", err)
	}
	if v.IsSet("from") {
		from := v.GetTime("from")
		o.from = &from
	}

	return nil
}

func (o *jobExecutionCmdOptions) run() error {
	conn, err := connection.Connect()
	if err != nil {
		return err
	}
	defer conn.Close()

	client := reportV1.NewReportClient(conn)

	request, err := o.createJobExecutionsListRequest()
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	r, err := client.ListJobExecutions(context.Background(), request)
	if err != nil {
		return fmt.Errorf("error from controller: %w", err)
	}

	w, err := format.ResolveWriter(o.file)
	if err != nil {
		return fmt.Errorf("unable to open output file: %v: %w", o.file, err)
	}
	s, err := format.NewFormatter("JobExecutionStatus", w, o.format, o.goTemplate)
	if err != nil {
		return fmt.Errorf("error creating formatter: %w", err)
	}
	defer s.Close()

	for {
		msg, err := r.Recv()
		if errors.Is(err, io.EOF) {
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
}

func newJobExecutionCmd() *cobra.Command {
	o := &jobExecutionCmdOptions{}
	var cmd = &cobra.Command{
		Use:     "jobexecution [ID ...]",
		Short:   "Get current status for job executions",
		Long:    `Get current status for job executions.`,
		PreRunE: o.complete,
		RunE: func(cmd *cobra.Command, args []string) error {
			// set silence usage to true to avoid printing usage when an error occurs
			cmd.SilenceUsage = true
			return o.run()
		},
	}

	cmd.Flags().Int32VarP(&o.pageSize, "pagesize", "s", 10, "Number of objects to get")
	cmd.Flags().Int32VarP(&o.page, "page", "p", 0, "The page number")
	cmd.Flags().StringVarP(&o.format, "output", "o", "table", "Output format (table|wide|json|yaml|template|template-file)")
	cmd.Flags().StringVarP(&o.goTemplate, "template", "t", "", "A Go template used to format the output")
	cmd.Flags().StringSliceVarP(&o.filters, "filter", "q", nil, "Filter objects by field (i.e. meta.description=foo")
	cmd.Flags().StringSliceVar(&o.states, "state", nil, "Filter objects by state(s)")
	cmd.Flags().StringVarP(&o.file, "filename", "f", "", "Filename to write to")
	cmd.Flags().StringVar(&o.orderByPath, "order-by", "", "Order by path")
	cmd.Flags().BoolVar(&o.orderDesc, "desc", false, "Order descending")
	cmd.Flags().String("to", "", "To start time")
	cmd.Flags().String("from", "", "From start time")
	cmd.Flags().BoolVarP(&o.watch, "watch", "w", false, "Get a continous stream of changes")

	return cmd
}

func (o *jobExecutionCmdOptions) createJobExecutionsListRequest() (*reportV1.JobExecutionsListRequest, error) {
	request := &reportV1.JobExecutionsListRequest{
		Id:              o.ids,
		Watch:           o.watch,
		PageSize:        o.pageSize,
		Offset:          o.page,
		OrderByPath:     o.orderByPath,
		OrderDescending: o.orderDesc,
	}
	if o.watch {
		request.PageSize = 0
	}

	if len(o.states) > 0 {
		for _, state := range o.states {
			if s, ok := frontierV1.JobExecutionStatus_State_value[state]; !ok {
				return nil, fmt.Errorf("not a jobexecution state: %s", state)
			} else {
				request.State = append(request.State, frontierV1.JobExecutionStatus_State(s))
			}
		}
	}

	if len(o.filters) > 0 {
		queryMask := new(commonsV1.FieldMask)
		queryTemplate := new(frontierV1.JobExecutionStatus)
		for _, filter := range o.filters {
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
