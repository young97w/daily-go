package afterclass

import (
	"fmt"
	"testing"
)

func TestMySet_Put(t *testing.T) {
	var set MySet
	set.NewSet()
	for _, v := range []string{"a", "b", "c"} {
		set.Put(v)
	}
	fmt.Println(set.M)
	fmt.Println(set.Contains("a"))
	set.Remove("a")
	fmt.Println(set.Keys())
	absent, a := set.PutIfAbsent("a")
	fmt.Println(absent, a)
	ifAbsent, b := set.PutIfAbsent("a")
	fmt.Println(ifAbsent, b)
}
