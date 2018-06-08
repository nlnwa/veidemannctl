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
	"io"
	log "github.com/sirupsen/logrus"
	"os"
	"reflect"
	"strings"
	"io/ioutil"
	"github.com/golang/protobuf/ptypes"
	"text/template"
	tspb "github.com/golang/protobuf/ptypes/timestamp"
)

var jsonMarshaler = jsonpb.Marshaler{EmitDefaults: true}
var jsonUnMarshaler = jsonpb.Unmarshaler{}

type MarshalSpec struct {
	Filename string
	Format   string
	Msg      proto.Message
	Template string
	writer   io.Writer
}

func (s *MarshalSpec) resolve() {
	if s.Filename == "" {
		s.writer = os.Stdout
	} else {
		f, err := os.Create(s.Filename)
		if err != nil {
			log.Fatalf("Could not create file '%s': %v", s.Filename, err)
		}
		defer f.Close()
		s.writer = f
	}

	switch s.Format {
	case "template":
		if s.Template == "" {
			log.Fatal("Format is 'template', but template is missing")
		}
		s.Format = "json"
		s.writer = newTemplateWriter(s.writer, s.Template)
	case "template-file":
		if s.Template == "" {
			log.Fatal("Format is 'template-file', but template is missing")
		}
		data, err := ioutil.ReadFile(s.Template)
		if err != nil {
			log.Fatalf("Template not found: %v", err)
		}
		s.Template = string(data)
		s.Format = "json"
		s.writer = newTemplateWriter(s.writer, s.Template)
	default:
		s.Template = ""
	}
}

type templateWriter struct {
	writer   io.Writer
	template string
	pin      *io.PipeReader
	pout     *io.PipeWriter
}

func newTemplateWriter(writer io.Writer, template string) *templateWriter {
	t := &templateWriter{}
	t.writer = writer
	t.template = template
	t.pin, t.pout = io.Pipe()
	c := make(chan bool)
	go t.unmarshalJson(c)
	return t
}

func (t *templateWriter) Write(p []byte) (n int, err error) {
	return t.pout.Write(p)
}

func (t *templateWriter) unmarshalJson(c chan bool) {
	dec := json.NewDecoder(t.pin)
	for dec.More() {
		var val interface{}
		err := dec.Decode(&val)
		if err != nil {
			log.Fatal(err)
		}
		t.applyTemplate(val)
	}
}

func (t *templateWriter) applyTemplate(val interface{}) {
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
	}

	tmpl, err := template.New("Template").Funcs(funcMap).Parse(t.template)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(t.writer, val)
	if err != nil {
		panic(err)
	}
}

func Marshal(spec *MarshalSpec) error {
	spec.resolve()

	switch spec.Format {
	case "json":
		err := MarshalJson(spec.writer, spec.Msg)
		if err != nil {
			return err
		}
	case "yaml":
		err := MarshalYaml(spec.writer, spec.Msg)
		if err != nil {
			return err
		}
	case "table":
		err := MarshalTable(spec.writer, spec.Msg)
		if err != nil {
			return err
		}
	default:
		log.Fatalf("Illegal format %s", spec.Format)
	}
	return nil
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

		var val interface{}
		err := yaml.Unmarshal(data, &val)
		if err != nil {
			log.Fatal(err)
		}

		v := val.(map[string]interface{})
		k := v["kind"]
		if k == nil {
			return nil, fmt.Errorf("Missing kind")
		}
		kind := k.(string)
		delete(v, "kind")

		b, _ := json.Marshal(&v)

		buf := bytes.NewBuffer(b)
		t := GetObjectType(kind)
		if t == nil {
			return nil, fmt.Errorf("Unknown kind '%v'", kind)
		}

		target := reflect.New(t.Elem()).Interface().(proto.Message)

		jsonUnMarshaler.Unmarshal(buf, target)
		result = append(result, target)
	}
	return result, nil
}

func UnmarshalJson(r io.Reader, result []proto.Message) ([]proto.Message, error) {
	dec := json.NewDecoder(r)
	for dec.More() {
		var val interface{}
		err := dec.Decode(&val)
		if err != nil {
			log.Fatal(err)
		}
		v := val.(map[string]interface{})
		kind := v["kind"].(string)
		delete(v, "kind")

		b, _ := json.Marshal(&v)

		buf := bytes.NewBuffer(b)
		t := GetObjectType(kind)

		target := reflect.New(t.Elem()).Interface().(proto.Message)

		jsonUnMarshaler.Unmarshal(buf, target)
		result = append(result, target)
	}

	return result, nil
}
