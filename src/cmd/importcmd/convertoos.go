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
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	configV1 "github.com/nlnwa/veidemann-api-go/config/v1"
	"github.com/nlnwa/veidemannctl/src/connection"
	"github.com/nlnwa/veidemannctl/src/importutil"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/net/html"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var convertFlags struct {
	filename        string
	errorFile       string
	outFile         string
	toplevel        bool
	checkUri        bool
	checkUriTimeout int64
	dbDir           string
	resetDb         bool
}

// convertOosCmd represents the convertoos command
var convertOosCmd = &cobra.Command{
	Use:   "convertoos",
	Short: "Convert Out of Scope file(s) to seed import file",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		c := &converter{}
		var err error

		// Check inputflag
		if convertFlags.filename == "" {
			fmt.Printf("Import file is required. See --filename\n")
			cmd.Usage()
			os.Exit(1)
		}

		// Check ouput flag
		if convertFlags.outFile == "" {
			fmt.Printf("Output file is required. See --outfile\n")
			cmd.Usage()
			os.Exit(1)
		}

		// Create output file
		var out io.Writer
		out, err = os.Create(convertFlags.outFile)
		defer out.(io.Closer).Close()
		if err != nil {
			log.Fatalf("Unable to open out file: %v, cause: %v", convertFlags.outFile, err)
			os.Exit(1)
		}

		// Create error writer (file or stdout)
		var errFile io.Writer
		if convertFlags.errorFile == "" {
			errFile = os.Stdout
		} else {
			errFile, err = os.Create(convertFlags.errorFile)
			defer errFile.(io.Closer).Close()
			if err != nil {
				log.Fatalf("Unable to open error file: %v, cause: %v", convertFlags.errorFile, err)
				os.Exit(1)
			}
		}

		// Create Veidemann config client
		client, conn := connection.NewConfigClient()
		defer conn.Close()

		// Create http client
		c.httpClient = importutil.NewHttpClient(convertFlags.checkUriTimeout)

		// Create state Database based on seeds in Veidemann
		impf := importutil.NewImportDb(client, convertFlags.dbDir, convertFlags.resetDb)
		impf.ImportExisting()
		defer impf.Close()

		// Create Record reader for file input
		rr, err := importutil.NewRecordReader(convertFlags.filename, &importutil.LineAsStringDecoder{}, "*.txt")
		if err != nil {
			log.Fatalf("Parse error: %v", err)
			os.Exit(1)
		}

		// Processor for converting oos records into import records
		proc := func(value interface{}) error {
			v := value.(string)
			if exists, err := impf.Check(v); err != nil {
				return err
			} else if exists.Code > importutil.NEW {
				return fmt.Errorf("%v already exists", v)
			}

			seed := &seed{
				Uri:         v,
				SeedLabel:   []*configV1.Label{{Key: "source", Value: "oosh"}},
				EntityLabel: []*configV1.Label{{Key: "source", Value: "oosh"}},
			}

			err := c.checkUri(seed)
			if err != nil {
				return err
			}

			if seed.EntityName == "" {
				seed.EntityName = seed.Uri
			}

			json, err := json.Marshal(seed)
			if err != nil {
				return err
			}
			fmt.Fprintf(out, "%s\n", json)

			return nil
		}

		// Error handler
		errorHandler := func(state *importutil.StateVal) {
			_, _ = fmt.Fprintf(errFile, "ERR: %v %v %v\n", state.GetFilename(), state.GetRecordNum(), state.GetError())
		}

		// Create multithreaded executor
		conv := importutil.NewExecutor(512, proc, errorHandler)

		// Process
		var ts string
		for {
			state, err := rr.Next(&ts)
			if err == io.EOF {
				break
			}
			if err != nil {
				_, _ = fmt.Fprintf(errFile, "Error decoding record: %v, cause: %v", state, err)
				os.Exit(1)
			}
			conv.Do(state, ts)
		}

		conv.Finish()
		_, _ = fmt.Fprintf(os.Stderr, "\nRecords read: %v, imported: %v, Failed: %v\n", conv.Count, conv.Success, conv.Failed)
	},
}

func init() {
	ImportCmd.AddCommand(convertOosCmd)

	convertOosCmd.PersistentFlags().StringVarP(&convertFlags.filename, "filename", "f", "", "Filename or directory to read from. "+
		"If input is a directory, all files ending in .yaml or .json will be tried. An input of '-' will read from stdin. (required)")
	convertOosCmd.PersistentFlags().StringVarP(&convertFlags.errorFile, "errorfile", "e", "", "File to write errors to.")
	convertOosCmd.PersistentFlags().StringVarP(&convertFlags.outFile, "outfile", "o", "", "File to write result to. (required)")
	convertOosCmd.PersistentFlags().BoolVarP(&convertFlags.toplevel, "toplevel", "", true, "Convert URI to toplevel by removing path.")
	convertOosCmd.PersistentFlags().BoolVarP(&convertFlags.checkUri, "checkuri", "", true, "Check the uri for liveness and follow 301")
	convertOosCmd.PersistentFlags().Int64VarP(&convertFlags.checkUriTimeout, "checkuri-timeout", "", 2000, "Timeout in ms when checking uri for liveness.")
	convertOosCmd.PersistentFlags().StringVarP(&convertFlags.dbDir, "db-directory", "b", "/tmp/veidemannctl", "Directory for storing state db")
	convertOosCmd.PersistentFlags().BoolVarP(&convertFlags.resetDb, "reset-db", "r", false, "Clean state db")
}

type converter struct {
	httpClient *http.Client
}

func (c *converter) checkUri(s *seed) (err error) {
	uri, err := url.Parse(s.Uri)
	if err != nil {
		return fmt.Errorf("unparseable URL '%v', cause: %v", s.Uri, err)
	}

	if uri.Host == "" {
		return errors.New("unparseable URL")
	}

	if convertFlags.toplevel {
		uri.Path = ""
		uri.RawQuery = ""
		uri.Fragment = ""
		s.Uri = uri.Scheme + "://" + uri.Host
	}

	if convertFlags.checkUri {
		c.checkRedirect(s.Uri, s, 0)
	}
	return
}

func (c *converter) checkRedirect(uri string, s *seed, count int) {
	if count > 5 {
		return
	}
	count++

	resp, err := c.httpClient.Head(uri)
	if err != nil {
		uerr := err.(*url.Error)
		if uerr.Timeout() {
			err = nil
		} else {
			switch v := uerr.Err.(type) {
			case *net.OpError:
				if t, ok := v.Err.(*net.DNSError); ok {
					err = fmt.Errorf("no such host %s", t.Name)
				}
				return
			case x509.HostnameError:
			case x509.UnknownAuthorityError:
			case x509.CertificateInvalidError:
				return
			default:
				return
			}
		}
	} else {
		resp.Body.Close()
		if resp.StatusCode == 301 {
			uri = resp.Header.Get("Location")
			if uri != "" {
				c.checkRedirect(uri, s, count)
			}
		} else {
			s.EntityName = c.getTitle(uri)
		}
	}
}

func (c *converter) getTitle(uri string) string {
	resp, err := c.httpClient.Get(uri)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return ""
	}
	var title string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" && n.FirstChild != nil {
			title = strings.TrimSpace(n.FirstChild.Data)
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return title
}
