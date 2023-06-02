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
	configV1 "github.com/nlnwa/veidemann-api/go/config/v1"
	"github.com/nlnwa/veidemannctl/connection"
	"github.com/spf13/cobra"
)

type deleteLoggerCmdOptions struct {
	logger string
}

func (o *deleteLoggerCmdOptions) complete(cmd *cobra.Command, args []string) error {
	o.logger = args[0]

	return nil
}

func (o *deleteLoggerCmdOptions) run() error {
	conn, err := connection.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer conn.Close()

	client := configV1.NewConfigClient(conn)

	r, err := client.GetLogConfig(context.Background(), &empty.Empty{})
	if err != nil {
		return fmt.Errorf("failed to get log config: %w", err)
	}

	loggers := make(map[string]configV1.LogLevels_Level)
	for _, l := range r.LogLevel {
		if l.Logger != "" && l.Logger != o.logger {
			loggers[l.Logger] = l.Level
		}
	}

	n := &configV1.LogLevels{}
	for k, v := range loggers {
		n.LogLevel = append(n.LogLevel, &configV1.LogLevels_LogLevel{Logger: k, Level: v})
	}

	_, err = client.SaveLogConfig(context.Background(), n)
	if err != nil {
		return fmt.Errorf("failed to save log config: %w", err)
	}
	return nil
}

func newDeleteLoggerCmd() *cobra.Command {
	o := &deleteLoggerCmdOptions{}

	return &cobra.Command{
		Use:   "delete LOGGER",
		Short: "Delete a logger",
		Long:  `Delete a logger.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.complete(cmd, args); err != nil {
				return err
			}
			// silence usage to prevent printing usage when error occurs
			cmd.SilenceUsage = true

			return o.run()
		},
	}
}
