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
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	tspb "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/nlnwa/veidemann-api-go/config/v1"
	"github.com/nlnwa/veidemannctl/bindata"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"text/template"
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
				defer f.Close()
				s.rWriter = f
			}
		} else {
			s.rWriter = s.Writer
		}

		switch s.Format {
		case "template":
			if s.Template == "" {
				if (s.defaultTemplateName != "") {
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
		if t, ok := s.rWriter.(io.Closer); ok {
			if err := t.Close(); err != nil {
				return err
			}
		}
	}
	s.wg.Wait()
	return nil
}

type templateWriter struct {
	writer         io.Writer
	template       string
	headerTemplate string
	pin            *io.PipeReader
	pout           *io.PipeWriter
	wg             *sync.WaitGroup
}

func newTemplateWriter(writer io.Writer, template string, headerTemplate string, wg *sync.WaitGroup) *templateWriter {
	t := &templateWriter{}
	t.writer = writer
	t.template = template
	t.headerTemplate = headerTemplate
	t.pin, t.pout = io.Pipe()
	t.wg = wg
	t.wg.Add(1)
	go t.unmarshalJson()
	return t
}

func (t *templateWriter) Write(p []byte) (n int, err error) {
	return t.pout.Write(p)
}

func (t *templateWriter) Close() error {
	return t.pout.Close()
}

func (t *templateWriter) unmarshalJson() {
	defer t.wg.Done()
	dec := json.NewDecoder(t.pin)
	for dec.More() {
		var val interface{}
		err := dec.Decode(&val)
		if err != nil {
			log.Fatal("Failed decoding json: ", err)
		}
		//t.applyTemplate(val, "json.template")
		t.applyRecordTemplate(val)
	}
	if c, ok := t.writer.(io.Closer); ok {
		c.Close()
	}
}

func (t *templateWriter) applyHeaderTemplate(msg interface{}) {
	t.applyTemplate(nil, t.headerTemplate)
}

func (t *templateWriter) applyRecordTemplate(val interface{}) {
	t.applyTemplate(val, t.template)
}

func (t *templateWriter) applyTemplate(val interface{}, templateString string) {
	ESC := string(0x1b)
	funcMap := template.FuncMap{
		"reset":         func() string { return ESC + "[0m" },
		"bold":          func() string { return ESC + "[1m" },
		"inverse":       func() string { return ESC + "[7m" },
		"red":           func() string { return ESC + "[31m" },
		"green":         func() string { return ESC + "[32m" },
		"yellow":        func() string { return ESC + "[33m" },
		"blue":          func() string { return ESC + "[34m" },
		"magenta":       func() string { return ESC + "[35m" },
		"cyan":          func() string { return ESC + "[36m" },
		"brightred":     func() string { return ESC + "[1;31m" },
		"brightgreen":   func() string { return ESC + "[1;32m" },
		"brightyellow":  func() string { return ESC + "[1;33m" },
		"brightblue":    func() string { return ESC + "[1;34m" },
		"brightmagenta": func() string { return ESC + "[1;35m" },
		"brightcyan":    func() string { return ESC + "[1;36m" },
		"bgwhite":       func() string { return ESC + "[47m" },
		"bgbrightblack": func() string { return ESC + "[100m" },
		"time": func(ts *tspb.Timestamp) string {
			if ts == nil {
				return "                        "
			} else {
				return fmt.Sprintf("%-24.24s", ptypes.TimestampString(ts))
			}
		},
		"rethinktime": func(ts map[string]interface{}) string {
			if ts == nil {
				return "                        "
			} else {
				dateTime := ts["dateTime"].(map[string]interface{})
				date := dateTime["date"].(map[string]interface{})
				time := dateTime["time"].(map[string]interface{})
				return fmt.Sprintf("%04.f-%02.f-%02.fT%02.f:%02.f:%02.f", date["year"], date["month"], date["day"], time["hour"], time["minute"], time["second"])
			}
		},
		"prettyJson": func(v interface{}) string {
			if v == nil {
				return ""
			} else {
				json, err := json.MarshalIndent(v, "", "  ")
				if err != nil {
					log.Fatal(err)
				}
				return string(json)
			}
		},
		"nl": func() string { return "\n" },
	}

	tmpl, err := template.New("Template").Funcs(funcMap).Parse(templateString)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(t.writer, val)
	if err != nil {
		panic(err)
	}
}

func Unmarshal(filename string) ([]proto.Message, error) {
	result := make([]proto.Message, 0, 16)
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
				if !fi.IsDir() && (strings.HasSuffix(fi.Name(), ".yaml") || strings.HasSuffix(fi.Name(), ".json")) {
					fmt.Println("Reading file: ", fi.Name())
					f, err = os.Open(fi.Name())
					if err != nil {
						log.Fatalf("Could not open file '%s': %v", filename, err)
						return nil, err
					}
					defer f.Close()

					if strings.HasSuffix(f.Name(), ".yaml") {
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
			if strings.HasSuffix(f.Name(), ".yaml") {
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

func UnmarshalYaml(r io.Reader, result []proto.Message) ([]proto.Message, error) {
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

func UnmarshalJson(r io.Reader, result []proto.Message) ([]proto.Message, error) {
	dec := json.NewDecoder(r)
	for dec.More() {
		target := &config.ConfigObject{}
		jsonUnMarshaler.UnmarshalNext(dec, target)
		result = append(result, target)
	}

	return result, nil
}
