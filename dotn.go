package dotn

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// object is realizing dot notation
type Object struct {
	data map[string]interface{}
}

// func set value by dot notation
func (obj *Object) Set(key string, value interface{}) {
	if value == nil {
		delete(obj.data, key)
	} else {
		obj.data[key] = value
	}
}

// get value by dot notation
func (obj *Object) Get(key string) (interface{}, bool) {
	v, has := obj.data[key]
	return v, has
}

// return data as string format
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

// create dotn.Object from custom struct or map
func New(v interface{}) (*Object, error) {

	codec := NewJsonCodec()

	data, err := codec.Marshal(v)
	if err != nil {
		return nil, err
	}

	var res map[string]interface{}

	err = codec.Unmarshal(data, &res)
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

// return subnode via dot notation
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

// return top level fields or array indexes
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

// check object "is array"
func (obj *Object) IsArray() bool {

	fields := obj.Fields()

	for _, v := range fields {
		if _, err := strconv.Atoi(v); err != nil {
			return false
		}
	}

	return true
}

// return data object interface
func (obj *Object) Interface() interface{} {

	fields := obj.Fields()

	if obj.IsArray() {

		list := []interface{}{}

		indexes := make([]int, len(fields))

		for i, f := range fields {
			indexes[i], _ = strconv.Atoi(f)
		}

		sort.Ints(indexes)

		for _, idx := range indexes {
			key := strconv.Itoa(idx)

			if v, has := obj.data[key]; has {
				list = append(list, v)
			} else {
				nobj := obj.Node(key)
				list = append(list, nobj.Interface())
			}
		}

		return list
	}

	res := make(map[string]interface{})

	for _, key := range fields {
		if v, has := obj.data[key]; has {
			res[key] = v
		} else {
			nobj := obj.Node(key)
			res[key] = nobj.Interface()
		}
	}

	return res
}

// delete sub node via dot notation
func (obj *Object) Delete(path string) {

	if _, has := obj.data[path]; has {
		delete(obj.data, path)
	} else {
		res := make(map[string]interface{})

		for k, v := range obj.data {
			if !strings.HasPrefix(k, path+".") {
				res[k] = v
			}
		}

		obj.data = res
	}
}

// decode object to input value via json.Unmarshal
func (obj *Object) Decode(v interface{}) error {
	codec := NewJsonCodec()

	data, err := codec.Marshal(obj.data)
	if err != nil {
		return err
	}

	return codec.Unmarshal(data, v)
}
