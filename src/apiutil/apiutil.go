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
	"fmt"
	api "github.com/nlnwa/veidemann-api-go/config/v1"
	log "github.com/sirupsen/logrus"
	"reflect"
	"strconv"

	"strings"
)

func createSelector(labelString string) []string {
	var result []string
	if labelString != "" {
		result = strings.Split(labelString, ",")
	}
	return result
}

func createTemplateFilter(filterString string) (*api.FieldMask, *api.ConfigObject, error) {
	q := strings.Split(filterString, "=")
	mask := &api.FieldMask{}
	mask.Paths = append(mask.Paths, q[0])
	fmt.Println(mask)
	obj := &api.ConfigObject{}
	tokens := strings.Split(q[0], ".")

	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	for _, token := range tokens {
		token = strings.Title(token)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
			v = v.Elem()
		}
		if t.Kind() == reflect.Struct {
			if x, ok := t.FieldByName(token); ok {
				t = x.Type
				v = v.FieldByName(token)
				if t.Kind() == reflect.Ptr && v.IsNil() {
					y := reflect.New(t.Elem())
					v.Set(y)
				}
			} else {
				if x, ok := t.FieldByName("Spec"); ok {
					t = x.Type
					v = v.FieldByName("Spec")
					if t.Kind() == reflect.Interface && v.IsNil() {
						val, newVal := makeInstance(token)
						v.Set(val)
						v = newVal
						t = v.Type()
					}
				} else {
					return nil, nil, fmt.Errorf("Field not found: %v", token)
				}
			}
		}
		switch t.Kind() {
		case reflect.Ptr:
			// Nothing to do
		case reflect.String:
			v.Set(reflect.ValueOf(q[1]))
		case reflect.Int64:
			i, _ := strconv.ParseInt(q[1], 10, 64)
			v.Set(reflect.ValueOf(i))
		default:
			log.Fatalf("Field %v of type %v is not implemented yet", token, t.Kind())
		}
	}
	return mask, obj, nil
}

func CreateListRequest(kind api.Kind, ids []string, name string, labelString string, filterString string, pageSize int32, page int32) (*api.ListRequest, error) {
	selector := createSelector(labelString)

	request := &api.ListRequest{}
	request.Kind = kind
	request.Id = ids
	request.NameRegex = name
	request.LabelSelector = selector

	request.Offset = page
	request.PageSize = pageSize

	if filterString != "" {
		m, o, err := createTemplateFilter(filterString)
		if err != nil {
			return nil, err
		}
		request.QueryMask = m
		request.QueryTemplate = o
	}

	return request, nil
}

var typeRegistry = make(map[string]reflect.Type)

func init() {
	co := api.ConfigObject{}
	for _, b := range co.XXX_OneofWrappers() {
		t := reflect.TypeOf(b).Elem()
		n := strings.TrimPrefix(t.String(), "config.ConfigObject_")
		typeRegistry[n] = t
	}
}

func makeInstance(name string) (reflect.Value, reflect.Value) {
	v := reflect.New(typeRegistry[name])
	innerT := v.Elem().Field(0).Type().Elem()
	innerV := reflect.New(innerT)
	v.Elem().Field(0).Set(innerV)
	return v, innerV
}
