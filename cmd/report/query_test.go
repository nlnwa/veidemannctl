package report

import (
	reportV1 "github.com/nlnwa/veidemann-api/go/report/v1"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_queryCmdOptions_parseQuery(t *testing.T) {
	tests := []struct {
		name             string
		fieldsFromFlags  *queryCmdOptions
		args             []string
		fieldsAfterParse *queryCmdOptions
		want             *query
		wantErr          assert.ErrorAssertionFunc
	}{
		{"template file",
			&queryCmdOptions{},
			[]string{"testdata/template1.yaml"},
			&queryCmdOptions{
				format: "template",
			},
			&query{
				Name:        "template1",
				Description: "Example template\n",
				Query:       "r.db('veidemann').table('config_crawl_entities')\n",
				Template:    "{{.id}} {{.meta.name}}\n",
				queryOrFile: "testdata/template1.yaml",
				queryArgs:   make([]any, 0),
				request: &reportV1.ExecuteDbQueryRequest{
					Query: "r.db('veidemann').table('config_crawl_entities')\n",
				},
			},
			assert.NoError},
		{"template file with template flag",
			&queryCmdOptions{
				goTemplate: "{{.id}}",
			},
			[]string{"testdata/template1.yaml"},
			&queryCmdOptions{
				goTemplate: "{{.id}}",
				format:     "template",
			},
			&query{
				Name:        "template1",
				Description: "Example template\n",
				Query:       "r.db('veidemann').table('config_crawl_entities')\n",
				Template:    "{{.id}}",
				queryOrFile: "testdata/template1.yaml",
				queryArgs:   make([]any, 0),
				request: &reportV1.ExecuteDbQueryRequest{
					Query: "r.db('veidemann').table('config_crawl_entities')\n",
				},
			},
			assert.NoError},
		{"template file with format flag",
			&queryCmdOptions{
				format: "yaml",
			},
			[]string{"testdata/template1.yaml"},
			&queryCmdOptions{
				format: "yaml",
			},
			&query{
				Name:        "template1",
				Description: "Example template\n",
				Query:       "r.db('veidemann').table('config_crawl_entities')\n",
				Template:    "{{.id}} {{.meta.name}}\n",
				queryOrFile: "testdata/template1.yaml",
				queryArgs:   make([]any, 0),
				request: &reportV1.ExecuteDbQueryRequest{
					Query: "r.db('veidemann').table('config_crawl_entities')\n",
				},
			},
			assert.NoError},
		{"nonexisting template file",
			&queryCmdOptions{},
			[]string{"missing.yaml"},
			&queryCmdOptions{},
			nil,
			assert.Error},
		{"query",
			&queryCmdOptions{},
			[]string{"r.db('veidemann').table('config_crawl_entities')"},
			&queryCmdOptions{},
			&query{
				Query: "r.db('veidemann').table('config_crawl_entities')",
				request: &reportV1.ExecuteDbQueryRequest{
					Query: "r.db('veidemann').table('config_crawl_entities')",
				},
			},
			assert.NoError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := tt.fieldsFromFlags
			got, err := o.parseQuery(tt.args)
			if !tt.wantErr(t, err, "parseQuery()") {
				return
			}
			assert.Equalf(t, tt.fieldsAfterParse, tt.fieldsFromFlags, "Fields after parseQuery()")
			assert.Equalf(t, tt.want, got, "parseQuery()")
		})
	}
}
