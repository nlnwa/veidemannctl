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
	"log"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/nlnwa/veidemann-api-go/controller/v1"
	"github.com/nlnwa/veidemannctl/src/connection"
	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get crawler status",
	Long:  `Get crawler status.`,
	Run: func(cmd *cobra.Command, args []string) {
		crawlerStatus, err := getCrawlerStatus()
		if err != nil {
			log.Fatalf("could not get crawler status: %v", err)
		}
		fmt.Printf("Status: %v, Url queue size: %v, Busy crawl host groups: %v\n",
			crawlerStatus.RunStatus, crawlerStatus.QueueSize, crawlerStatus.BusyCrawlHostGroupCount)
	},
}

func getCrawlerStatus() (*controller.CrawlerStatus, error) {
	client, conn := connection.NewControllerClient()
	defer func() {
		_ = conn.Close()
	}()

	request := empty.Empty{}
	return client.Status(context.Background(), &request)
}

func init() {
	RootCmd.AddCommand(statusCmd)
}
