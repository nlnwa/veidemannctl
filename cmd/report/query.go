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
	"os"
	"path/filepath"
	"strings"

	"github.com/invopop/yaml"
	reportV1 "github.com/nlnwa/veidemann-api/go/report/v1"
	"github.com/nlnwa/veidemannctl/config"
	"github.com/nlnwa/veidemannctl/connection"
	"github.com/nlnwa/veidemannctl/format"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type queryCmdOptions struct {
	queryOrFile string
	queryArgs   []any
	pageSize    int32
	page        int32
	goTemplate  string
	file        string
	format      string
}

func (o *queryCmdOptions) complete(cmd *cobra.Command, args []string) error {
	o.queryOrFile = args[0]

	for _, arg := range args[1:] {
		o.queryArgs = append(o.queryArgs, arg)
	}

	return nil
}

// run runs the query command
func (o *queryCmdOptions) run() error {
	q, err := o.parseQuery()
	if err != nil {
		return fmt.Errorf("failed to parse query: %w", err)
	}

	request := reportV1.ExecuteDbQueryRequest{
		Query: q.query,
		Limit: o.pageSize,
	}

	log.Debug().Msgf("Executing query: %s", request.GetQuery())

	var w io.Writer

	if o.file == "" || o.file == "-" {
		w = os.Stdout
	} else {
		w, err = os.Create(o.file)
		if err != nil {
			return fmt.Errorf("failed to create file: %s: %w", o.file, err)
		}
	}

	var formatter format.Formatter
	if q.template == "" {
		formatter, err = format.NewFormatter("", w, o.format, o.goTemplate)
	} else {
		formatter, err = format.NewFormatter("", w, "template", q.template)
	}
	if err != nil {
		return fmt.Errorf("failed to create formatter: %w", err)
	}

	conn, err := connection.Connect()
	if err != nil {
		return err
	}
	defer conn.Close()

	client := reportV1.NewReportClient(conn)

	stream, err := client.ExecuteDbQuery(context.Background(), &request)
	if err != nil {
		return fmt.Errorf("failed executing query: %w", err)
	}

	for {
		value, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to receive a value: %w", err)
		}
		if err := formatter.WriteRecord(value.GetRecord()); err != nil {
			return err
		}
	}

	return nil
}

func newQueryCmd() *cobra.Command {
	o := &queryCmdOptions{}

	cmd := &cobra.Command{
		Use:   "query (QUERY-STRING | FILENAME) [ARGS ...]",
		Short: "Run a database query",
		Long:  `Run a database query. The query should be a JavaScript string like the ones used by RethinkDB JavaScript driver.`,
		Args:  cobra.MinimumNArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			d, err := config.GetConfigPath("query")
			if err != nil {
				return nil, cobra.ShellCompDirectiveError
			}
			q := listStoredQueries(d)
			var names []string
			for _, s := range q {
				if strings.HasPrefix(s.name, toComplete) {
					names = append(names, s.name+"\t"+s.description)
				}
			}
			return names, cobra.ShellCompDirectiveDefault
		},
		PreRunE: o.complete,
		RunE: func(cmd *cobra.Command, args []string) error {
			// set silence usage to true to avoid printing usage when an error occurs
			cmd.SilenceUsage = true

			return o.run()
		},
	}

	cmd.Flags().Int32VarP(&o.pageSize, "pagesize", "s", 10, "Number of objects to get")
	cmd.Flags().Int32VarP(&o.page, "page", "p", 0, "The page number")
	cmd.Flags().StringVarP(&o.format, "output", "o", "json", "Output format (json|yaml|template|template-file)")
	cmd.Flags().StringVarP(&o.goTemplate, "template", "t", "", "A Go template used to format the output")
	cmd.Flags().StringVarP(&o.file, "filename", "f", "", "Filename to write to")

	return cmd
}

// query is a struct for holding query definitions
type query struct {
	name        string
	description string
	query       string
	template    string
}

func (o *queryCmdOptions) parseQuery() (*query, error) {
	var q *query

	if strings.HasPrefix(o.queryOrFile, "r.") {
		q = &query{
			query: o.queryOrFile,
		}
	} else {
		filename, err := findQueryFile(o.queryOrFile)
		if err != nil {
			return nil, err
		}
		log.Debug().Msgf("Using query definition from file '%s'", filename)
		q, err = readQuery(filename)
		if err != nil {
			return nil, err
		}
	}

	q.query = fmt.Sprintf(q.query, o.queryArgs...)

	return q, nil
}

// findQueryFile finds a query file
func findQueryFile(name string) (string, error) {
	filename := name
	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		return filename, nil
	}

	queryDir, err := config.GetConfigPath("query")
	if err != nil {
		return "", fmt.Errorf("failed to resolve query dir: %w", err)
	}

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

// readQuery reads a query definition from a file
func readQuery(name string) (*query, error) {
	data, err := os.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %s: %w", name, err)
	}
	qd := new(query)
	// Found file
	if strings.HasSuffix(name, ".yml") || strings.HasSuffix(name, ".yaml") {
		err := yaml.Unmarshal(data, qd)
		if err != nil {
			return nil, err
		}
	} else {
		qd.query = string(data)
	}
	qd.name = strings.TrimSuffix(filepath.Base(name), filepath.Ext(name))

	return qd, nil
}

// listStoredQueries returns a list of query definitions stored in the query directory
func listStoredQueries(path string) []*query {
	var r []*query

	if files, err := os.ReadDir(path); err == nil {
		for _, f := range files {
			if !f.IsDir() {
				q, err := readQuery(filepath.Join(path, f.Name()))
				if err != nil {
					return nil
				}
				r = append(r, q)
			}
		}
	}

	return r
}
