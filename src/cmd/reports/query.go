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
	"github.com/ghodss/yaml"
	reportV1 "github.com/nlnwa/veidemann-api/go/report/v1"
	"github.com/nlnwa/veidemannctl/src/configutil"
	"github.com/nlnwa/veidemannctl/src/connection"
	"github.com/nlnwa/veidemannctl/src/format"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type queryConf struct {
	pageSize   int32
	page       int32
	goTemplate string
	file       string
	format     string
}

var queryFlags = &queryConf{}

// queryCmd represents the query command
var queryCmd = &cobra.Command{
	Use:   "query [queryString|file] args...",
	Short: "Run a database query",
	Long:  `Run a database query. The query should be a java script string like the ones used by RethinkDb javascript driver.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			client, conn := connection.NewReportClient()
			defer conn.Close()

			request := reportV1.ExecuteDbQueryRequest{}

			queryDef, err := getQueryDef(args[0])
			if err != nil {
				return err
			}
			defer queryDef.marshalSpec.Close()

			var params []interface{}

			for _, v := range args[1:] {
				params = append(params, v)
			}

			request.Query = fmt.Sprintf(queryDef.Query, params...)

			request.Limit = queryFlags.pageSize
			log.Debugf("Executing query: %s", request.GetQuery())

			// from now on we don't want usage when error occurs
			cmd.SilenceUsage = true

			stream, err := client.ExecuteDbQuery(context.Background(), &request)
			if err != nil {
				return fmt.Errorf("Failed executing query: %v", err)
			}

			for {
				value, err := stream.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					return fmt.Errorf("Query error: %v", err)
				}
				if err := queryDef.marshalSpec.WriteRecord(value.GetRecord()); err != nil {
					return err
				}
			}
		} else {
			d := configutil.GetConfigDir("query")
			if configutil.GlobalFlags.IsShellCompletion {
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
		return nil
	},
}

func init() {
	queryCmd.PersistentFlags().Int32VarP(&queryFlags.pageSize, "pagesize", "s", 10, "Number of objects to get")
	queryCmd.PersistentFlags().Int32VarP(&queryFlags.page, "page", "p", 0, "The page number")
	queryCmd.PersistentFlags().StringVarP(&queryFlags.format, "output", "o", "json", "Output format (json|yaml|template|template-file)")
	queryCmd.PersistentFlags().StringVarP(&queryFlags.goTemplate, "template", "t", "", "A Go template used to format the output")
	queryCmd.Flags().StringVarP(&queryFlags.file, "filename", "f", "", "File name to write to")
}

type queryDef struct {
	Name        string
	Description string
	Query       string
	Template    string
	marshalSpec format.Formatter
}

func getQueryDef(queryArg string) (queryDef, error) {
	var queryDef queryDef

	if strings.HasPrefix(queryArg, "r.") {
		queryDef.Query = queryArg
	} else {
		filename, err := findFile(queryArg)
		if err != nil {
			return queryDef, err
		}
		log.Debugf("Using query definition from file '%s'", filename)
		err = readFile(filename, &queryDef)
		if err != nil {
			return queryDef, err
		}
	}

	out, err := format.ResolveWriter(queryFlags.file)
	if err != nil {
		return queryDef, fmt.Errorf("Could not resolve output '%v': %v", queryFlags.file, err)
	}
	if queryDef.Template == "" {
		queryDef.marshalSpec, err = format.NewFormatter("", out, queryFlags.format, queryFlags.goTemplate)
	} else {
		queryDef.marshalSpec, err = format.NewFormatter("", out, "template", queryDef.Template)
	}
	if err != nil {
		return queryDef, err
	}
	return queryDef, nil
}

func findFile(name string) (string, error) {
	filename := name
	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		return filename, nil
	}

	queryDir := configutil.GetConfigDir("query")

	filename = filepath.Join(queryDir, name)
	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		return filename, nil
	}
	filename = filepath.Join(queryDir, name) + ".yml"
	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		return filename, nil
	}
	filename = filepath.Join(queryDir, name) + ".yaml"
	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		return filename, nil
	}
	return "", fmt.Errorf("query not found: %s", name)
}

func readFile(name string, queryDef *queryDef) error {
	data, err := ioutil.ReadFile(name)
	if err != nil {
		return fmt.Errorf("failed to read file: %s: %w", name, err)
	}
	// Found file
	if strings.HasSuffix(name, ".yml") || strings.HasSuffix(name, ".yaml") {
		yaml.Unmarshal(data, &queryDef)
	} else {
		queryDef.Query = string(data)
	}
	queryDef.Name = strings.TrimSuffix(filepath.Base(name), filepath.Ext(name))
	return nil
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
