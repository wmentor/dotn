package dotn

import (
	"sort"
	"strings"
	"testing"
)

func TestFilter(t *testing.T) {

	obj, _ := New(map[string]interface{}{"te": 123,
		"map1": map[string]interface{}{
			"hello": "world",
			"map2": map[string]interface{}{
				"key":  123,
				"list": []interface{}{1, 2, 3},
			},
		},
	})

	nobj := obj.Node("map1.map2.list")

	tFields := func(obj *Object, expect []string) {
		list := obj.Fields()
		sort.Strings(list)
		if strings.Join(list, "|") != strings.Join(expect, "|") {
			t.Fatal("invalid list")
		}
	}

	tString := func(obj *Object, wait string) {
		if obj.String() != wait {
			t.Fatalf("Invalid String() result")
		}
	}

	tIsArray := func(obj *Object, wait bool) {
		if obj.IsArray() != wait {
			t.Fatalf("IsArray failed for: %s", obj.String())
		}
	}

	tString(obj, `map1.hello=world
map1.map2.key=123
map1.map2.list.0=1
map1.map2.list.1=2
map1.map2.list.2=3
te=123
`)

	tString(nobj, `0=1
1=2
2=3
`)

	tFields(obj, []string{"map1", "te"})
	tFields(nobj, []string{"0", "1", "2"})

	tIsArray(obj, false)
	tIsArray(nobj, true)
}
