/*
 * Copyright 2019 National Library of Norway.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *       http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package format

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	tspb "github.com/golang/protobuf/ptypes/timestamp"
	log "github.com/sirupsen/logrus"
	"text/template"
)

type templateFormatter struct {
	*MarshalSpec
	headerWritten  bool
	parsedTemplate *template.Template
}

func newTemplateFormatter(s *MarshalSpec) (Formatter, error) {
	t := &templateFormatter{
		MarshalSpec: s,
	}
	pt, err := parseTemplate(t.rTemplate)
	if err != nil {
		return nil, err
	}
	t.parsedTemplate = pt
	return t, nil
}

func (tf *templateFormatter) WriteRecord(record interface{}) error {
	if !tf.headerWritten {
		tf.headerWritten = true
		tpl := tf.parsedTemplate.Lookup("HEADER")
		if tpl != nil {
			err := tpl.Execute(tf.rWriter, nil)
			if err != nil {
				log.Fatal("Failed applying header template: ", err)
			}
		}
	}

	tpl := tf.parsedTemplate
	if tpl != nil {
		if r, ok := record.(string); ok {
			var j interface{}
			err := json.Unmarshal([]byte(r), &j)
			if err != nil {
				return fmt.Errorf("failed to parse json: %v", err)
			}
			record = j
		}
		err := tpl.Execute(tf.rWriter, record)
		if err != nil {
			return fmt.Errorf("failed applying template to '%v': %v", record, err)
		}
	}
	return nil
}

func parseTemplate(templateString string) (*template.Template, error) {
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
		"json": func(v interface{}) string {
			if v == nil {
				return ""
			} else {
				var buf bytes.Buffer
				encoder := json.NewEncoder(&buf)
				encoder.SetEscapeHTML(false)
				err := encoder.Encode(v)
				if err != nil {
					log.Fatal(err)
				}
				return string(buf.Bytes())
			}
		},
		"prettyJson": func(v interface{}) string {
			if v == nil {
				return ""
			} else {
				var buf bytes.Buffer
				encoder := json.NewEncoder(&buf)
				encoder.SetEscapeHTML(false)
				encoder.SetIndent("", "  ")
				err := encoder.Encode(v)
				if err != nil {
					log.Fatal(err)
				}
				return buf.String()
			}
		},
		"nl": func() string { return "\n" },
	}

	return template.New("Template").Funcs(funcMap).Parse(templateString)
}
