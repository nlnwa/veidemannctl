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
	controllerV1 "github.com/nlnwa/veidemann-api/go/controller/v1"

	"context"
	"github.com/nlnwa/veidemannctl/src/connection"
	"github.com/spf13/cobra"
	"log"
)

// abortCmd represents the abort command
var abortJobExecutionCmd = &cobra.Command{
	Use:   "abortjobexecution",
	Short: "Abort one or more job executions",
	Long:  `Abort one or more job executions.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			client, conn := connection.NewControllerClient()
			defer conn.Close()

			for _, arg := range args {
				request := controllerV1.ExecutionId{Id: arg}
				_, err := client.AbortJobExecution(context.Background(), &request)
				if err != nil {
					log.Fatalf("could not abort job execution '%v': %v", arg, err)
				}
			}
		} else {
			cmd.Usage()
		}
	},
}

func init() {
	RootCmd.AddCommand(abortJobExecutionCmd)
}
