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

	"github.com/nlnwa/veidemannctl/src/configutil"
	"github.com/nlnwa/veidemannctl/src/connection"
	"github.com/spf13/cobra"
)

var manualLogin bool

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Initiate browser session for logging in to Veidemann",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		idp, _ := connection.GetIdp()
		if idp == "" {
			return
		}

		a := connection.NewAuth(idp)

		a.Login(manualLogin)
		claims := a.Claims()
		configutil.WriteConfig()
		fmt.Printf("Hello %s\n", claims.Name)
	},
}

func init() {
	RootCmd.AddCommand(loginCmd)
	loginCmd.PersistentFlags().BoolVarP(&manualLogin, "manual", "m", false,
		"Manually copy and paste login url and code. Use this to log in from a remote terminal.")
}
