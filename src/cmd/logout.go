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
	"github.com/nlnwa/veidemannctl/src/configutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out of Veidemann",
	Long:  `Log out of Veidemann.`,
	Run: func(cmd *cobra.Command, args []string) {
		viper.Set("accessToken", "")
		viper.Set("nonce", "")
		configutil.WriteConfig()
	},
}

func init() {
	RootCmd.AddCommand(logoutCmd)
}
