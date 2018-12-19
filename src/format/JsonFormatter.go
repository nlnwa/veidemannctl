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
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"io"
	"log"
	"reflect"
)

type jsonFormatter struct {
	*MarshalSpec
}

func newJsonFormatter(s *MarshalSpec) Formatter {
	return &jsonFormatter{
		MarshalSpec: s,
	}
}

func (jf *jsonFormatter) WriteHeader() error {
	t, ok := jf.rWriter.(*templateWriter)
	if ok && jf.HeaderTemplate != "" {
		t.applyHeaderTemplate(jf.HeaderTemplate)
	}
	return nil
}

func (jf *jsonFormatter) WriteRecord(record interface{}) error {
	switch v := record.(type) {
	case string:
		fmt.Fprint(jf.rWriter, v)
	case proto.Message:
		var values reflect.Value
		values = reflect.ValueOf(v).Elem().FieldByName("Value")
		if values.IsValid() {
			if values.Len() == 0 {
				fmt.Println("Empty result")
				return nil
			}

			for i := 0; i < values.Len(); i++ {
				m := values.Index(i).Interface().(proto.Message)

				err := marshalElementJson(jf.rWriter, m)
				if err != nil {
					return err
				}
			}
		} else {
			m := reflect.ValueOf(v).Interface().(proto.Message)

			err := marshalElementJson(jf.rWriter, m)
			if err != nil {
				return err
			}
		}
	default:
		log.Fatalf("Illegal format %s", jf.rFormat)
	}
	return nil
}

func marshalElementJson(w io.Writer, msg proto.Message) error {
	r, err := EncodeJson(msg)
	if err != nil {
		log.Fatalf("Could not convert %v to JSON: %v", msg, err)
		return err
	}

	var out bytes.Buffer
	json.Indent(&out, r, "", "  ")

	fmt.Fprintln(w, out.String())
	return nil
}

func EncodeJson(msg proto.Message) ([]byte, error) {
	var buf bytes.Buffer
	if err := jsonMarshaler.Marshal(&buf, msg); err != nil {
		log.Fatalf("Could not convert %v to YAML: %v", msg, err)
		return nil, err
	}

	values := make(map[string]interface{})
	json.Unmarshal(buf.Bytes(), &values)

	b, _ := json.Marshal(values)

	return b, nil
}
