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
	"errors"
	"fmt"
	configV1 "github.com/nlnwa/veidemann-api-go/config/v1"
	"github.com/nlnwa/veidemannctl/src/connection"
	"github.com/nlnwa/veidemannctl/src/importutil"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
)

type seedDesc struct {
	EntityName        string
	Uri               string
	EntityLabel       []*configV1.Label
	SeedLabel         []*configV1.Label
	EntityDescription string
	SeedDescription   string
	Description       string
	fileName          string
	recNum            int
	err               error
	crawlJobRef       []*configV1.ConfigRef
}

var importFlags struct {
	filename        string
	errorFile       string
	toplevel        bool
	checkUri        bool
	checkUriTimeout int64
	crawlJobId      string
	dbDir           string
	resetDb         bool
}

var httpClient *http.Client

// importSeedCmd represents the import command
var importSeedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Import seeds",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		i := &importer{}
		var err error

		// Check inputflag
		if importFlags.filename == "" {
			fmt.Printf("Import file is required. See --filename\n")
			cmd.Usage()
			os.Exit(1)
		}

		// Create error writer (file or stdout)
		var errFile io.Writer
		if importFlags.errorFile == "" {
			errFile = os.Stdout
		} else {
			errFile, err = os.Create(importFlags.errorFile)
			defer errFile.(io.Closer).Close()
			if err != nil {
				log.Fatalf("Unable to open error file: %v, cause: %v", importFlags.errorFile, err)
				os.Exit(1)
			}
		}

		// Create Veidemann config client
		client, conn := connection.NewConfigClient()
		defer conn.Close()

		// Create http client
		i.httpClient = importutil.NewHttpClient(importFlags.checkUriTimeout)

		// Create state Database based on seeds in Veidemann
		impf := importutil.NewImportDb(client, importFlags.dbDir, importFlags.resetDb)
		impf.ImportExisting()
		defer impf.Close()

		// Create Record reader for file input
		rr, err := importutil.NewRecordReader(importFlags.filename, &importutil.JsonYamlDecoder{}, "*.json")
		if err != nil {
			log.Fatalf("Unable to create RecordReader: %v", err)
			os.Exit(1)
		}

		// Processor for converting oos records into import records
		proc := func(value interface{}) error {
			sd := value.(*seedDesc)
			if err := i.topLevelUri(sd); err != nil {
				return err
			}
			exists, err := impf.Check(sd.Uri)
			if exists.Code > importutil.DUPLICATE_NEW {
				return fmt.Errorf("seed already exists: %v", sd.Uri)
			}
			if err := i.checkUri(sd); err != nil {
				return err
			}

			exists, err = impf.CheckAndUpdateVeidemann(sd.Uri, sd, func(client configV1.ConfigClient, data interface{}) (id string, err error) {
				obj := data.(*seedDesc)
				e := &configV1.ConfigObject{
					ApiVersion: "v1",
					Kind:       configV1.Kind_crawlEntity,
					Meta: &configV1.Meta{
						Name:        obj.EntityName,
						Description: obj.EntityDescription,
						Label:       obj.EntityLabel,
					},
				}
				ctx := context.Background()
				log.Debugf("store entity: %v", e)
				e, err = client.SaveConfigObject(ctx, e)
				if err != nil {
					//client.DeleteConfigObject(context.Background(), e)
					return "", fmt.Errorf("Error writing crawl entity: %v", err)
				}

				s := &configV1.ConfigObject{
					ApiVersion: "v1",
					Kind:       configV1.Kind_seed,
					Meta: &configV1.Meta{
						Name:        obj.Uri,
						Description: obj.SeedDescription,
						Label:       obj.SeedLabel,
					},
					Spec: &configV1.ConfigObject_Seed{
						Seed: &configV1.Seed{
							EntityRef: &configV1.ConfigRef{
								Kind: configV1.Kind_crawlEntity,
								Id:   e.Id,
							},
							JobRef: obj.crawlJobRef,
						},
					},
				}
				log.Debugf("store seed: %v", s)
				_, err = client.SaveConfigObject(ctx, s)
				if err != nil {
					if d, err := client.DeleteConfigObject(ctx, e); err == nil {
						fmt.Println("Delete entity: ", d)
					} else {
						fmt.Println("Failed deletion of entity: ", err)
					}
					return "", fmt.Errorf("Error writing seed: %v", err)
				}
				return s.Id, nil
			})
			if err != nil {
				return err
			}
			if exists.Code > importutil.NEW {
				return fmt.Errorf("seed already exists: %v", sd.Uri)
			}

			return nil
		}

		// Error handler
		errorHandler := func(state *importutil.StateVal) {
			var uri string
			if state.Val != nil {
				uri = state.Val.(*seedDesc).Uri
			}
			_, _ = fmt.Fprintf(errFile, "{\"uri\": \"%s\", \"err\": \"%s\", \"file\": \"%s\", \"recNum\": %v}\n", uri, state.GetError(), state.GetFilename(), state.GetRecordNum())
		}

		// Create multithreaded executor
		conv := importutil.NewExecutor(32, proc, errorHandler)

		crawlJobRef := []*configV1.ConfigRef{
			{
				Kind: configV1.Kind_crawlJob,
				Id:   importFlags.crawlJobId,
			},
		}

		// Process
		for {
			var sd seedDesc
			state, err := rr.Next(&sd)
			if err == io.EOF {
				break
			}
			if err != nil {
				_, _ = fmt.Fprintf(errFile, "Error decoding record: %v, cause: %v", state, err)
				os.Exit(1)
			}
			if importFlags.crawlJobId != "" {
				sd.crawlJobRef = crawlJobRef
			}

			conv.Do(state, &sd)
		}

		conv.Finish()
		_, _ = fmt.Fprintf(os.Stderr, "\nRecords read: %v, imported: %v, Failed: %v\n", conv.Count, conv.Success, conv.Failed)
	},
}

type importer struct {
	httpClient *http.Client
}

func init() {
	ImportCmd.AddCommand(importSeedCmd)

	importSeedCmd.PersistentFlags().StringVarP(&importFlags.filename, "filename", "f", "", "Filename or directory to read from. "+
		"If input is a directory, all files ending in .yaml or .json will be tried. An input of '-' will read from stdin.")
	importSeedCmd.PersistentFlags().StringVarP(&importFlags.errorFile, "errorfile", "e", "", "File to write errors to.")
	importSeedCmd.PersistentFlags().BoolVarP(&importFlags.toplevel, "toplevel", "", false, "Convert URI to toplevel by removing path.")
	importSeedCmd.PersistentFlags().BoolVarP(&importFlags.checkUri, "checkuri", "", false, "Check the uri for liveness and follow 301")
	importSeedCmd.PersistentFlags().Int64VarP(&importFlags.checkUriTimeout, "checkuri-timeout", "", 500, "Timeout in ms when checking uri for liveness.")
	importSeedCmd.PersistentFlags().StringVarP(&importFlags.crawlJobId, "crawljob-id", "", "", "Set crawlJob ID for new seeds.")
	importSeedCmd.PersistentFlags().StringVarP(&importFlags.dbDir, "db-directory", "b", "/tmp/veidemannctl", "Directory for storing state db")
	importSeedCmd.PersistentFlags().BoolVarP(&importFlags.resetDb, "reset-db", "r", false, "Clean state db")
}

func (i *importer) topLevelUri(s *seedDesc) (err error) {
	uri, err := url.Parse(s.Uri)
	if err != nil {
		return fmt.Errorf("unparseable URL '%v', cause: %v", s.Uri, err)
	}

	if uri.Host == "" {
		return errors.New("unparseable URL")
	}

	if importFlags.toplevel {
		uri.Path = ""
		uri.RawQuery = ""
		uri.Fragment = ""
		s.Uri = uri.Scheme + "://" + uri.Host
	}
	return
}

func (i *importer) checkUri(s *seedDesc) (err error) {
	if importFlags.checkUri {
		i.checkRedirect(s.Uri, s, 0)
	}
	return
}

func (i *importer) checkRedirect(uri string, s *seedDesc, count int) {
	if count > 5 {
		return
	}
	count++

	resp, err := i.httpClient.Head(uri)
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
				i.checkRedirect(uri, s, count)
			}
		} else {
			s.Uri = uri
		}
	}
}
