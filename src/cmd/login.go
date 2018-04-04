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

		authCodeURL := a.CreateAuthCodeURL()
		fmt.Println("A login screen should now open in your browser. Follow the login steps and paste the code here.")
		fmt.Println("In case the browser window won't open, paste this uri in a browser window:")
		fmt.Println(authCodeURL)
		a.Openbrowser(authCodeURL)
		fmt.Print("Code: ")
		var code string
		fmt.Scan(&code)
		a.VerifyCode(code)
		claims := a.Claims()
		configutil.WriteConfig()
		fmt.Printf("Hello %s\n", claims.Name)
	},
}

func init() {
	RootCmd.AddCommand(loginCmd)
}
