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
	"errors"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/golang/protobuf/jsonpb"
	"github.com/nlnwa/veidemann-api-go/config/v1"
	"github.com/nlnwa/veidemannctl/bindata"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

var jsonMarshaler = jsonpb.Marshaler{EmitDefaults: true}
var jsonUnMarshaler = jsonpb.Unmarshaler{}

type Formatter interface {
	WriteRecord(record interface{}) error
	Close() error
}

type MarshalSpec struct {
	ObjectType          string
	Format              string
	Template            string
	defaultTemplateName string

	rFormat   string
	rTemplate string
	rWriter   io.Writer
	resolved  bool
}

func NewFormatter(objectType string, out io.Writer, format string, template string) (formatter Formatter, err error) {
	if out == nil {
		return nil, errors.New("missing writer")
	}

	s := &MarshalSpec{
		ObjectType: objectType,
		Format:     format,
		Template:   template,
		rWriter:    out,
	}

	if err := s.resolve(); err != nil {
		return nil, err
	}

	switch s.rFormat {
	case "json":
		formatter = newJsonFormatter(s)
	case "yaml":
		formatter = newYamlFormatter(s)
	case "template":
		formatter, err = newTemplateFormatter(s)
	case "table":
		formatter, err = newTemplateFormatter(s)
	case "wide":
		formatter, err = newTemplateFormatter(s)
	default:
		return nil, errors.New(fmt.Sprintf("Illegal or missing format '%s'", s.rFormat))
	}
	return
}

// ResolveWriter creates a file and returns a io.Writer for the file.
// If filename is empty, os.StdOut is returned
func ResolveWriter(filename string) (io.Writer, error) {
	if filename == "" {
		return os.Stdout, nil
	} else {
		f, err := os.Create(filename)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Could not create file '%s': %v", filename, err))
		}
		return f, nil
	}
}

func (s *MarshalSpec) resolve() (err error) {
	if !s.resolved {
		switch s.Format {
		case "template":
			if s.Template == "" {
				if s.defaultTemplateName != "" {
					data, err := bindata.Asset(s.defaultTemplateName)
					if err != nil {
						return err
					}
					s.rTemplate = string(data)
				} else {
					return errors.New("Format is 'template', but template is missing")
				}
			} else {
				s.rTemplate = s.Template
			}
			s.rFormat = s.Format
		case "template-file":
			if s.Template == "" {
				return errors.New("format is 'template-file', but template is missing")
			}
			data, err := ioutil.ReadFile(s.Template)
			if err != nil {
				return fmt.Errorf("template not found: %v", err)
			}
			s.rTemplate = string(data)
			s.rFormat = s.Format
		case "json":
			s.rTemplate = ""
			s.rFormat = s.Format
		case "yaml":
			s.rTemplate = ""
			s.rFormat = s.Format
		case "table":
			if s.ObjectType == "" {
				return fmt.Errorf("format is table, but object type is missing")
			}

			templateName := s.ObjectType + "_table.template"
			data, err := bindata.Asset(templateName)
			if err != nil {
				return err
			}
			s.rTemplate = string(data)
			s.rFormat = s.Format
		case "wide":
			if s.ObjectType == "" {
				return fmt.Errorf("format is wide, but object type is missing")
			}

			templateName := s.ObjectType + "_wide.template"
			data, err := bindata.Asset(templateName)
			if err != nil {
				return err
			}
			s.rTemplate = string(data)
			s.rFormat = s.Format
		default:
			s.rTemplate = s.Template
			s.rFormat = s.Format
		}
	}
	s.resolved = true
	return nil
}

func (s *MarshalSpec) Close() error {
	if nil != s {
		if c, ok := s.rWriter.(io.Closer); ok && c != os.Stdout {
			if err := c.Close(); err != nil {
				return err
			}
		}
	}
	return nil
}

func Unmarshal(filename string) ([]*config.ConfigObject, error) {
	result := make([]*config.ConfigObject, 0, 16)
	if filename == "" {
		return UnmarshalYaml(os.Stdin, result)
	} else {
		f, err := os.Open(filename)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Could not open file '%s': %v", filename, err))
		}
		defer f.Close()
		if fi, _ := f.Stat(); fi.IsDir() {
			fis, _ := f.Readdir(0)
			for _, fi = range fis {
				if !fi.IsDir() && HasSuffix(fi.Name(), ".yaml", ".yml", ".json") {
					fmt.Println("Reading file: ", fi.Name())
					f, err = os.Open(fi.Name())
					if err != nil {
						return nil, errors.New(fmt.Sprintf("Could not open file '%s': %v", filename, err))
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
			return nil, errors.New(fmt.Sprintf("Error decoding: %v, %v", err, data))
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

func HasSuffix(s string, suffix ...string) bool {
	for _, suf := range suffix {
		if strings.HasSuffix(s, suf) {
			return true
		}
	}
	return false
}
