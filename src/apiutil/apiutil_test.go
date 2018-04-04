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

package apiutil

import (
	"reflect"
	"testing"

	api "github.com/nlnwa/veidemannctl/veidemann_api"
)

func TestCreateSelector(t *testing.T) {
	type args struct {
		labelString string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Empty label",
			args: args{""},
			want: nil,
		},
		{
			name: "Single label",
			args: args{"foo:bar"},
			want: []string{"foo:bar"},
		},
		{
			name: "Multiple labels",
			args: args{"foo:bar,lab"},
			want: []string{"foo:bar", "lab"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateSelector(tt.args.labelString); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateSelector() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateListRequest(t *testing.T) {
	type args struct {
		ids         []string
		name        string
		labelString string
		pageSize    int32
		page        int32
	}
	tests := []struct {
		name string
		args args
		want api.ListRequest
	}{
		{
			name: "One Id",
			args: args{
				[]string{"id1"},
				"",
				"",
				0,
				0,
			},
			want: api.ListRequest{Id: []string{"id1"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateListRequest(tt.args.ids, tt.args.name, tt.args.labelString, tt.args.pageSize, tt.args.page); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateListRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
