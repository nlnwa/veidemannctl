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
