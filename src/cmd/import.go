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
	"encoding/json"
	"errors"
	"fmt"
	configV1 "github.com/nlnwa/veidemann-api-go/config/v1"
	"github.com/nlnwa/veidemannctl/src/connection"
	"github.com/nlnwa/veidemannctl/src/format"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"gopkg.in/yaml.v2"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type seed struct {
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
}

var importFlags struct {
	errorfile       string
	toplevel        bool
	checkUri        bool
	checkUriTimeout int64
}

var httpClient *http.Client

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import seeds",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		httpTimeout := time.Duration(importFlags.checkUriTimeout) * time.Millisecond
		httpClient = &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
			Transport: &http.Transport{
				DisableKeepAlives:     true,
				ResponseHeaderTimeout: httpTimeout,
			},
			Timeout: httpTimeout,
		}

		if filename == "" {
			cmd.Usage()
			os.Exit(1)
		} else if filename == "-" {
			filename = ""
		}

		dr, err := NewDocReader(filename)
		if err != nil {
			log.Fatalf("Parse error: %v", err)
			os.Exit(1)
		}

		var out io.Writer
		if importFlags.errorfile == "" {
			out = os.Stdout
		} else {
			out, err = os.Create(importFlags.errorfile)
			defer out.(io.Closer).Close()
			if err != nil {
				log.Fatalf("Unable to open file: %v, cause: %v", importFlags.errorfile, err)
				os.Exit(1)
			}
		}

		client, conn := connection.NewConfigClient()
		defer conn.Close()

		var (
			count   int
			success int
			failed  int
		)

		subCount := 8
		cData := make(chan *seed)
		cComplete := make(chan subResponse)
		for i := 0; i < subCount; i++ {
			go storeRecord(client, out, cData, cComplete)
		}

		for {
			var s seed
			obj := &s
			err = dr.Next(obj)
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				if err == io.ErrUnexpectedEOF {
					obj.err = err
					count++
					failed++
					printErr(out, obj)
				}
				break
			}

			// Print progress
			count++
			fmt.Fprint(os.Stderr, ".")
			if count%100 == 0 {
				fmt.Fprintln(os.Stderr, count)
			}

			if err != nil {
				obj.err = err
				failed++
				printErr(out, obj)
				continue
			}

			if obj != nil {
				cData <- obj
			}
		}

		close(cData)
		subCompleted := 0
		for i := range cComplete {
			success += i.success
			failed += i.failed
			subCompleted++
			if subCompleted >= subCount {
				break
			}
		}
		_, _ = fmt.Fprintf(os.Stderr, "\nRecords read: %v, imported: %v, failed: %v\n", count, success, failed)

	},
}

type subResponse struct {
	count   int
	success int
	failed  int
}

func storeRecord(client configV1.ConfigClient, out io.Writer, cin chan *seed, cout chan subResponse) {
	res := subResponse{}

	for obj := range cin {
		res.count++
		err := checkSeed(obj, client)
		if err != nil {
			obj.err = err
			res.failed++
			printErr(out, obj)
		} else {
			e := &configV1.ConfigObject{
				ApiVersion: "v1",
				Kind:       configV1.Kind_crawlEntity,
				Meta: &configV1.Meta{
					Name:        obj.EntityName,
					Description: obj.EntityDescription,
					Label:       obj.EntityLabel,
				},
			}
			log.Debugf("store entity: %v", e)
			e, err = client.SaveConfigObject(context.Background(), e)
			if err != nil {
				obj.err = fmt.Errorf("Error writing crawl entity: %v", err)
				res.failed++
				printErr(out, obj)
				client.DeleteConfigObject(context.Background(), e)
				continue
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
					},
				},
			}
			log.Debugf("store seed: %v", s)
			_, err = client.SaveConfigObject(context.Background(), s)
			if err != nil {
				obj.err = fmt.Errorf("Error writing seed: %v", err)
				res.failed++
				printErr(out, obj)
				continue
			}

			res.success++
		}
	}
	cout <- res
}

func printErr(out io.Writer, obj *seed) {
	var uri string
	if obj != nil {
		uri = obj.Uri
	}
	_, _ = fmt.Fprintf(out, "{\"uri\": \"%s\", \"err\": \"%s\", \"file\": \"%s\", \"recNum\": %v}\n", uri, obj.err, obj.fileName, obj.recNum)
}

func init() {
	RootCmd.AddCommand(importCmd)

	importCmd.PersistentFlags().StringVarP(&filename, "filename", "f", "", "Filename or directory to read from. "+
		"If input is a directory, all files ending in .yaml or .json will be tried. An input of '-' will read from stdin.")
	importCmd.PersistentFlags().StringVarP(&importFlags.errorfile, "errorfile", "e", "", "File to write errors to.")
	importCmd.PersistentFlags().BoolVarP(&importFlags.toplevel, "toplevel", "", false, "Convert URI to toplevel by removing path.")
	importCmd.PersistentFlags().BoolVarP(&importFlags.checkUri, "checkuri", "", false, "Check the uri for liveness and follow 301")
	importCmd.PersistentFlags().Int64VarP(&importFlags.checkUriTimeout, "checkuri-timeout", "", 500, "Timeout in ms when checking uri for liveness.")
}

func checkSeed(s *seed, client configV1.ConfigClient) (err error) {
	uri, err := checkUri(s, true)
	if err != nil {
		return
	}

	req := &configV1.ListRequest{
		Kind:      configV1.Kind_seed,
		NameRegex: uri.Host,
	}

	r, err := client.ListConfigObjects(context.Background(), req)
	if err != nil {
		return
	}

	o, err := r.Recv()
	if err != nil && err != io.EOF {
		return
	}
	if o != nil {
		return fmt.Errorf("seed already exists: %v", o.Meta.Name)
	}

	if s.SeedDescription == "" && s.Description != "" {
		s.SeedDescription = s.Description
	}

	if s.EntityDescription == "" && s.Description != "" {
		s.EntityDescription = s.Description
	}

	return nil
}

func checkUri(s *seed, followRedirect bool) (uri *url.URL, err error) {
	uri, err = url.Parse(s.Uri)
	if err != nil {
		return uri, fmt.Errorf("unparseable URL '%v', cause: %v", s.Uri, err)
	}

	if uri.Host == "" {
		return uri, errors.New("unparseable URL")
	}

	if s.EntityName == "" {
		return uri, errors.New("entityName cannot be empty")
	}

	if importFlags.toplevel {
		uri.Path = ""
		uri.RawQuery = ""
		uri.Fragment = ""
		s.Uri = uri.Scheme + "://" + uri.Host
	}

	if importFlags.checkUri {
		var resp *http.Response
		resp, err = httpClient.Head(s.Uri)
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
				default:
					return
				}
			}
		} else {
			resp.Body.Close()
			if resp.StatusCode == 301 && followRedirect {
				loc := resp.Header.Get("Location")
				if loc != "" {
					s.Uri = loc
					return checkUri(s, false)
				}
			}
		}
	}
	return
}

type docReader struct {
	decoder     yjDecoder
	curFile     *os.File
	dir         *os.File
	curFileName string
	curRecNum   int
}

type yjDecoder interface {
	Decode(v interface{}) (err error)
}

func (d *docReader) Next(target *seed) (err error) {
	err = d.decoder.Decode(target)
	if err == io.EOF {
		if err = d.nextDecoder(); err != nil {
			return
		}
		return d.Next(target)
	}
	d.curRecNum++
	target.fileName = d.curFileName
	target.recNum = d.curRecNum
	target.err = err
	return
}

func (d *docReader) nextDecoder() (err error) {
	if d.curFile != nil {
		d.curFile.Close()
	}

	if d.dir == nil {
		return io.EOF
	}

	var f []os.FileInfo
	f, err = d.dir.Readdir(1)
	if err != nil {
		return
	}
	fi := f[0]

	if !fi.IsDir() && format.HasSuffix(fi.Name(), ".yaml", ".yml", ".json") {
		fmt.Println("Reading file: ", fi.Name())
		d.curFile, err = os.Open(filepath.Join(d.dir.Name(), fi.Name()))
		if err != nil {
			log.Fatalf("Could not open file '%s': %v", filename, err)
			return
		}
		d.curRecNum = 0
		d.curFileName = d.curFile.Name()

		if strings.HasSuffix(d.curFile.Name(), ".yaml") || strings.HasSuffix(d.curFile.Name(), ".yml") {
			dec := yaml.NewDecoder(d.curFile)
			dec.SetStrict(true)
			d.decoder = dec
			return
		} else {
			dec := json.NewDecoder(d.curFile)
			dec.DisallowUnknownFields()
			d.decoder = dec
			return
		}
	}
	return
}

func NewDocReader(filename string) (d *docReader, err error) {
	d = &docReader{}

	if filename == "" {
		d.decoder = yaml.NewDecoder(os.Stdin)
		return
	} else {
		var f *os.File
		f, err = os.Open(filename)
		if err != nil {
			//log.Fatalf("Could not open file '%s': %v", filename, err)
			return
		}
		if fi, _ := f.Stat(); fi.IsDir() {
			d.dir = f
			d.nextDecoder()
			return
		} else {
			d.curFileName = f.Name()
			if format.HasSuffix(f.Name(), ".yaml", ".yml") {
				dec := yaml.NewDecoder(f)
				dec.SetStrict(true)
				d.decoder = dec
				return
			} else {
				dec := json.NewDecoder(f)
				dec.DisallowUnknownFields()
				d.decoder = dec
				return
			}
		}
	}
}
