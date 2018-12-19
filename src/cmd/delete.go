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
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	configV1 "github.com/nlnwa/veidemann-api-go/config/v1"
	"github.com/nlnwa/veidemannctl/src/connection"
	"github.com/nlnwa/veidemannctl/src/format"
	"golang.org/x/net/context"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete [kind] [id]",
	Short: "Delete a config object",
	Long: `Delete a config object.

` +
		printValidObjectTypes() +
		`Examples:
  #List all seeds.
  veidemannctl delete seed 407a9600-4f25-4f17-8cff-ee1b8ee950f6`,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			configClient, conn := connection.NewConfigClient()
			defer conn.Close()

			k := configV1.Kind(configV1.Kind_value[args[0]])

			if k == configV1.Kind_undefined {
				fmt.Printf("Unknown object type\n")
				cmd.Usage()
				return
			}

			//var selector []string
			var ids []string

			if len(args) > 1 {
				ids = args[1:]
				fmt.Println("ID: ", ids)

				fmt.Println("KIND: ", k, args[0], flags.format)
				for _, id := range ids {
					request := &configV1.ConfigObject{
						ApiVersion: "v1",
						Kind:       k,
						Id:         id,
					}

					r, err := configClient.DeleteConfigObject(context.Background(), request)
					if err != nil {
						log.Fatalf("could not get crawl config: %v", err)
					}
					fmt.Println(r)
				}
			} else {
				fmt.Println("Missing id(s)")
				fmt.Println("See 'veidemannctl get -h' for help")
			}
		} else {
			fmt.Println("You must specify the object type to get. ")
			//fmt.Println(printValidObjectTypes())
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

	// Here you will define your flags and configuration settings.

}
