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
	commonsV1 "github.com/nlnwa/veidemann-api/go/commons/v1"
	"github.com/nlnwa/veidemannctl/src/apiutil"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	configV1 "github.com/nlnwa/veidemann-api/go/config/v1"
	"github.com/nlnwa/veidemannctl/src/connection"
	"github.com/nlnwa/veidemannctl/src/format"
	"golang.org/x/net/context"
)

var updateflags struct {
	label       string
	name        string
	filter      string
	updateField string
	pageSize    int32
}

// updateCmd represents the get command
var updateCmd = &cobra.Command{
	Use:   "update [object_type]",
	Short: "Update the value(s) for an object type",
	Long:  `Update one or many objects with new values.`,
	Example: `# Add CrawlJob for a seed.
veidemannctl update seed -n "http://www.gwpda.org/" -u seed.jobRef+=crawlJob:e46863ae-d076-46ca-8be3-8a8ef72e709

# Replace all configured CrawlJobs for a seed with a new one.
veidemannctl update seed -n "http://www.gwpda.org/" -u seed.jobRef=crawlJob:e46863ae-d076-46ca-8be3-8a8ef72e709`,

	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.ExactArgs(1)(cmd, args); err != nil {
			msg := "An object type must be specified. %sSee 'veidemannctl update -h' for help"
			return fmt.Errorf(msg, printValidObjectTypes())
		}
		if err := cobra.OnlyValidArgs(cmd, args); err != nil {
			msg := "%s is not a valid object type. %vSee 'veidemannctl update -h' for help"
			return fmt.Errorf(msg, args[0], printValidObjectTypes())
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		configClient, conn := connection.NewConfigClient()
		defer conn.Close()

		k := format.GetKind(args[0])

		var ids []string

		if len(args) == 2 {
			ids = args[1:]
		}

		if k == configV1.Kind_undefined {
			fmt.Printf("Unknown object type\n")
			_ = cmd.Usage()
			return
		}

		selector, err := apiutil.CreateListRequest(k, ids, updateflags.name, updateflags.label, updateflags.filter, updateflags.pageSize, 0)
		if err != nil {
			log.Fatalf("Error creating request: %v", err)
		}

		updateMask := new(commonsV1.FieldMask)
		updateTemplate := new(configV1.ConfigObject)
		if err := apiutil.CreateTemplateFilter(updateflags.updateField, updateTemplate, updateMask); err != nil {
			log.Fatalf("Error creating request: %v", err)
		}

		updateRequest := &configV1.UpdateRequest{
			ListRequest:    selector,
			UpdateMask:     updateMask,
			UpdateTemplate: updateTemplate,
		}

		u, err := configClient.UpdateConfigObjects(context.Background(), updateRequest)
		if err != nil {
			log.Fatalf("Error from controller: %v", err)
		}
		fmt.Printf("Objects updated: %v\n", u.Updated)
	},
	ValidArgs:     format.GetObjectNames(),
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	RootCmd.AddCommand(updateCmd)

	updateCmd.PersistentFlags().StringVarP(&updateflags.label, "label", "l", "", "List objects by label (<type>:<value> | <value>)")

	updateCmd.PersistentFlags().StringVarP(&updateflags.name, "name", "n", "", "List objects by name (accepts regular expressions)")
	annotation := make(map[string][]string)
	annotation[cobra.BashCompCustom] = []string{"__veidemannctl_get_name"}
	updateCmd.PersistentFlags().Lookup("name").Annotations = annotation

	updateCmd.PersistentFlags().StringVarP(&updateflags.filter, "filter", "q", "", "Filter objects by field (i.e. meta.description=foo)")
	updateCmd.PersistentFlags().StringVarP(&updateflags.updateField, "updateField", "u", "", "Filter objects by field (i.e. meta.description=foo)")
	_ = updateCmd.MarkPersistentFlagRequired("updateField")
	updateCmd.PersistentFlags().Int32VarP(&updateflags.pageSize, "limit", "s", 0, "Limit the number of objects to update. 0 = no limit")
}
