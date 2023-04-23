package unsafe

import (
	"fmt"
	"reflect"
)

func iterateFields(entity any) {
	typ := reflect.TypeOf(entity)
	numField := typ.NumField()
	for i := 0; i < numField; i++ {
		f := typ.Field(i)
		fmt.Println(fmt.Sprintf("Field:%s,Offset:%v", f.Name, f.Offset))
	}
}
