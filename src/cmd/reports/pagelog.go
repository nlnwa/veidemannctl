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
	logV1 "github.com/nlnwa/veidemann-api/go/log/v1"
	"github.com/nlnwa/veidemannctl/src/apiutil"
	"github.com/nlnwa/veidemannctl/src/connection"
	"github.com/nlnwa/veidemannctl/src/format"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
	"os"
)

var pagelogFlags struct {
	filter     string
	pageSize   int32
	page       int32
	goTemplate string
	format     string
	file       string
	watch      bool
}

// pagelogCmd represents the pagelog command
var pagelogCmd = &cobra.Command{
	Use:   "pagelog",
	Short: "View page log",
	Long:  `View page log.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, conn := connection.NewLogClient()
		defer conn.Close()

		var ids []string

		if len(args) > 0 {
			ids = args
		}

		request, err := CreatePageLogListRequest(ids)
		if err != nil {
			log.Fatalf("Error creating request: %v", err)
		}

		r, err := client.ListPageLogs(context.Background(), request)
		if err != nil {
			log.Fatalf("Error from controller: %v", err)
		}

		out, err := format.ResolveWriter(pagelogFlags.file)
		if err != nil {
			log.Fatalf("Could not resolve output '%v': %v", pagelogFlags.file, err)
		}
		s, err := format.NewFormatter("PageLog", out, pagelogFlags.format, pagelogFlags.goTemplate)
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
			log.Debugf("Outputting page log record with WARC id '%s'", msg.WarcId)
			if s.WriteRecord(msg) != nil {
				os.Exit(1)
			}
		}
	},
}

func init() {
	pagelogCmd.Flags().Int32VarP(&pagelogFlags.pageSize, "pagesize", "s", 10, "Number of objects to get")
	pagelogCmd.Flags().Int32VarP(&pagelogFlags.page, "page", "p", 0, "The page number")
	pagelogCmd.Flags().StringVarP(&pagelogFlags.format, "output", "o", "table", "Output format (table|wide|json|yaml|template|template-file)")
	pagelogCmd.Flags().StringVarP(&pagelogFlags.goTemplate, "template", "t", "", "A Go template used to format the output")
	pagelogCmd.Flags().StringVarP(&pagelogFlags.filter, "filter", "q", "", "Filter objects by field (i.e. meta.description=foo")
	pagelogCmd.Flags().StringVarP(&pagelogFlags.file, "filename", "f", "", "File name to write to")
	pagelogCmd.Flags().BoolVarP(&pagelogFlags.watch, "watch", "w", false, "Get a continous stream of changes")

	ReportCmd.AddCommand(pagelogCmd)
}

func CreatePageLogListRequest(ids []string) (*logV1.PageLogListRequest, error) {
	request := &logV1.PageLogListRequest{}
	request.WarcId = ids
	request.Watch = pagelogFlags.watch
	if pagelogFlags.watch {
		pagelogFlags.pageSize = 0
	}

	request.Offset = pagelogFlags.page
	request.PageSize = pagelogFlags.pageSize

	if pagelogFlags.filter != "" {
		m, o, err := apiutil.CreateTemplateFilter(pagelogFlags.filter, &logV1.PageLog{})
		if err != nil {
			return nil, err
		}
		request.QueryMask = m
		request.QueryTemplate = o.(*logV1.PageLog)
	}

	return request, nil
}
