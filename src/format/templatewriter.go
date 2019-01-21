package format

import (
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"io"
	log "github.com/sirupsen/logrus"
	"sync"
	"text/template"
	tspb "github.com/golang/protobuf/ptypes/timestamp"
)

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
