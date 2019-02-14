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
	configV1 "github.com/nlnwa/veidemann-api-go/config/v1"
	"github.com/nlnwa/veidemannctl/src/apiutil"
	"github.com/nlnwa/veidemannctl/src/connection"
	"github.com/nlnwa/veidemannctl/src/format"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"io"
)

var deleteFlags struct {
	label  string
	filter string
	dryRun bool
}

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete <kind> [id] ...",
	Short: "Delete a config object",
	Long: `Delete a config object.

` +
		printValidObjectTypes() +
		`Examples:
  #Delete a seed.
  veidemannctl delete seed 407a9600-4f25-4f17-8cff-ee1b8ee950f6`,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			configClient, conn := connection.NewConfigClient()
			defer conn.Close()

			k := format.GetKind(args[0])

			var ids []string

			if len(args) > 1 {
				ids = args[1:]
			}

			if k == configV1.Kind_undefined {
				fmt.Printf("Unknown object type\n")
				cmd.Usage()
				return
			}

			if len(ids) == 0 && deleteFlags.filter == "" && deleteFlags.label == "" {
				fmt.Printf("At least one of the -f or -l flags must be set or one or more id's\n")
				cmd.Usage()
				return
			}

			selector, err := apiutil.CreateListRequest(k, ids, "", deleteFlags.label, deleteFlags.filter, 0, 0)
			if err != nil {
				log.Fatalf("Error creating request: %v", err)
			}
			if err != nil {
				log.Fatalf("Error creating request: %v", err)
			}

			r, err := configClient.ListConfigObjects(context.Background(), selector)
			if err != nil {
				log.Fatalf("Error from controller: %v", err)
			}

			count, err := configClient.CountConfigObjects(context.Background(), selector)
			if err != nil {
				log.Fatalf("Error from controller: %v", err)
			}

			if deleteFlags.dryRun {
				for {
					msg, err := r.Recv()
					if err == io.EOF {
						break
					}
					if err != nil {
						log.Fatalf("Error getting object: %v", err)
					}
					log.Debugf("Outputing record of kind '%s' with name '%s'", msg.Kind, msg.Meta.Name)
					fmt.Printf("%s\n", msg.Meta.Name)
				}
				fmt.Printf("Requested count: %v\nTo actually delete, add: --dry-run=false\n", count.Count)
			} else {
				var deleted int
				for {
					msg, err := r.Recv()
					if err == io.EOF {
						break
					}
					if err != nil {
						log.Fatalf("Error getting object: %v", err)
					}
					log.Debugf("Deleting record of kind '%s' with name '%s'", msg.Kind, msg.Meta.Name)

					request := &configV1.ConfigObject{
						ApiVersion: "v1",
						Kind:       k,
						Id:         msg.Id,
					}

					r, err := configClient.DeleteConfigObject(context.Background(), request)
					if err != nil {
						log.Fatalf("could not delete '%v': %v\n", msg.Id, err)
					}
					if r.Deleted {
						deleted++
						fmt.Print(".")
					} else {
						fmt.Printf("\nCould not delete %v: %v\n", k, msg.Id)
					}
				}
				fmt.Printf("\nRequested count: %v\n", count.Count)
				fmt.Printf("Deleted count: %v\n", deleted)
			}
		} else {
			fmt.Println("You must specify the object type to delete. ")
			for _, k := range configV1.Kind_name {
				fmt.Println(k)
			}
			fmt.Println("See 'veidemannctl get -h' for help")
		}
	},
	ValidArgs:
	format.GetObjectNames(),
}

func init() {
	RootCmd.AddCommand(deleteCmd)

	deleteCmd.PersistentFlags().StringVarP(&deleteFlags.label, "label", "l", "", "Delete objects by label (<type>:<value> | <value>)")
	deleteCmd.PersistentFlags().StringVarP(&deleteFlags.filter, "filter", "q", "", "Delete objects by field (i.e. meta.description=foo)")
	deleteCmd.PersistentFlags().BoolVarP(&deleteFlags.dryRun, "dry-run", "", true, "Set to false to execute delete")
}
