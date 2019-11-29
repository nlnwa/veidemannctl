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

package format

import (
	"errors"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/golang/protobuf/proto"
	"io"
	"reflect"
)

type yamlFormatter struct {
	*MarshalSpec
}

func newYamlFormatter(s *MarshalSpec) Formatter {
	return &yamlFormatter{
		MarshalSpec: s,
	}
}

func (yf *yamlFormatter) WriteRecord(record interface{}) error {
	switch v := record.(type) {
	case string:
		final, err := yaml.JSONToYAML([]byte(v))
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return err
		}

		fmt.Fprint(yf.rWriter, string(final))
		fmt.Fprintln(yf.rWriter, "---")
	case proto.Message:
		var values reflect.Value
		values = reflect.ValueOf(v).Elem().FieldByName("Value")
		if values.IsValid() {
			if values.Len() == 0 {
				fmt.Println("Empty result")
				return nil
			}

			for i := 0; i < values.Len(); i++ {
				if i > 0 {
					yf.rWriter.Write([]byte("---\n"))
				}

				m := values.Index(i).Interface().(proto.Message)

				err := marshalElementYaml(yf.rWriter, m)
				if err != nil {
					return err
				}
			}
		} else {
			m := reflect.ValueOf(v).Interface().(proto.Message)

			err := marshalElementYaml(yf.rWriter, m)
			if err != nil {
				return err
			}
			yf.rWriter.Write([]byte("---\n"))
		}
	default:
		return errors.New(fmt.Sprintf("Illegal record type '%T'", record))
	}
	return nil
}

func marshalElementYaml(w io.Writer, msg proto.Message) error {
	r, err := EncodeJson(msg)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not convert %v to JSON: %v", msg, err))
	}

	final, err := yaml.JSONToYAML(r)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not convert %v to YAML: %v", r, err))
	}

	fmt.Fprintln(w, string(final))

	return nil
}
