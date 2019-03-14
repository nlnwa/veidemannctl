// Copyright © 2019 National Library of Norway
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

package importutil

import (
	configV1 "github.com/nlnwa/veidemann-api-go/config/v1"
	"reflect"
	"testing"
)

func TestImportDb_stringArrayToBytes(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want []byte
	}{
		{"1", []string{"foo", "bar"}, []byte{12, 255, 129, 2, 1, 2, 255, 130, 0, 1, 12, 0, 0, 12, 255, 130, 0, 2, 3, 102, 111, 111, 3, 98, 97, 114}},
		{"2", []string{}, []byte{12, 255, 129, 2, 1, 2, 255, 130, 0, 1, 12, 0, 0, 4, 255, 130, 0, 0}},
		{"3", []string{"", "foo", "bar"}, []byte{12, 255, 129, 2, 1, 2, 255, 130, 0, 1, 12, 0, 0, 13, 255, 130, 0, 3, 0, 3, 102, 111, 111, 3, 98, 97, 114}},
		{"4", []string{""}, []byte{12, 255, 129, 2, 1, 2, 255, 130, 0, 1, 12, 0, 0, 5, 255, 130, 0, 1, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &ImportDb{}
			if got := d.stringArrayToBytes(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ImportDb.stringArrayToBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImportDb_bytesToStringArray(t *testing.T) {
	tests := []struct {
		name string
		args []byte
		want []string
	}{
		{"1", []byte{12, 255, 129, 2, 1, 2, 255, 130, 0, 1, 12, 0, 0, 12, 255, 130, 0, 2, 3, 102, 111, 111, 3, 98, 97, 114}, []string{"foo", "bar"}},
		{"2", []byte{}, []string{}},
		{"3", []byte{12, 255, 129, 2, 1, 2, 255, 130, 0, 1, 12, 0, 0, 4, 255, 130, 0, 0}, []string{}},
		{"4", []byte{12, 255, 129, 2, 1, 2, 255, 130, 0, 1, 12, 0, 0, 13, 255, 130, 0, 3, 0, 3, 102, 111, 111, 3, 98, 97, 114}, []string{"", "foo", "bar"}},
		{"5", []byte{12, 255, 129, 2, 1, 2, 255, 130, 0, 1, 12, 0, 0, 5, 255, 130, 0, 1, 0}, []string{""}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &ImportDb{}
			if got := d.bytesToStringArray(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ImportDb.bytesToStringArray() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExistsCode_ExistsInVeidemann(t *testing.T) {
	tests := []struct {
		name string
		e    ExistsCode
		want bool
	}{
		{"error", ERROR, false},
		{"new", NEW, false},
		{"duplicate_new", DUPLICATE_NEW, false},
		{"exists_veidemann", EXISTS_VEIDEMANN, true},
		{"duplicate_veidemann", DUPLICATE_VEIDEMANN, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.ExistsInVeidemann(); got != tt.want {
				t.Errorf("ExistsCode.ExistsInVeidemann() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImportDb_CheckAndUpdateVeidemann(t *testing.T) {
	s1 := "https://www.eiksenteret.no"
	s2 := "http://www.eiksenteret.no"
	s3 := "http://www.eiksenteret.no/foo"
	s4 := "https://www.foo.no"
	s5 := "http://www.foo.no"
	s6 := "http://www.foo.no/foo"
	f := func(client configV1.ConfigClient, data interface{}) (id string, err error) {
		switch data {
		case s1:
			return "s1", nil
		case s2:
			return "s2", nil
		case s3:
			return "s3", nil
		case s5:
			return "s5", nil
		case s6:
			return "s6", nil
		default:
			return "", nil
		}
	}

	type args struct {
		uri        string
		data       interface{}
		createFunc func(client configV1.ConfigClient, data interface{}) (id string, err error)
	}
	tests := []struct {
		name    string
		args    args
		want    *ExistsResponse
		wantErr bool
	}{
		{"first", args{s1, s1, f}, &ExistsResponse{Code: NEW}, false},
		{"duplicate", args{s2, s2, f}, &ExistsResponse{EXISTS_VEIDEMANN, []string{"s1"}}, false},
		{"duplicate_with_path", args{s3, s3, f}, &ExistsResponse{EXISTS_VEIDEMANN, []string{"s1"}}, false},
		{"no_id_first", args{s4, s4, f}, &ExistsResponse{Code: NEW}, false},
		{"no_id_duplicate", args{s5, s5, f}, &ExistsResponse{Code: DUPLICATE_NEW}, false},
		{"no_id_duplicate_with_path", args{s6, s6, f}, &ExistsResponse{EXISTS_VEIDEMANN, []string{"s5"}}, false},
	}

	d := NewImportDb(nil, "/tmp/vmtest", true)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := d.CheckAndUpdateVeidemann(tt.args.uri, tt.args.data, tt.args.createFunc)
			if (err != nil) != tt.wantErr {
				t.Errorf("ImportDb.CheckAndUpdateVeidemann() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ImportDb.CheckAndUpdateVeidemann() = %v, want %v", got, tt.want)
			}
		})
	}
}
