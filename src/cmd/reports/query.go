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
	"github.com/nlnwa/veidemann-api-go/config/v1"
	"github.com/nlnwa/veidemannctl/src/connection"
	api "github.com/nlnwa/veidemann-api-go/veidemann_api"
	"github.com/spf13/cobra"
	log "github.com/sirupsen/logrus"
	"fmt"
	"io"
	"strings"
	"io/ioutil"
	"github.com/ghodss/yaml"
	//"github.com/nlnwa/veidemannctl/bindata"
	"path/filepath"
	"os"
	"github.com/nlnwa/veidemannctl/src/configutil"
	"github.com/nlnwa/veidemannctl/src/format"
)

// queryCmd represents the query command
var queryCmd = &cobra.Command{
	Use:   "query [queryString|file] args...",
	Short: "Run a databse query",
	Long:  `Run a databse query. The query should be a java script string like the ones used by RethinkDb javascript driver.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			client, conn := connection.NewReportClient()
			defer conn.Close()

			request := api.ExecuteDbQueryRequest{}

			queryDef := getQueryDef(args[0])
			defer queryDef.marshalSpec.Close()

			var params []interface{}

			for _, v := range args[1:] {
				params = append(params, v)
			}

			request.Query = fmt.Sprintf(queryDef.Query, params...)

			request.Limit = flags.pageSize
			log.Debugf("Executing query: %s", request.GetQuery())

			stream, err := client.ExecuteDbQuery(context.Background(), &request)
			if err != nil {
				log.Fatalf("Failed executing query: %v", err)
			}

			if queryDef.Header != "" {
				queryDef.marshalSpec.WriteHeader()
			}

			for {
				value, err := stream.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					log.Fatalf("Query error: %v", err)
				}
				queryDef.marshalSpec.WriteRecord(value.GetRecord())
			}
		} else {
			d := configutil.GetConfigDir("query")
			if flags.quiet {
				q := listStoredQueries(d)
				for _, s := range q {
					fmt.Println(s.Name)
				}
			} else {
				fmt.Println("Missing query.\nSee 'veidemannctl report query -h' for help")
				q := listStoredQueries(d)
				if len(q) > 0 {
					fmt.Printf("\nStored queries in '%s':\n", d)
					for _, s := range q {
						fmt.Printf(" * %-20s - %s", s.Name, s.Description)
					}
				}
			}
		}
	},
}

func init() {
	ReportCmd.AddCommand(queryCmd)
}

type queryDef struct {
	Name        string
	Description string
	Query       string
	Header      string
	Template    string
	marshalSpec format.Formatter
}

func getQueryDef(queryArg string) queryDef {
	var queryDef queryDef

	if strings.HasPrefix(queryArg, "r.") {
		queryDef.Query = queryArg
	} else {
		filename := findFile(queryArg)
		log.Debugf("Using query definition from file '%s'", filename)
		readFile(filename, &queryDef)
	}

	if queryDef.Template == "" {
		queryDef.marshalSpec = format.NewFormatter(config.Kind_undefined, "", flags.format, flags.goTemplate, queryDef.Header)
	} else {
		queryDef.marshalSpec = format.NewFormatter(config.Kind_undefined, "", "template", queryDef.Template, queryDef.Header)
	}
	return queryDef
}

func findFile(name string) string {
	filename := name
	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		return filename
	}

	queryDir := configutil.GetConfigDir("query")

	filename = filepath.Join(queryDir, name)
	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		return filename
	}
	filename = filepath.Join(queryDir, name) + ".yml"
	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		return filename
	}
	filename = filepath.Join(queryDir, name) + ".yaml"
	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		return filename
	}
	log.Fatalf("Query not found: %s", name)
	return ""
}

func readFile(name string, queryDef *queryDef) {
	data, err := ioutil.ReadFile(name)
	if err != nil {
		log.Fatalf("Query not found: %v", err)
	}
	// Found file
	if strings.HasSuffix(name, ".yml") || strings.HasSuffix(name, ".yaml") {
		yaml.Unmarshal(data, &queryDef)
	} else {
		queryDef.Query = string(data)
	}
	queryDef.Name = strings.TrimSuffix(filepath.Base(name), filepath.Ext(name))
}

func listStoredQueries(path string) []queryDef {
	var r []queryDef

	if files, err := ioutil.ReadDir(path); err == nil {
		for _, f := range files {
			if !f.IsDir() {
				var q queryDef
				readFile(filepath.Join(path, f.Name()), &q)
				r = append(r, q)
			}
		}
	}
	return r
}
