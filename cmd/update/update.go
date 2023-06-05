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

package update

import (
	"context"
	"fmt"

	commonsV1 "github.com/nlnwa/veidemann-api/go/commons/v1"
	configV1 "github.com/nlnwa/veidemann-api/go/config/v1"
	"github.com/nlnwa/veidemannctl/apiutil"
	"github.com/nlnwa/veidemannctl/connection"
	"github.com/nlnwa/veidemannctl/format"
	"github.com/spf13/cobra"
)

type updateCmdOptions struct {
	kind        configV1.Kind
	ids         []string
	name        string
	label       string
	filters     []string
	updateField string
	pageSize    int32
}

func (o *updateCmdOptions) complete(cmd *cobra.Command, args []string) error {
	// first arg is kind
	k := args[0]
	o.kind = format.GetKind(k)

	if o.kind == configV1.Kind_undefined {
		return fmt.Errorf("undefined kind '%s'", k)
	}

	// rest of args are ids
	o.ids = args[1:]

	return nil
}

// run runs the update command
func (o *updateCmdOptions) run() error {
	conn, err := connection.Connect()
	if err != nil {
		return err
	}
	defer conn.Close()

	client := configV1.NewConfigClient(conn)

	selector, err := apiutil.CreateListRequest(o.kind, o.ids, o.name, o.label, o.filters, o.pageSize, 0)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	updateMask := new(commonsV1.FieldMask)
	updateTemplate := new(configV1.ConfigObject)
	if err := apiutil.CreateTemplateFilter(o.updateField, updateTemplate, updateMask); err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	updateRequest := &configV1.UpdateRequest{
		ListRequest:    selector,
		UpdateMask:     updateMask,
		UpdateTemplate: updateTemplate,
	}

	u, err := client.UpdateConfigObjects(context.Background(), updateRequest)
	if err != nil {
		return fmt.Errorf("error from controller: %w", err)
	}
	fmt.Printf("Objects updated: %v\n", u.Updated)
	return nil
}

func NewUpdateCmd() *cobra.Command {
	o := &updateCmdOptions{}

	cmd := &cobra.Command{
		GroupID: "basic",
		Use:     "update KIND [ID ...]",
		Short:   "Update fields of config objects of the same kind",
		Long:    `Update a field of one or many config objects of the same kind`,
		Example: `# Add CrawlJob for a seed.
veidemannctl update seed -n "https://www.gwpda.org/" -u seed.jobRef+=crawlJob:e46863ae-d076-46ca-8be3-8a8ef72e709

# Replace all configured CrawlJobs for a seed with a new one.
veidemannctl update seed -n "https://www.gwpda.org/" -u seed.jobRef=crawlJob:e46863ae-d076-46ca-8be3-8a8ef72e709`,

		Args: cobra.MatchAll(
			cobra.MinimumNArgs(1),
			func(cmd *cobra.Command, args []string) error {
				return cobra.OnlyValidArgs(cmd, args[:1])
			},
		),
		ValidArgs: format.GetObjectNames(),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.complete(cmd, args); err != nil {
				return err
			}
			// silence usage to prevent printing usage when an error occurs
			cmd.SilenceUsage = true
			return o.run()
		},
	}

	// update-field is required
	cmd.Flags().StringVarP(&o.updateField, "update-field", "u", "", "Which field to update (i.e. meta.description=foo)")
	_ = cmd.MarkFlagRequired("update-field")

	// label is optional
	cmd.Flags().StringVarP(&o.label, "label", "l", "", "Filter objects by label (<type>:<value> | <value>)")

	// name is optional
	cmd.Flags().StringVarP(&o.name, "name", "n", "", "Filter objects by name (accepts regular expressions)")
	// register name flag completion func
	_ = cmd.RegisterFlagCompletionFunc("name", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		names, err := apiutil.CompleteName(args[0], toComplete)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		return names, cobra.ShellCompDirectiveDefault
	})

	// filters is optional
	cmd.Flags().StringArrayVarP(&o.filters, "filter", "q", nil, "Filter objects by field (i.e. meta.description=foo)")

	// limit is optional
	cmd.Flags().Int32VarP(&o.pageSize, "limit", "s", 0, "Limit the number of objects to update. 0 = no limit")

	return cmd
}
