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
	"fmt"
	tm "github.com/buger/goterm"
	"github.com/nlnwa/veidemann-api-go/config/v1"
	"reflect"
	"strings"
)

var debugEnabled = false

type tableFormatter struct {
	*MarshalSpec
	table *tm.Table
}

func newTableFormatter(s *MarshalSpec) Formatter {
	return &tableFormatter{
		MarshalSpec: s,
		table:       tm.NewTable(2, 10, 2, ' ', 0),
	}
}

func (tf *tableFormatter) WriteHeader() error {
	var header string
	for idx, tab := range GetTableDefForKind(tf.Kind) {
		if idx > 0 {
			header += "\t"
		}
		header += tab
	}

	fmt.Fprintf(tf.table, "\n%s\n", header)
	return nil
}

func (tf *tableFormatter) WriteRecord(record interface{}) error {
	m := record.(*config.ConfigObject)

	tableDef := GetTableDefForKind(m.Kind)

	err := formatData(tf.table, tableDef, m)
	if err != nil {
		return err
	}
	return nil
}

func (tf *tableFormatter) Close() error {
	tm.Println(tf.table)
	tm.Flush()

	return tf.MarshalSpec.Close()
}

func formatData(t *tm.Table, tableDef []string, msg *config.ConfigObject) error {
	var line string
	for idx, tab := range tableDef {
		if idx > 0 {
			line += "\t"
		}
		line += fmt.Sprint(getField(msg, tab))
	}

	fmt.Fprintf(t, "%s\n", line)
	return nil
}

func getField(msg *config.ConfigObject, fieldName string) reflect.Value {
	tokens := strings.Split(fieldName, ".")
	v := reflect.ValueOf(msg)
	for _, tok := range tokens {
		v = reflect.Indirect(v)
		if v.Kind() == reflect.Interface {
			v = reflect.Indirect(reflect.ValueOf(v.Interface()))
		}
		v = v.FieldByName(tok)
		if v.Kind() == reflect.Interface && v.IsNil() {
			return v
		}
	}
	return v
}
