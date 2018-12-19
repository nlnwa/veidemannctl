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
	"fmt"

	tm "github.com/buger/goterm"
	"github.com/golang/protobuf/ptypes"
	"github.com/nlnwa/veidemannctl/src/connection"
	api "github.com/nlnwa/veidemann-api-go/veidemann_api"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"io"
	"log"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get a realtime view of crawling activity",
	Long:  `Get a realtime view of crawling activity.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, conn := connection.NewStatusClient()
		defer conn.Close()

		request := api.RunningExecutionsRequest{PageSize: 50}

		stream, err := client.GetRunningExecutions(context.Background(), &request)
		if err != nil {
			log.Fatalf("could not get status: %v", err)
		}

		tm.Clear()
		for {
			value, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("%v.ListFeatures(_) = _, %v", client, err)
			}

			tm.MoveCursor(1, 1)

			rest := tm.Width() - 142
			tm.Println(
				tm.Background(
					tm.Color(
						fmt.Sprintf("%-36s %-18.18s %-10.10s %5s %5s %10s %6s %5s %5s %5s %5s %-20s %-*s",
							"ID", "Start time", "State", "Docs", "URIs", "Bytes", "Failed",
							"OoS", "Deny", "Retry", "Queue", "Seed", rest, "Current Uri"),
						tm.BLACK),
					tm.WHITE))

			for _, e := range value.GetValue() {
				s, _ := ptypes.Timestamp(e.StartTime)
				start := s.Format("02.01.06 15:04:05Z")
				color := tm.BLACK
				switch e.State {
				case api.CrawlExecutionStatus_FETCHING:
					color = tm.GREEN
				case api.CrawlExecutionStatus_SLEEPING:
					color = tm.BLUE
				}
				tm.Println(tm.Background(
					fmt.Sprintf("%-36.36s %18s %-10.10s %5d %5d %10d %6d %5d %5d %5d %5d %-20.20s %-*.*s",
						e.Id, start, e.State, e.DocumentsCrawled, e.UrisCrawled, e.BytesCrawled, e.DocumentsFailed,
						e.DocumentsOutOfScope, e.DocumentsDenied, e.DocumentsRetried, e.QueueSize, e.Seed,
						rest, rest, e.CurrentUri),
					color))
			}
			tm.Flush()
		}
	},
}

func init() {
	RootCmd.AddCommand(statusCmd)
}
