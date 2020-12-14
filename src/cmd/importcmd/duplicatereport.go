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
	configV1 "github.com/nlnwa/veidemann-api/go/config/v1"
	"github.com/nlnwa/veidemannctl/src/connection"
	"github.com/nlnwa/veidemannctl/src/importutil"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
	"os"
)

var dupFlags struct {
	outFile      string
	dbDir        string
	resetDb      bool
	toplevel     bool
	ignoreScheme bool
}

// duplicateReportCmd represents the duplicatereport command
var duplicateReportCmd = &cobra.Command{
	Use:   "duplicatereport [kind]",
	Short: "List duplicated seeds in Veidemann",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		kind := configV1.Kind(configV1.Kind_value[args[0]])

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
		var keyNormalizer importutil.KeyNormalizer
		if kind == configV1.Kind_seed {
			keyNormalizer = &UriKeyNormalizer{toplevel: dupFlags.toplevel, ignoreScheme: dupFlags.ignoreScheme}
		} else {
			keyNormalizer = &importutil.NoopKeyNormalizer{}
		}
		impf := importutil.NewImportDb(client, dupFlags.dbDir, kind, keyNormalizer, dupFlags.resetDb)
		impf.ImportExisting()
		defer impf.Close()

		switch kind {
		case configV1.Kind_seed:
			if err = impf.SeedDuplicateReport(out); err != nil {
				log.Errorf("failed creating seed duplicate report: %v", err)
			}
		case configV1.Kind_crawlEntity:
			if err = impf.CrawlEntityDuplicateReport(out); err != nil {
				log.Errorf("failed creating crawl entity duplicate report: %v", err)
			}
		}
	},
}

func init() {
	ImportCmd.AddCommand(duplicateReportCmd)

	duplicateReportCmd.PersistentFlags().StringVarP(&dupFlags.outFile, "outFile", "o", "", "File to write output.")
	duplicateReportCmd.PersistentFlags().StringVarP(&dupFlags.dbDir, "db-directory", "b", "/tmp/veidemannctl", "Directory for storing state db")
	duplicateReportCmd.PersistentFlags().BoolVarP(&dupFlags.resetDb, "reset-db", "r", false, "Clean state db")
	duplicateReportCmd.PersistentFlags().BoolVarP(&dupFlags.toplevel, "toplevel", "", false, "Convert URI to toplevel by removing path before checking for duplicates.")
	duplicateReportCmd.PersistentFlags().BoolVarP(&dupFlags.toplevel, "ignore-scheme", "", false, "Ignore the URL's scheme when checking for duplicates.")
}
