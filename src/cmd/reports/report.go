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
	"fmt"
	api "github.com/nlnwa/veidemann-api-go/veidemann_api"
	"io/ioutil"
	"os"

	"encoding/json"
	"github.com/golang/protobuf/ptypes"
	tspb "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/nlnwa/veidemannctl/bindata"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"strings"
	"text/template"
)

var flags struct {
	executionId string
	pageSize    int32
	page        int32
	goTemplate  string
	filter      []string
	format      string
	quiet       bool
}

// reportCmd represents the report command
var ReportCmd = &cobra.Command{
	Use:   "report",
	Short: "Get log report",
	Long:  `Request a report.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	ReportCmd.PersistentFlags().StringVarP(&flags.executionId, "executionid", "e", "", "All objects by Execution ID")
	ReportCmd.PersistentFlags().Int32VarP(&flags.pageSize, "pagesize", "s", 10, "Number of objects to get")
	ReportCmd.PersistentFlags().Int32VarP(&flags.page, "page", "p", 0, "The page number")
	ReportCmd.PersistentFlags().StringVarP(&flags.format, "output", "o", "json", "Output format (json|yaml|template|template-file)")
	ReportCmd.PersistentFlags().StringVarP(&flags.goTemplate, "template", "t", "", "A Go template used to format the output")
	ReportCmd.PersistentFlags().StringSliceVarP(&flags.filter, "filter", "f", nil, "Filters")
	ReportCmd.PersistentFlags().BoolVarP(&flags.quiet, "quiet", "q", false, "Quiet")
	ReportCmd.PersistentFlags().MarkHidden("quiet")
}

func applyFilter(filter []string) []*api.Filter {
	var result []*api.Filter
	for _, f := range filter {
		tokens := strings.SplitN(f, " ", 3)
		op := api.Filter_Operator(api.Filter_Operator_value[strings.ToUpper(tokens[1])])
		result = append(result, &api.Filter{FieldName: tokens[0], Op: op, Value: tokens[2]})
	}
	return result
}

func ApplyTemplate(msg interface{}, defaultTemplate string) {
	var data []byte
	var err error
	if flags.goTemplate == "" {
		data, err = bindata.Asset(defaultTemplate)
		if err != nil {
			panic(err)
		}
	} else {
		data, err = ioutil.ReadFile(flags.goTemplate)
		if err != nil {
			panic(err)
		}
	}

	RunTemplate(msg, string(data))
}

func RunTemplate(msg interface{}, templateString string) {
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
		"time": func(ts *tspb.Timestamp) string {
			if ts == nil {
				return "                        "
			} else {
				return fmt.Sprintf("%-24.24s", ptypes.TimestampString(ts))
			}
		},
		"json": func(v interface{}) string {
			if v == nil {
				return ""
			} else {
				json, err := json.Marshal(v)
				if err != nil {
					log.Fatal(err)
				}
				return string(json)
			}
		},
		"prettyJson": func(v interface{}) string {
			if v == nil {
				return ""
			} else {
				json, err := json.MarshalIndent(v, "", "  ")
				if err != nil {
					log.Fatal(err)
				}
				return string(json)
			}
		},
	}

	tmpl, err := template.New("Template").Funcs(funcMap).Parse(templateString)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(os.Stdout, msg)
	if err != nil {
		panic(err)
	}
}
