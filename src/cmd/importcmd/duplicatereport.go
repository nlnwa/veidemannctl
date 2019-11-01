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
	"github.com/nlnwa/veidemannctl/src/connection"
	"github.com/nlnwa/veidemannctl/src/importutil"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
	"os"
)

var dupFlags struct {
	outFile         string
	dbDir           string
	resetDb         bool
}

// duplicateReportCmd represents the duplicatereport command
var duplicateReportCmd = &cobra.Command{
	Use:   "duplicatereport",
	Short: "List duplicated seeds in Veidemann",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Create output writer (file or stdout)
		var out io.Writer
		if dupFlags.outFile == "" {
			out = os.Stdout
		} else {
			out, err = os.Create(dupFlags.outFile)
			defer out.(io.Closer).Close()
			if err != nil {
				log.Fatalf("Unable to open out file: %v, cause: %v", dupFlags.outFile, err)
				os.Exit(1)
			}
		}

		// Create Veidemann config client
		client, conn := connection.NewConfigClient()
		defer conn.Close()

		// Create state Database based on seeds in Veidemann
		impf := importutil.NewImportDb(client, dupFlags.dbDir, dupFlags.resetDb)
		impf.ImportExisting()
		defer impf.Close()

		impf.DuplicateReport(out)
	},
}

func init() {
	ImportCmd.AddCommand(duplicateReportCmd)

	duplicateReportCmd.PersistentFlags().StringVarP(&dupFlags.outFile, "outFile", "o", "", "File to write output.")
	duplicateReportCmd.PersistentFlags().StringVarP(&dupFlags.dbDir, "db-directory", "b", "/tmp/veidemannctl", "Directory for storing state db")
	duplicateReportCmd.PersistentFlags().BoolVarP(&dupFlags.resetDb, "reset-db", "r", false, "Clean state db")
}