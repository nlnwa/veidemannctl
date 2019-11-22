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
	"github.com/spf13/cobra"

	configV1 "github.com/nlnwa/veidemann-api-go/config/v1"
	"github.com/nlnwa/veidemannctl/src/connection"
	"github.com/nlnwa/veidemannctl/src/format"
	"golang.org/x/net/context"
)

var updateflags struct {
	label       string
	name        string
	filter      string
	updateField string
	file        string
	format      string
	goTemplate  string
	pageSize    int32
	page        int32
}

// updateCmd represents the get command
var updateCmd = &cobra.Command{
	Use:   "update [object_type]",
	Short: "Update the value(s) for an object type",
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

			selector, err := apiutil.CreateListRequest(k, ids, updateflags.name, updateflags.label, updateflags.filter, updateflags.pageSize, updateflags.page)
			if err != nil {
				log.Fatalf("Error creating request: %v", err)
			}

			mask, obj, err := apiutil.CreateTemplateFilter(updateflags.updateField, &configV1.ConfigObject{})
			if err != nil {
				log.Fatalf("Error creating request: %v", err)
			}

			updateRequest := &configV1.UpdateRequest{
				ListRequest:    selector,
				UpdateMask:     mask,
				UpdateTemplate: obj.(*configV1.ConfigObject),
			}

			u, err := configClient.UpdateConfigObjects(context.Background(), updateRequest)
			if err != nil {
				log.Fatalf("Error from controller: %v", err)
			}
			fmt.Printf("Objects updated: %v\n", u.Updated)
		} else {
			fmt.Print("You must specify the object type to update. ")
			fmt.Println(printValidObjectTypes())
			fmt.Println("See 'veidemannctl get -h' for help")
		}
	},
	ValidArgs: format.GetObjectNames(),
}

func init() {
	RootCmd.AddCommand(updateCmd)

	updateCmd.PersistentFlags().StringVarP(&updateflags.label, "label", "l", "", "List objects by label (<type>:<value> | <value>)")

	updateCmd.PersistentFlags().StringVarP(&updateflags.name, "name", "n", "", "List objects by name (accepts regular expressions)")
	annotation := make(map[string][]string)
	annotation[cobra.BashCompCustom] = []string{"__veidemannctl_get_name"}
	updateCmd.PersistentFlags().Lookup("name").Annotations = annotation

	updateCmd.PersistentFlags().StringVarP(&updateflags.filter, "filter", "q", "", "Filter objects by field (i.e. meta.description=foo")
	updateCmd.PersistentFlags().StringVarP(&updateflags.updateField, "updateField", "u", "", "Filter objects by field (i.e. meta.description=foo")
	updateCmd.PersistentFlags().StringVarP(&updateflags.format, "output", "o", "table", "Output format (table|json|yaml|template|template-file)")
	updateCmd.PersistentFlags().StringVarP(&updateflags.goTemplate, "template", "t", "", "A Go template used to format the output")
	updateCmd.PersistentFlags().StringVarP(&updateflags.file, "filename", "f", "", "File name to write to")
	updateCmd.PersistentFlags().Int32VarP(&updateflags.pageSize, "pagesize", "s", 10, "Number of objects to get")
	updateCmd.PersistentFlags().Int32VarP(&updateflags.page, "page", "p", 0, "The page number")
}
