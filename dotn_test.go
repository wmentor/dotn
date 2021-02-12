package dotn

import (
	"fmt"
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
				"list2": []interface{}{map[string]interface{}{"tl1": true, "tl2": false},
					"13",
				},
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
			t.Fatalf("Invalid String() result:\n%s", obj.String())
		}
	}

	tIsArray := func(obj *Object, wait bool) {
		if obj.IsArray() != wait {
			t.Fatalf("IsArray failed for: %s", obj.String())
		}
	}

	tInterface := func(obj *Object, wait string) {
		str := fmt.Sprintf("%v", obj.Interface())
		if str != wait {
			t.Fatalf("result=%s", str)
		}
	}

	tString(obj, `map1.hello=world
map1.map2.key=123
map1.map2.list.0=1
map1.map2.list.1=2
map1.map2.list.2=3
map1.map2.list2.0.tl1=true
map1.map2.list2.0.tl2=false
map1.map2.list2.1=13
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

	tInterface(obj, "map[map1:map[hello:world map2:map[key:123 list:[1 2 3] list2:[map[tl1:true tl2:false] 13]]] te:123]")
	tInterface(nobj, "[1 2 3]")

	obj.Delete("map1.map2.list")
	tInterface(obj, "map[map1:map[hello:world map2:map[key:123 list2:[map[tl1:true tl2:false] 13]]] te:123]")
	obj.Delete("map1.hello")
	tInterface(obj, "map[map1:map[map2:map[key:123 list2:[map[tl1:true tl2:false] 13]]] te:123]")
}
