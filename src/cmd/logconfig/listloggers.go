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

package logconfig

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/nlnwa/veidemannctl/src/connection"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// reportCmd represents the report command
var ListLoggersCmd = &cobra.Command{
	Use:   "list",
	Short: "List configured loggers",
	Long:  `List configured loggers.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, conn := connection.NewControllerClient()
		defer conn.Close()

		r, err := client.GetLogConfig(context.Background(), &empty.Empty{})
		if err != nil {
			log.Fatalf("could not get log config: %v", err)
		}

		fmt.Printf("%-45s %s\n", "Logger", "Level")
		fmt.Println("---------------------------------------------------")
		for _, l := range r.LogLevel {
			fmt.Printf("%-45s %s\n", l.Logger, l.Level)
		}
	},
}

func init() {
	LogconfigCmd.AddCommand(ListLoggersCmd)
}
