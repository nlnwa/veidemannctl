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

package util

import (
	api "github.com/nlnwa/veidemannctl/veidemann_api"
	"strings"
)

func CreateSelector(labelString string) []string {
	var result []string
	if labelString != "" {
		result = strings.Split(labelString, ",")
	}
	return result
}

func CreateListRequest(ids []string, name string, labelString string, pageSize int32, page int32) api.ListRequest {
	selector := CreateSelector(labelString)

	request := api.ListRequest{}
	request.Id = ids
	request.Name = name
	request.LabelSelector = selector
	request.Page = page
	request.PageSize = pageSize

	return request
}
