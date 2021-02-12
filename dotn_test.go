package dotn

import (
	"fmt"
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

	fmt.Print(obj)

	fmt.Println("-----------")
	fmt.Println(obj.Fields())

	fmt.Println("-----------")

	nobj := obj.Node("map1.map2.list")
	fmt.Print(nobj)

	fmt.Println("-----------")
	fmt.Println(nobj.Fields())
}
