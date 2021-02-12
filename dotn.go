package dotn

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

type Object struct {
	data map[string]interface{}
}

func (obj *Object) Set(key string, value interface{}) {
	if value == nil {
		delete(obj.data, key)
	} else {
		obj.data[key] = value
	}
}

func (obj *Object) Get(key string) (interface{}, bool) {
	v, has := obj.data[key]
	return v, has
}

func (obj *Object) String() string {

	list := make([]string, 0, len(obj.data))

	for k := range obj.data {
		list = append(list, k)
	}

	sort.Strings(list)

	buf := bytes.NewBuffer(nil)

	for _, k := range list {
		fmt.Fprintf(buf, "%s=%v\n", k, obj.data[k])
	}

	return buf.String()
}

func New(v interface{}) (*Object, error) {

	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	var res map[string]interface{}

	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}

	obj := &Object{data: make(map[string]interface{})}

	if err = obj.nodeWork(res, ""); err != nil {
		return nil, err
	}

	return obj, nil
}

func (obj *Object) nodeWork(v interface{}, base string) error {

	pref := base
	if base != "" {
		pref = pref + "."
	}

	switch cv := v.(type) {

	case map[string]interface{}:

		for key, val := range cv {
			obj.nodeWork(val, pref+key)
		}

	default:

		robj := reflect.ValueOf(v)
		rtype := robj.Type()

		switch rtype.Kind() {

		case reflect.Slice, reflect.Array:
			for i := 0; i < robj.Len(); i++ {
				obj.nodeWork(robj.Index(i).Elem().Interface(), pref+strconv.Itoa(i))
			}

		case reflect.Bool, reflect.Float32, reflect.Float64, reflect.Int, reflect.Int8, reflect.Int16,
			reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.String:

			obj.data[base] = v

		default:

			if reflect.ValueOf(v).IsNil() {
				obj.data[base] = nil
			} else {
				return errors.New("invalid")
			}
		}
	}

	return nil
}

func (obj *Object) Node(path string) *Object {

	if !strings.HasSuffix(path, ".") {
		path = path + "."
	}

	size := len(path)

	newObj := &Object{
		data: make(map[string]interface{}),
	}

	for k, v := range obj.data {
		if strings.HasPrefix(k, path) {
			newObj.data[k[size:]] = v
		}
	}

	return newObj
}

func (obj *Object) Fields() []string {

	keys := map[string]bool{}

	var fields []string

	for k := range obj.data {
		key := k
		if idx := strings.Index(key, "."); idx >= 0 {
			key = k[:idx]
		}

		if _, has := keys[key]; !has {
			fields = append(fields, key)
			keys[key] = true
		}
	}

	return fields
}

func (obj *Object) IsArray() bool {

	fields := obj.Fields()

	for _, v := range fields {
		if _, err := strconv.Atoi(v); err != nil {
			return false
		}
	}

	return true
}
