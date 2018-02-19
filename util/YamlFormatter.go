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

package util

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/golang/protobuf/proto"
	"io"
	"log"
	"reflect"
)

func MarshalYaml(w io.Writer, msg proto.Message) error {
	var values reflect.Value
	values = reflect.ValueOf(msg).Elem().FieldByName("Value")
	if values.IsValid() {
		if values.Len() == 0 {
			fmt.Println("Empty result")
			return nil
		}

		for i := 0; i < values.Len(); i++ {
			if i > 0 {
				w.Write([]byte("---\n"))
			}

			m := values.Index(i).Interface().(proto.Message)

			err := marshalElementYaml(w, m)
			if err != nil {
				return err
			}
		}
	} else {
		m := reflect.ValueOf(msg).Interface().(proto.Message)

		err := marshalElementYaml(w, m)
		if err != nil {
			return err
		}
	}

	return nil
}

func marshalElementYaml(w io.Writer, msg proto.Message) error {
	r, err := EncodeJson(msg)
	if err != nil {
		log.Fatalf("Could not convert %v to JSON: %v", msg, err)
		return err
	}

	final, err := yaml.JSONToYAML(r)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}

	fmt.Fprintln(w, string(final))

	return nil
}
