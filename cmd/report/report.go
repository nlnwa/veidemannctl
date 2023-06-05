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

package report

import (
	"github.com/spf13/cobra"
)

func NewReportCmd() *cobra.Command {
	cmd := &cobra.Command{
		GroupID: "advanced",
		Use:     "report",
		Short:   "Request a report",
	}

	cmd.AddCommand(newJobExecutionCmd())   // jobexecution
	cmd.AddCommand(newCrawlExecutionCmd()) // crawlexecution
	cmd.AddCommand(newCrawlLogCmd())       // crawllog
	cmd.AddCommand(newQueryCmd())          // query
	cmd.AddCommand(newPageLogCmd())        // pagelog

	return cmd
}
