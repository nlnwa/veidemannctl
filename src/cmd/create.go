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

package cmd

import (
	"context"
	"fmt"
	configV1 "github.com/nlnwa/veidemann-api/go/config/v1"
	"github.com/nlnwa/veidemannctl/src/connection"
	"github.com/nlnwa/veidemannctl/src/format"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var filename string

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create or update a config object",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		if filename == "" {
			cmd.Usage()
			os.Exit(1)
		} else if filename == "-" {
			filename = ""
		}
		result, err := format.Unmarshal(filename)
		if err != nil {
			log.Fatalf("Parse error: %v", err)
			os.Exit(1)
		}

		client, conn := connection.NewConfigClient()
		defer conn.Close()

		for _, co := range result {
			if co.ApiVersion == "" {
				handleError(co, fmt.Errorf("Missing apiVersion"))
			}
			if co.Kind == configV1.Kind_undefined {
				handleError(co, fmt.Errorf("Missing kind"))
			}
			r, err := client.SaveConfigObject(context.Background(), co)
			handleError(co, err)
			fmt.Printf("Saved %v: %v %v\n", co.Kind, r.Meta.Name, r.Id)
		}
	},
}

func handleError(msg *configV1.ConfigObject, err error) {
	if err != nil {
		fmt.Printf("Could not save %v: %v. Cause: %v\n", msg.Kind, msg.Meta.Name, err)
		os.Exit(2)
	}
}

func init() {
	RootCmd.AddCommand(createCmd)

	createCmd.PersistentFlags().StringVarP(&filename, "filename", "f", "", "Filename or directory to read from. "+
		"If input is a directory, all files ending in .yaml or .json will be tried. An input of '-' will read from stdin.")
}
