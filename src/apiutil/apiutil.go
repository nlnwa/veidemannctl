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
	commonsV1 "github.com/nlnwa/veidemann-api/go/commons/v1"
	configV1 "github.com/nlnwa/veidemann-api/go/config/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
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

func CreateTemplateFilter(filterString string, templateObj proto.Message) (*commonsV1.FieldMask, proto.Message, error) {
	q := strings.Split(filterString, "=")
	mask := &commonsV1.FieldMask{}
	mask.Paths = append(mask.Paths, q[0])
	obj := templateObj
	path := strings.TrimRight(q[0], "+-")
	value := q[1]
	tokens := strings.Split(path, ".")

	var fieldType protoreflect.FieldDescriptor
	msgType := obj.ProtoReflect().Descriptor()

	for _, token := range tokens {
		fieldType = msgType.Fields().ByJSONName(token)
		if fieldType == nil {
			sa := make([]string, msgType.Fields().Len())
			for i := 0; i < msgType.Fields().Len(); i++ {
				sa[i] = msgType.Fields().Get(i).JSONName()
			}
			return nil, templateObj, fmt.Errorf("no field with name '%s' in '%s'. Valid field names: %s", token, msgType.FullName(), strings.Join(sa, ", "))
		}
		switch fieldType.Kind() {
		case protoreflect.BoolKind:
			v, err := strconv.ParseBool(value)
			if err != nil {
				return nil, templateObj, fmt.Errorf("error converting %v to boolean: %w", value, err)
			}
			setValue(protoreflect.ValueOfBool(v), obj, fieldType)

		case protoreflect.Int32Kind:
			fallthrough
		case protoreflect.Sint32Kind:
			fallthrough
		case protoreflect.Sfixed32Kind:
			v, err := strconv.ParseInt(value, 0, 32)
			if err != nil {
				return nil, templateObj, fmt.Errorf("error converting %v to int: %w", value, err)
			}
			setValue(protoreflect.ValueOfInt32(int32(v)), obj, fieldType)

		case protoreflect.Int64Kind:
			fallthrough
		case protoreflect.Sint64Kind:
			fallthrough
		case protoreflect.Sfixed64Kind:
			v, err := strconv.ParseInt(value, 0, 64)
			if err != nil {
				return nil, templateObj, fmt.Errorf("error converting %v to int: %w", value, err)
			}
			setValue(protoreflect.ValueOfInt64(v), obj, fieldType)

		case protoreflect.Uint32Kind:
			fallthrough
		case protoreflect.Fixed32Kind:
			v, err := strconv.ParseUint(value, 0, 32)
			if err != nil {
				return nil, templateObj, fmt.Errorf("error converting %v to uint: %w", value, err)
			}
			setValue(protoreflect.ValueOfUint32(uint32(v)), obj, fieldType)

		case protoreflect.Uint64Kind:
			fallthrough
		case protoreflect.Fixed64Kind:
			v, err := strconv.ParseUint(value, 0, 64)
			if err != nil {
				return nil, templateObj, fmt.Errorf("error converting %v to uint: %w", value, err)
			}
			setValue(protoreflect.ValueOfUint64(v), obj, fieldType)

		case protoreflect.FloatKind:
			v, err := strconv.ParseFloat(value, 32)
			if err != nil {
				return nil, templateObj, fmt.Errorf("error converting %v to float: %w", value, err)
			}
			setValue(protoreflect.ValueOfFloat32(float32(v)), obj, fieldType)

		case protoreflect.DoubleKind:
			v, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, templateObj, fmt.Errorf("error converting %v to float: %w", value, err)
			}
			setValue(protoreflect.ValueOfFloat64(v), obj, fieldType)

		case protoreflect.StringKind:
			setValue(protoreflect.ValueOfString(value), obj, fieldType)

		case protoreflect.BytesKind:
			setValue(protoreflect.ValueOfBytes([]byte(value)), obj, fieldType)

		case protoreflect.EnumKind:
			enumVal := fieldType.Enum().Values().ByName(protoreflect.Name(value))
			if enumVal == nil {
				sa := make([]string, fieldType.Enum().Values().Len())
				for i := 0; i < fieldType.Enum().Values().Len(); i++ {
					sa[i] = string(fieldType.Enum().Values().Get(i).Name())
				}
				return nil, templateObj, fmt.Errorf("not a valid enum value '%s' for '%s'. Valid values: %s", value, fieldType.FullName(), strings.Join(sa, ", "))
			}
			setValue(protoreflect.ValueOfEnum(enumVal.Number()), obj, fieldType)

		case protoreflect.MessageKind:
			item := obj.ProtoReflect().NewField(fieldType)
			var message protoreflect.Message

			if fieldType.IsList() {
				v := item.List().NewElement()
				item.List().Append(v)
				message = v.Message()
			} else {
				message = item.Message()
			}
			msgType = fieldType.Message()

			switch m := message.Interface().(type) {
			case *configV1.ConfigRef:
				kindId := strings.SplitN(value, ":", 2)
				if len(kindId) != 2 {
					return nil, templateObj, fmt.Errorf("not a valid configRef value %v for: %v. ConfigRef should have format [kind]:[id]", value, fieldType.FullName())
				}
				m.Kind = configV1.Kind(configV1.Kind_value[kindId[0]])
				m.Id = kindId[1]
			case *configV1.Label:
				keyVal := strings.SplitN(value, ":", 2)
				if len(keyVal) != 2 {
					return nil, templateObj, fmt.Errorf("not a valid label value %v for: %v. Label should have format [name]:[value]", value, fieldType.FullName())
				}
				m.Key = keyVal[0]
				m.Value = keyVal[1]
			}
			obj.ProtoReflect().Set(fieldType, item)
			obj = message.Interface()
		}
	}

	return mask, templateObj, nil
}

func setValue(v protoreflect.Value, obj proto.Message, fieldType protoreflect.FieldDescriptor) {
	if fieldType.IsList() {
		list := obj.ProtoReflect().NewField(fieldType)
		list.List().Append(v)
		v = list
	}
	obj.ProtoReflect().Set(fieldType, v)
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
