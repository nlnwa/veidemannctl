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

import "testing"

func TestHasSuffix(t *testing.T) {
	type args struct {
		s      string
		suffix []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"1", args{"foo.yml", []string{".yml"}}, true},
		{"2", args{"foo.yml", []string{".yaml", ".yml"}}, true},
		{"3", args{"foo.yml", []string{".yaml", ".json"}}, false},
		{"4", args{"foo.yml", []string{".yaml"}}, false},
		{"5", args{"foo.yml", []string{}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasSuffix(tt.args.s, tt.args.suffix...); got != tt.want {
				t.Errorf("HasSuffix() = %v, want %v", got, tt.want)
			}
		})
	}
}
