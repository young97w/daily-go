package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	s := "你好"
	//get length of bytes
	fmt.Println(len(s))
	//get length of characters
	fmt.Println(utf8.RuneCountInString(s))
}
