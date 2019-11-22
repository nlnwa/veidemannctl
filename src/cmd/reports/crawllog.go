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

var crawllogFlags struct {
	executionId string
	filter      string
	pageSize    int32
	page        int32
	goTemplate  string
	format      string
	file        string
	watch       bool
}

// crawllogCmd represents the crawllog command
var crawllogCmd = &cobra.Command{
	Use:   "crawllog",
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

		request, err := CreateCrawlLogListRequest(ids)
		if err != nil {
			log.Fatalf("Error creating request: %v", err)
		}

		r, err := client.ListCrawlLogs(context.Background(), request)
		if err != nil {
			log.Fatalf("Error from controller: %v", err)
		}

		out, err := format.ResolveWriter(crawllogFlags.file)
		if err != nil {
			log.Fatalf("Could not resolve output '%v': %v", crawlExecFlags.file, err)
		}
		s, err := format.NewFormatter("CrawlLog", out, crawllogFlags.format, crawllogFlags.goTemplate)
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
			log.Debugf("Outputting crawl log record with WARC id '%s'", msg.WarcId)
			if s.WriteRecord(msg) != nil {
				os.Exit(1)
			}
		}
	},
}

func init() {
	crawllogCmd.Flags().Int32VarP(&crawllogFlags.pageSize, "pagesize", "s", 10, "Number of objects to get")
	crawllogCmd.Flags().Int32VarP(&crawllogFlags.page, "page", "p", 0, "The page number")
	crawllogCmd.Flags().StringVarP(&crawllogFlags.format, "output", "o", "json", "Output format (json|yaml|template|template-file)")
	crawllogCmd.Flags().StringVarP(&crawllogFlags.goTemplate, "template", "t", "", "A Go template used to format the output")
	crawllogCmd.Flags().StringVarP(&crawllogFlags.filter, "filter", "q", "", "Filter objects by field (i.e. meta.description=foo")
	crawllogCmd.Flags().StringVarP(&crawllogFlags.file, "filename", "f", "", "File name to write to")
	crawllogCmd.Flags().BoolVarP(&crawllogFlags.watch, "watch", "w", false, "Get a continous stream of changes")

	ReportCmd.AddCommand(crawllogCmd)
}

func CreateCrawlLogListRequest(ids []string) (*reportV1.CrawlLogListRequest, error) {
	request := &reportV1.CrawlLogListRequest{}
	request.WarcId = ids
	request.Watch = crawllogFlags.watch
	if crawllogFlags.watch {
		crawllogFlags.pageSize = 0
	}

	request.Offset = crawllogFlags.page
	request.PageSize = crawllogFlags.pageSize

	if crawllogFlags.filter != "" {
		m, o, err := apiutil.CreateTemplateFilter(crawllogFlags.filter, &frontierV1.CrawlLog{})
		if err != nil {
			return nil, err
		}
		request.QueryMask = m
		request.QueryTemplate = o.(*frontierV1.CrawlLog)
	}

	return request, nil
}
