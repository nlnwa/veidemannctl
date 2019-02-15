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
	"github.com/spf13/cobra"
	"os"
	log "github.com/sirupsen/logrus"
	api "github.com/nlnwa/veidemann-api-go/veidemann_api"
)

// DeleteCmd represents the delete command
var DeleteLoggerCmd = &cobra.Command{
	Use:   "delete [logger]",
	Short: "Delete a logger",
	Long:  `Delete a logger.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Printf("Wrong number of arguments\n")
			cmd.Help()
			os.Exit(1)
		}

		logger := args[0]

		client, conn := connection.NewControllerClient()
		defer conn.Close()

		r, err := client.GetLogConfig(context.Background(), &empty.Empty{})
		if err != nil {
			log.Fatalf("could not get log config: %v", err)
		}

		var loggers map[string]api.LogLevels_Level
		loggers = make(map[string]api.LogLevels_Level)
		for _, l := range r.LogLevel {
			if l.Logger != "" && l.Logger != logger {
				loggers[l.Logger] = l.Level
			}
		}

		n := &api.LogLevels{}
		for k, v := range loggers {
			n.LogLevel = append(n.LogLevel, &api.LogLevels_LogLevel{Logger: k, Level: v})
		}

		_, err = client.SaveLogConfig(context.Background(), n)
		if err != nil {
			log.Fatalf("could not get log config: %v", err)
		}
	},
}

func init() {
	LogconfigCmd.AddCommand(DeleteLoggerCmd)
}