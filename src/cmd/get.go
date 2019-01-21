// Copyright © 2017 National Library of Norway.
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

package cmd

import (
	"fmt"
	"github.com/nlnwa/veidemannctl/src/apiutil"
	log "github.com/sirupsen/logrus"
	"io"
	"os"

	"github.com/spf13/cobra"

	configV1 "github.com/nlnwa/veidemann-api-go/config/v1"
	"github.com/nlnwa/veidemannctl/src/connection"
	"github.com/nlnwa/veidemannctl/src/format"
	"golang.org/x/net/context"
)

var flags struct {
	label      string
	name       string
	filter     string
	file       string
	format     string
	goTemplate string
	pageSize   int32
	page       int32
}

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get [object_type]",
	Short: "Get the value(s) for an object type",
	Long: `Display one or many objects.

` +
		printValidObjectTypes() +
		`Examples:
  #List all seeds.
  veidemannctl get seed

  #List all seeds in yaml output format.
  veidemannctl get seed -f yaml`,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			configClient, conn := connection.NewConfigClient()
			defer conn.Close()

			k := format.GetKind(args[0])

			var ids []string

			if len(args) == 2 {
				ids = args[1:]
			}

			if k == configV1.Kind_undefined {
				fmt.Printf("Unknown object type\n")
				cmd.Usage()
				return
			}

			request, err := apiutil.CreateListRequest(k, ids, flags.name, flags.label, flags.filter, flags.pageSize, flags.page)
			if err != nil {
				log.Fatalf("Error creating request: %v", err)
			}

			r, err := configClient.ListConfigObjects(context.Background(), request)
			if err != nil {
				log.Fatalf("Error from controller: %v", err)
			}

			s := format.NewFormatter(k, flags.file, flags.format, flags.goTemplate, "")
			defer s.Close()

			s.WriteHeader()
			for {
				msg, err := r.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					log.Fatalf("Error getting object: %v", err)
				}
				log.Debugf("Outputing record of kind '%s' with name '%s'", msg.Kind, msg.Meta.Name)
				if s.WriteRecord(msg) != nil {
					os.Exit(1)
				}
			}
		} else {
			fmt.Print("You must specify the object type to get. ")
			fmt.Println(printValidObjectTypes())
			fmt.Println("See 'veidemannctl get -h' for help")
		}
	},
	ValidArgs: format.GetObjectNames(),
}

func printValidObjectTypes() string {
	var names string
	for _, v := range format.GetObjectNames() {
		names += fmt.Sprintf("  * %s\n", v)
	}
	return fmt.Sprintf("Valid object types include:\n%s\n", names)
}

func init() {
	RootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	getCmd.PersistentFlags().StringVarP(&flags.label, "label", "l", "", "List objects by label (<type>:<value> | <value>)")

	getCmd.PersistentFlags().StringVarP(&flags.name, "name", "n", "", "List objects by name (accepts regular expressions)")
	annotation := make(map[string][]string)
	annotation[cobra.BashCompCustom] = []string{"__veidemannctl_get_name"}
	getCmd.PersistentFlags().Lookup("name").Annotations = annotation

	getCmd.PersistentFlags().StringVarP(&flags.filter, "filter", "q", "", "Filter objects by field (i.e. meta.description=foo")
	getCmd.PersistentFlags().StringVarP(&flags.format, "output", "o", "table", "Output format (table|json|yaml|template|template-file)")
	getCmd.PersistentFlags().StringVarP(&flags.goTemplate, "template", "t", "", "A Go template used to format the output")
	getCmd.PersistentFlags().StringVarP(&flags.file, "filename", "f", "", "File name to write to")
	getCmd.PersistentFlags().Int32VarP(&flags.pageSize, "pagesize", "s", 10, "Number of objects to get")
	getCmd.PersistentFlags().Int32VarP(&flags.page, "page", "p", 0, "The page number")
}
