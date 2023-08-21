package report

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_queryCmdOptions_parseQuery(t *testing.T) {
	tests := []struct {
		name             string
		fieldsFromFlags  *queryCmdOptions
		fieldsAfterParse *queryCmdOptions
		want             *query
		wantErr          assert.ErrorAssertionFunc
	}{
		{"template file",
			&queryCmdOptions{queryOrFile: "testdata/template1.yaml"},
			&queryCmdOptions{
				queryOrFile: "testdata/template1.yaml",
			},
			&query{
				Name:        "template1",
				Description: "Example template\n",
				Query:       "r.db('veidemann').table('config_crawl_entities')\n",
				Template:    "{{.id}} {{.meta.name}}\n",
			},
			assert.NoError},
		{"template file with template flag",
			&queryCmdOptions{
				queryOrFile: "testdata/template1.yaml",
				goTemplate:  "{{.id}}",
			},
			&queryCmdOptions{
				queryOrFile: "testdata/template1.yaml",
				goTemplate:  "{{.id}}",
			},
			&query{
				Name:        "template1",
				Description: "Example template\n",
				Query:       "r.db('veidemann').table('config_crawl_entities')\n",
				Template:    "{{.id}} {{.meta.name}}\n",
			},
			assert.NoError},
		{"template file with format flag",
			&queryCmdOptions{
				queryOrFile: "testdata/template1.yaml",
				format:      "yaml",
			},
			&queryCmdOptions{
				queryOrFile: "testdata/template1.yaml",
				format:      "yaml",
			},
			&query{
				Name:        "template1",
				Description: "Example template\n",
				Query:       "r.db('veidemann').table('config_crawl_entities')\n",
				Template:    "{{.id}} {{.meta.name}}\n",
			},
			assert.NoError},
		{"nonexisting template file",
			&queryCmdOptions{
				queryOrFile: "missing.yaml",
			},
			&queryCmdOptions{queryOrFile: "missing.yaml"},
			nil,
			assert.Error},
		{"query",
			&queryCmdOptions{
				queryOrFile: "r.db('veidemann').table('config_crawl_entities')",
			},
			&queryCmdOptions{
				queryOrFile: "r.db('veidemann').table('config_crawl_entities')",
			},
			&query{
				Name:        "",
				Description: "",
				Query:       "r.db('veidemann').table('config_crawl_entities')",
				Template:    "",
			},
			assert.NoError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fieldsFromFlags.parseQuery()
			if !tt.wantErr(t, err, "parseQuery()") {
				return
			}
			assert.Equalf(t, tt.fieldsAfterParse, tt.fieldsFromFlags, "Fields after parseQuery()")
			assert.Equalf(t, tt.want, got, "parseQuery()")
		})
	}
}
