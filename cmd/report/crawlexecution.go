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
	"google.golang.org/protobuf/types/known/timestamppb"
)

// crawlExecCmdOptions holds the crawl execution command line options
type crawlExecCmdOptions struct {
	ids         []string
	filters     []string
	pageSize    int32
	page        int32
	goTemplate  string
	format      string
	file        string
	orderByPath string
	orderDesc   bool
	to          *time.Time
	from        *time.Time
	watch       bool
	states      []string
}

// complete completes the crawl execution command options
func (o *crawlExecCmdOptions) complete(cmd *cobra.Command, args []string) error {
	o.ids = args

	// use viper to parse time for the flags "to" and "from"
	v := viper.New()

	if err := v.BindPFlag("to", cmd.Flag("to")); err != nil {
		return fmt.Errorf("failed to bind flag: %w", err)
	} else if v.IsSet("to") {
		to := v.GetTime("to")
		o.to = &to
	}
	if err := v.BindPFlag("from", cmd.Flag("from")); err != nil {
		return fmt.Errorf("failed to bind flag: %w", err)
	} else if v.IsSet("from") {
		from := v.GetTime("from")
		o.from = &from
	}
	return nil
}

// run runs the crawl execution command
func (o *crawlExecCmdOptions) run() error {
	conn, err := connection.Connect()
	if err != nil {
		return err
	}
	defer conn.Close()

	client := reportV1.NewReportClient(conn)

	request, err := o.createCrawlExecutionsListRequest()
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	r, err := client.ListExecutions(context.Background(), request)
	if err != nil {
		return fmt.Errorf("failed to list crawl executions: %w", err)
	}
	out, err := format.ResolveWriter(o.file)
	if err != nil {
		return fmt.Errorf("unable to open output file: %v: %w", o.file, err)
	}
	s, err := format.NewFormatter("CrawlExecutionStatus", out, o.format, o.goTemplate)
	if err != nil {
		return err
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

// newCrawlExecutionCmd creates the crawl execution command
func newCrawlExecutionCmd() *cobra.Command {
	o := &crawlExecCmdOptions{}

	cmd := &cobra.Command{
		Use:   "crawlexecution [ID ...]",
		Short: "Get current status for crawl executions",
		Long:  `Get current status for crawl executions.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.complete(cmd, args); err != nil {
				return err
			}
			cmd.SilenceUsage = true

			return o.run()
		},
	}
	cmd.Flags().Int32VarP(&o.pageSize, "pagesize", "s", 10, "Number of objects to get")
	cmd.Flags().Int32VarP(&o.page, "page", "p", 0, "The page number")
	cmd.Flags().StringVarP(&o.format, "output", "o", "table", "Output format (table|wide|json|yaml|template|template-file)")
	cmd.Flags().StringVarP(&o.goTemplate, "template", "t", "", "A Go template used to format the output")
	cmd.Flags().StringSliceVarP(&o.filters, "filter", "q", nil, "Filter objects by field (i.e. meta.description=foo)")
	cmd.Flags().StringSliceVar(&o.states, "state", nil, "Filter objects by state. Valid states are UNDEFINED, FETCHING, SLEEPING, FINISHED or FAILED")
	cmd.Flags().StringVarP(&o.file, "filename", "f", "", "Filename to write to")
	cmd.Flags().StringVar(&o.orderByPath, "order-by", "", "Order by path")
	cmd.Flags().String("to", "", "To start time")
	cmd.Flags().String("from", "", "From start time")
	cmd.Flags().BoolVar(&o.orderDesc, "desc", false, "Order descending")
	cmd.Flags().BoolVarP(&o.watch, "watch", "w", false, "Get a continous stream of changes")

	return cmd
}

// createCrawlExecutionsListRequest creates a crawl execution list request
func (o *crawlExecCmdOptions) createCrawlExecutionsListRequest() (*reportV1.CrawlExecutionsListRequest, error) {
	request := &reportV1.CrawlExecutionsListRequest{
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

	if o.from != nil {
		fmt.Println(o.from)
		request.StartTimeFrom = timestamppb.New(*o.from)
	}

	if o.to != nil {
		request.StartTimeTo = timestamppb.New(*o.to)
	}

	if len(o.states) > 0 {
		for _, state := range o.states {
			if s, ok := frontierV1.CrawlExecutionStatus_State_value[state]; !ok {
				return nil, fmt.Errorf("not a crawlexecution state: %s", state)
			} else {
				request.State = append(request.State, frontierV1.CrawlExecutionStatus_State(s))
			}
		}
	}

	if len(o.filters) > 0 {
		queryTemplate := new(frontierV1.CrawlExecutionStatus)
		queryMask := new(commonsV1.FieldMask)

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
