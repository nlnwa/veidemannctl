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
	frontierV1 "github.com/nlnwa/veidemann-api/go/frontier/v1"
	reportV1 "github.com/nlnwa/veidemann-api/go/report/v1"
	"github.com/nlnwa/veidemannctl/src/apiutil"
	"github.com/nlnwa/veidemannctl/src/connection"
	"github.com/nlnwa/veidemannctl/src/format"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
	"os"
)

type jobExecConf struct {
	filter     string
	pageSize   int32
	page       int32
	goTemplate string
	format     string
	file       string
	watch      bool
}

var jobExecFlags = &jobExecConf{}

// jobexecutionCmd represents the jobexecution command
var jobexecutionCmd = &cobra.Command{
	Use:   "jobexecution",
	Short: "Get current status for job executions",
	Long:  `Get current status for job executions.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, conn := connection.NewReportClient()
		defer conn.Close()

		var ids []string

		if len(args) > 0 {
			ids = args
		}

		request, err := CreateJobExecutionsListRequest(ids)
		if err != nil {
			log.Fatalf("Error creating request: %v", err)
		}

		r, err := client.ListJobExecutions(context.Background(), request)
		if err != nil {
			log.Fatalf("Error from controller: %v", err)
		}

		out, err := format.ResolveWriter(jobExecFlags.file)
		if err != nil {
			log.Fatalf("Could not resolve output '%v': %v", jobExecFlags.file, err)
		}
		s, err := format.NewFormatter("JobExecutionStatus", out, jobExecFlags.format, jobExecFlags.goTemplate)
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
	jobexecutionCmd.Flags().Int32VarP(&jobExecFlags.pageSize, "pagesize", "s", 10, "Number of objects to get")
	jobexecutionCmd.Flags().Int32VarP(&jobExecFlags.page, "page", "p", 0, "The page number")
	jobexecutionCmd.Flags().StringVarP(&jobExecFlags.format, "output", "o", "table", "Output format (table|wide|json|yaml|template|template-file)")
	jobexecutionCmd.Flags().StringVarP(&jobExecFlags.goTemplate, "template", "t", "", "A Go template used to format the output")
	jobexecutionCmd.Flags().StringVarP(&jobExecFlags.filter, "filter", "q", "", "Filter objects by field (i.e. meta.description=foo")
	jobexecutionCmd.Flags().StringVarP(&jobExecFlags.file, "filename", "f", "", "File name to write to")
	jobexecutionCmd.Flags().BoolVarP(&jobExecFlags.watch, "watch", "w", false, "Get a continous stream of changes")

	ReportCmd.AddCommand(jobexecutionCmd)
}

func CreateJobExecutionsListRequest(ids []string) (*reportV1.JobExecutionsListRequest, error) {
	request := &reportV1.JobExecutionsListRequest{}
	request.Id = ids
	request.Watch = jobExecFlags.watch
	if jobExecFlags.watch {
		jobExecFlags.pageSize = 0
	}

	request.Offset = jobExecFlags.page
	request.PageSize = jobExecFlags.pageSize

	if jobExecFlags.filter != "" {
		m, o, err := apiutil.CreateTemplateFilter(jobExecFlags.filter, &frontierV1.JobExecutionStatus{})
		if err != nil {
			return nil, err
		}
		request.QueryMask = m
		request.QueryTemplate = o.(*frontierV1.JobExecutionStatus)
	}

	return request, nil
}
