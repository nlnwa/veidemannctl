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

package unpause

import (
	"context"

	controllerV1 "github.com/nlnwa/veidemann-api/go/controller/v1"
	"github.com/nlnwa/veidemannctl/connection"
	"github.com/spf13/cobra"
	empty "google.golang.org/protobuf/types/known/emptypb"
)

func NewCmd() *cobra.Command {
	return &cobra.Command{
		GroupID:      "status",
		Use:          "unpause",
		Short:        "Request crawler to unpause",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			conn, err := connection.Connect()
			if err != nil {
				return err
			}
			defer conn.Close()

			client := controllerV1.NewControllerClient(conn)

			_, err = client.UnPauseCrawler(context.Background(), &empty.Empty{})
			return err
		},
	}
}
