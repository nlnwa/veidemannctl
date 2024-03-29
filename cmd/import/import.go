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

package importcmd

import (
	"github.com/nlnwa/veidemannctl/cmd/import/convertoos"
	"github.com/nlnwa/veidemannctl/cmd/import/duplicatereport"
	"github.com/nlnwa/veidemannctl/cmd/import/seeds"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		GroupID: "advanced",
		Use:     "import",
		Short:   "Import data into Veidemann using subcommands",
		Long:    ``,
	}

	cmd.AddCommand(convertoos.NewCmd())      // convertoos
	cmd.AddCommand(seeds.NewCmd())           // seed
	cmd.AddCommand(duplicatereport.NewCmd()) // duplicate

	return cmd
}
