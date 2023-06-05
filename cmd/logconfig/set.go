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

type setLoggerCmdOptions struct {
	logger string
	level  configV1.LogLevels_Level
}

func (o *setLoggerCmdOptions) complete(cmd *cobra.Command, args []string) error {
	o.logger = args[0]

	l, ok := configV1.LogLevels_Level_value[args[1]]
	if !ok {
		return fmt.Errorf("invalid log level: %s", args[1])
	}

	level := configV1.LogLevels_Level(l)
	if level == configV1.LogLevels_UNDEFINED {
		return fmt.Errorf("invalid log level: %s", args[1])
	}

	return nil
}

func (o *setLoggerCmdOptions) run() error {
	conn, err := connection.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer conn.Close()

	client := configV1.NewConfigClient(conn)

	r, err := client.GetLogConfig(context.Background(), &empty.Empty{})
	if err != nil {
		return fmt.Errorf("could not get log config: %w", err)
	}

	loggers := make(map[string]configV1.LogLevels_Level)
	for _, l := range r.LogLevel {
		if l.Logger != "" {
			loggers[l.Logger] = l.Level
		}
	}

	loggers[o.logger] = o.level
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

func newSetLoggerCmd() *cobra.Command {

	return &cobra.Command{
		Use:   "set LOGGER LEVEL",
		Short: "Configure logger",
		Long:  `Configure logger.`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			o := &setLoggerCmdOptions{}
			if err := o.complete(cmd, args); err != nil {
				return err
			}

			// silence usage to avoid printing usage when an error occurs
			cmd.SilenceUsage = true

			return o.run()
		},
	}
}
