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

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/nlnwa/veidemannctl/src/connection"
	"golang.org/x/net/context"
)

// activerolesCmd represents the activeroles command
var activerolesCmd = &cobra.Command{
	Use:   "activeroles",
	Short: "Get the active roles for the currently logged in user",
	Long:  `Get the active roles for the currently logged in user.`,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			client, conn := connection.NewControllerClient()
			defer conn.Close()

			r, err := client.GetRolesForActiveUser(context.Background(), &empty.Empty{})
			if err != nil {
				log.Fatalf("could not get active role: %v", err)
			}

			for _, role := range r.Role {
				fmt.Println(role)
			}
		} else {
			fmt.Println("activeroles takes no arguments.")
			fmt.Println("See 'veidemannctl activeroles -h' for help")
		}
	},
}

func init() {
	RootCmd.AddCommand(activerolesCmd)
}
