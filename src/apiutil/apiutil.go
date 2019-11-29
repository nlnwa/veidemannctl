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
	commonsV1 "github.com/nlnwa/veidemann-api-go/commons/v1"
	configV1 "github.com/nlnwa/veidemann-api-go/config/v1"
	log "github.com/sirupsen/logrus"
	"reflect"
	"strconv"

	"strings"
)

func CreateSelector(labelString string) []string {
	var result []string
	if labelString != "" {
		result = strings.Split(labelString, ",")
	}
	return result
}

func CreateTemplateFilter(filterString string, templateObj interface{}) (*commonsV1.FieldMask, interface{}, error) {
	q := strings.Split(filterString, "=")
	mask := &commonsV1.FieldMask{}
	mask.Paths = append(mask.Paths, q[0])
	obj := templateObj
	path := strings.TrimRight(q[0], "+-")
	value := q[1]
	tokens := strings.Split(path, ".")

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
	}

	if t.Kind() == reflect.Slice {
		t = t.Elem()
		val, err := makeValue(t, value)
		if err != nil {
			return nil, nil, err
		}
		if val.IsValid() {
			v.Set(reflect.Append(v, val))
		}
	} else {
		val, err := makeValue(t, value)
		if err != nil {
			return nil, nil, err
		}
		if val.IsValid() {
			v.Set(val)
		}
	}
	return mask, obj, nil
}

func makeValue(t reflect.Type, v string) (val reflect.Value, err error) {
	switch t.Kind() {
	case reflect.Ptr:
		switch t.Elem() {
		case reflect.TypeOf(configV1.ConfigRef{}):
			val = reflect.New(t.Elem())
			cr := val.Interface().(*configV1.ConfigRef)

			kindId := strings.SplitN(v, ":", 2)
			cr.Kind = configV1.Kind(configV1.Kind_value[kindId[0]])
			cr.Id = kindId[1]
		case reflect.TypeOf(configV1.Label{}):
			val = reflect.New(t.Elem())
			cr := val.Interface().(*configV1.Label)

			keyVal := strings.SplitN(v, ":", 2)
			cr.Key = keyVal[0]
			cr.Value = keyVal[1]
		default:
			if typeRegistry[t.Elem().Name()] == nil {
				log.Fatalf("field '%v' of pointer type '%v' is not implemented yet", v, t.Elem())
			}
		}
	case reflect.String:
		val = reflect.ValueOf(v)
	case reflect.Int32:
		n, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			return val, err
		}
		val = reflect.ValueOf(int32(n))
	case reflect.Int64:
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return val, err
		}
		val = reflect.ValueOf(n)
	case reflect.Bool:
		n, err := strconv.ParseBool(v)
		if err != nil {
			return val, err
		}
		val = reflect.ValueOf(n)
	default:
		log.Fatalf("field '%v' of type '%v' is not implemented yet", v, t.Kind())
	}
	return
}

func CreateListRequest(kind configV1.Kind, ids []string, name string, labelString string, filterString string, pageSize int32, page int32) (*configV1.ListRequest, error) {
	selector := CreateSelector(labelString)

	request := &configV1.ListRequest{}
	request.Kind = kind
	request.Id = ids
	request.NameRegex = name
	request.LabelSelector = selector

	request.Offset = page
	request.PageSize = pageSize

	if filterString != "" {
		m, o, err := CreateTemplateFilter(filterString, &configV1.ConfigObject{})
		if err != nil {
			return nil, err
		}
		request.QueryMask = m
		request.QueryTemplate = o.(*configV1.ConfigObject)
	}

	return request, nil
}

var typeRegistry = make(map[string]reflect.Type)

func init() {
	co := configV1.ConfigObject{}
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
