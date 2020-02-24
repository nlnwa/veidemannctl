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
	"log"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/nlnwa/veidemannctl/src/connection"
	"github.com/spf13/cobra"
)

// pauseCmd represents the pause command
var pauseCmd = &cobra.Command{
	Use:   "pause",
	Short: "Pause crawler",
	Long:  `Pause crawler`,
	Run: func(cmd *cobra.Command, args []string) {
		err := pauseCrawler()
		if err != nil {
			log.Fatalf("could not pause crawler: %v", err)
		}
	},
}

func pauseCrawler() error {
	client, conn := connection.NewControllerClient()
	defer func() {
		_ = conn.Close()
	}()

	_, err := client.PauseCrawler(context.Background(), &empty.Empty{})
	return err
}

func init() {
	RootCmd.AddCommand(pauseCmd)
}
