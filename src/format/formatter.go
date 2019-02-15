// Copyright © 2017 National Library of Norway.
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

package format

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/golang/protobuf/jsonpb"
	"github.com/nlnwa/veidemann-api-go/config/v1"
	"github.com/nlnwa/veidemannctl/bindata"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

var jsonMarshaler = jsonpb.Marshaler{EmitDefaults: true}
var jsonUnMarshaler = jsonpb.Unmarshaler{}

type Formatter interface {
	WriteHeader() error
	WriteRecord(record interface{}) error
	Close() error
}

type MarshalSpec struct {
	Kind                config.Kind
	Filename            string
	Format              string
	Template            string
	HeaderTemplate      string
	defaultTemplateName string
	Writer              io.Writer

	rFilename string
	rFormat   string
	rTemplate string
	rWriter   io.Writer
	resolved  bool
	wg        sync.WaitGroup
}

func NewFormatter(kind config.Kind, filename string, format string, template string, headerTemplate string) Formatter {
	s := &MarshalSpec{
		Kind:           kind,
		Filename:       filename,
		Format:         format,
		Template:       template,
		HeaderTemplate: headerTemplate,
	}

	s.resolve()

	var formatter Formatter

	switch s.rFormat {
	case "json":
		formatter = newJsonFormatter(s)
	case "yaml":
		formatter = newYamlFormatter(s)
	case "table":
		formatter = newTableFormatter(s)
	default:
		log.Fatalf("Illegal format %s", s.rFormat)
	}
	return formatter
}

func (s *MarshalSpec) WriteHeaderForKind(kind config.Kind) error {
	fmt.Println("HEAD")
	return nil
}

func (s *MarshalSpec) resolve() {
	if !s.resolved {
		if s.Writer == nil {
			if s.Filename == "" {
				s.rWriter = os.Stdout
			} else {
				f, err := os.Create(s.Filename)
				if err != nil {
					log.Fatalf("Could not create file '%s': %v", s.Filename, err)
				}
				s.rWriter = f
			}
		} else {
			s.rWriter = s.Writer
		}

		switch s.Format {
		case "template":
			if s.Template == "" {
				if s.defaultTemplateName != "" {
					data, err := bindata.Asset(s.defaultTemplateName)
					if err != nil {
						panic(err)
					}
					s.rTemplate = string(data)
				} else {
					log.Fatal("Format is 'template', but template is missing")
				}
			} else {
				s.rTemplate = s.Template
			}
			s.rFormat = "json"
			s.rWriter = newTemplateWriter(s.rWriter, s.rTemplate, s.HeaderTemplate, &s.wg)
		case "template-file":
			if s.Template == "" {
				log.Fatal("Format is 'template-file', but template is missing")
			}
			data, err := ioutil.ReadFile(s.Template)
			if err != nil {
				log.Fatalf("Template not found: %v", err)
			}
			s.rTemplate = string(data)
			s.rFormat = "json"
			s.rWriter = newTemplateWriter(s.rWriter, s.rTemplate, s.HeaderTemplate, &s.wg)
		case "json":
			s.rFormat = s.Format
			data, err := bindata.Asset("json.template")
			if err != nil {
				panic(err)
			}
			s.rTemplate = string(data)
			s.HeaderTemplate = ""
			s.rWriter = newTemplateWriter(s.rWriter, s.rTemplate, s.HeaderTemplate, &s.wg)
		case "yaml":
			s.rTemplate = ""
			s.rFormat = s.Format
		default:
			s.rTemplate = s.Template
			s.rFormat = s.Format
		}
	}
	s.resolved = true
}

func (s *MarshalSpec) Close() error {
	if nil != s {
		if c, ok := s.rWriter.(io.Closer); ok {
			if err := c.Close(); err != nil {
				return err
			}
		}
	}
	s.wg.Wait()
	return nil
}

func Unmarshal(filename string) ([]*config.ConfigObject, error) {
	result := make([]*config.ConfigObject, 0, 16)
	if filename == "" {
		return UnmarshalYaml(os.Stdin, result)
	} else {
		f, err := os.Open(filename)
		if err != nil {
			log.Fatalf("Could not open file '%s': %v", filename, err)
			return nil, err
		}
		defer f.Close()
		if fi, _ := f.Stat(); fi.IsDir() {
			fis, _ := f.Readdir(0)
			for _, fi = range fis {
				if !fi.IsDir() && HasSuffix(fi.Name(), ".yaml", ".yml", ".json") {
					fmt.Println("Reading file: ", fi.Name())
					f, err = os.Open(fi.Name())
					if err != nil {
						log.Fatalf("Could not open file '%s': %v", filename, err)
						return nil, err
					}
					defer f.Close()

					if HasSuffix(f.Name(), ".yaml", ".yml") {
						result, err = UnmarshalYaml(f, result)
					} else {
						result, err = UnmarshalJson(f, result)
					}
					if err != nil {
						return nil, err
					}
				}
			}
			return result, nil
		} else {
			if HasSuffix(f.Name(), ".yaml", ".yml") {
				return UnmarshalYaml(f, result)
			} else {
				return UnmarshalJson(f, result)
			}
		}
	}
	return result, nil
}

func ReadYamlDocument(r *bufio.Reader) ([]byte, error) {
	delim := []byte{'-', '-', '-'}
	var (
		inDoc  bool  = true
		err    error = nil
		l, doc []byte
	)
	for inDoc && err == nil {
		isPrefix := true
		ln := []byte{}
		for isPrefix && err == nil {
			l, isPrefix, err = r.ReadLine()
			ln = append(ln, l...)
		}

		if len(ln) >= 3 && bytes.Equal(delim, ln[:3]) {
			inDoc = false
		} else {
			doc = append(doc, ln...)
			doc = append(doc, '\n')
		}
	}
	return doc, err
}

func UnmarshalYaml(r io.Reader, result []*config.ConfigObject) ([]*config.ConfigObject, error) {
	br := bufio.NewReader(r)

	var (
		data    []byte
		readErr error = nil
	)
	for readErr == nil {
		data, readErr = ReadYamlDocument(br)
		if readErr != nil && readErr != io.EOF {
			return nil, readErr
		}

		b, err := yaml.YAMLToJSON(data)
		if err != nil {
			log.Fatalf("Error decoding: %v, %v", err, data)
		}
		buf := bytes.NewBuffer(b)

		target := &config.ConfigObject{}
		jsonUnMarshaler.Unmarshal(buf, target)
		result = append(result, target)
	}
	return result, nil
}

func UnmarshalJson(r io.Reader, result []*config.ConfigObject) ([]*config.ConfigObject, error) {
	dec := json.NewDecoder(r)
	for dec.More() {
		target := &config.ConfigObject{}
		jsonUnMarshaler.UnmarshalNext(dec, target)
		result = append(result, target)
	}

	return result, nil
}

func HasSuffix(s string, suffix... string) bool {
	for _, suf := range suffix {
		if strings.HasSuffix(s, suf) {
			return true
		}
	}
	return false
}
