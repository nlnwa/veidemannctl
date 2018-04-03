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

package reports

import (
	api "github.com/nlnwa/veidemannctl/veidemann_api"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	tspb "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/nlnwa/veidemannctl/bindata"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

var (
	executionId    string
	pageSize       int32
	page           int32
	goTemplate     string
	filter         []string
	//validArgs      = []string{"crawllog", "pagelog", "screenshot", "crawlexecution", "jobexecution"}
	//validArgs      = []string{"crawllog", "pagelog", "screenshot"}
)

// reportCmd represents the report command
var ReportCmd = &cobra.Command{
	Use:   "report",
	Short: "Get log report",
	Long: `Request a report.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	ReportCmd.PersistentFlags().StringVarP(&executionId, "executionid", "e", "", "All objects by Execution ID")
	ReportCmd.PersistentFlags().Int32VarP(&pageSize, "pagesize", "s", 10, "Number of objects to get")
	ReportCmd.PersistentFlags().Int32VarP(&page, "page", "p", 0, "The page number")
	ReportCmd.PersistentFlags().StringVarP(&goTemplate, "template", "t", "", "A Go template used to format the output")
	ReportCmd.PersistentFlags().StringSliceVarP(&filter, "filter", "f", nil, "Filters")
}

func applyFilter(filter []string) []*api.Filter {
	var result []*api.Filter
	for _, f := range filter {
		tokens := strings.SplitN(f, " ", 3)
		op := api.Filter_Operator(api.Filter_Operator_value[strings.ToUpper(tokens[1])])
		result = append(result, &api.Filter{tokens[0], op, tokens[2]})
	}
	return result
}

func ApplyTemplate(msg proto.Message, defaultTemplate string) {
	var data []byte
	var err error
	if goTemplate == "" {
		data, err = bindata.Asset(defaultTemplate)
		if err != nil {
			panic(err)
		}
	} else {
		data, err = ioutil.ReadFile(goTemplate)
		if err != nil {
			panic(err)
		}
	}

	ESC := string(0x1b)
	funcMap := template.FuncMap{
		"reset":         func() string { return ESC + "[0m" },
		"bold":          func() string { return ESC + "[1m" },
		"inverse":       func() string { return ESC + "[7m" },
		"red":           func() string { return ESC + "[31m" },
		"green":         func() string { return ESC + "[32m" },
		"yellow":        func() string { return ESC + "[33m" },
		"blue":          func() string { return ESC + "[34m" },
		"magenta":       func() string { return ESC + "[35m" },
		"cyan":          func() string { return ESC + "[36m" },
		"brightred":     func() string { return ESC + "[1;31m" },
		"brightgreen":   func() string { return ESC + "[1;32m" },
		"brightyellow":  func() string { return ESC + "[1;33m" },
		"brightblue":    func() string { return ESC + "[1;34m" },
		"brightmagenta": func() string { return ESC + "[1;35m" },
		"brightcyan":    func() string { return ESC + "[1;36m" },
		"bgwhite":       func() string { return ESC + "[47m" },
		"bgbrightblack": func() string { return ESC + "[100m" },
		"time":          func(ts *tspb.Timestamp) string { return ptypes.TimestampString(ts) },
	}

	tmpl, err := template.New(defaultTemplate).Funcs(funcMap).Parse(string(data))
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(os.Stdout, msg)
	if err != nil {
		panic(err)
	}
}
