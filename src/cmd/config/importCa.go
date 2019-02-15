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

package config

import (
	"fmt"
	"github.com/nlnwa/veidemannctl/src/configutil"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
)

// importCaCmd represents the import-ca command
var importCaCmd = &cobra.Command{
	Use:   "import-ca CA_CERT_FILE_NAME",
	Short: "Import file with trusted certificate chains for the idp and controller.",
	Long: `Import file with trusted certificate chains for the idp and controller. These are in addition to the default certs configured for the OS.`,
	Aliases: []string{"ca"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			rootCAs := args[0]
			if rootCAs != "" {
				rootCABytes, err := ioutil.ReadFile(rootCAs)
				if err != nil {
					log.Fatalf("failed to read root-ca: %v", err)
				}
				viper.Set("rootCAs", string(rootCABytes))
				configutil.WriteConfig()
			}
			fmt.Printf("Successfully imported root CA cert from file %s\n", rootCAs)
		} else {
			cmd.Usage()
		}
	},
}

func init() {
	ConfigCmd.AddCommand(importCaCmd)
}
