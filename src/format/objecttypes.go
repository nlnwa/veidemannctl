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
	configV1 "github.com/nlnwa/veidemann-api/go/config/v1"
	"sort"
	"strings"
)

// Get kind for string
func GetKind(Name string) configV1.Kind {
	Name = strings.ToLower(Name)
	for _, k := range configV1.Kind_name {
		if strings.ToLower(k) == Name {
			return configV1.Kind(configV1.Kind_value[k])
		}
	}
	return configV1.Kind_undefined
}

func GetObjectNames() []string {
	result := make([]string, len(configV1.Kind_name)-1)
	idx := 0
	for _, n := range configV1.Kind_name {
		if n != "undefined" {
			result[idx] = n
			idx++
		}
	}
	sort.Strings(result)
	return result
}
