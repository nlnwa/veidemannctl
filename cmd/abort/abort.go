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

package abort

import (
	"context"
	"fmt"

	controllerV1 "github.com/nlnwa/veidemann-api/go/controller/v1"
	"github.com/nlnwa/veidemannctl/connection"
	"github.com/spf13/cobra"
)

type abortCmdOptions struct {
	// ceids is a list of crawl execution ids to abort
	ceids []string
}

func (o *abortCmdOptions) complete(cmd *cobra.Command, args []string) error {
	o.ceids = args
	return nil
}

func (o *abortCmdOptions) run() error {
	conn, err := connection.Connect()
	if err != nil {
		return err
	}
	defer conn.Close()

	client := controllerV1.NewControllerClient(conn)

	for _, ceid := range o.ceids {
		request := controllerV1.ExecutionId{Id: ceid}
		_, err := client.AbortCrawlExecution(context.Background(), &request)
		if err != nil {
			return fmt.Errorf("failed to abort execution '%v': %w", ceid, err)
		}
	}
	return nil
}

func NewAbortCmd() *cobra.Command {
	o := &abortCmdOptions{}

	return &cobra.Command{
		GroupID: "run",
		Use:     "abort CRAWL-EXECUTION-ID ...",
		Short:   "Abort crawl executions",
		Long:    `Abort one or many crawl executions.`,
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = o.complete(cmd, args)

			// silence usage to avoid showing usage when an error occurs
			cmd.SilenceUsage = true

			return o.run()
		},
	}
}
