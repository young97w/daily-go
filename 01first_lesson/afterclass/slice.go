package main

import "fmt"

func main() {
	s := []int{1, 2, 4, 7}
	// 结果应该是 5, 1, 2, 4, 7
	s, _ = Add(s, 0, 5)

	// 结果应该是5, 9, 1, 2, 4, 7
	s, _ = Add(s, 1, 9)

	// 结果应该是5, 9, 1, 2, 4, 7, 13
	s, _ = Add(s, 6, 13)

	// 结果应该是5, 9, 2, 4, 7, 13
	s, _ = Delete(s, 2)

	// 结果应该是9, 2, 4, 7, 13
	s, _ = Delete(s, 0)

	// 结果应该是9, 2, 4, 7
	s, _ = Delete(s, 4)

}

func Add(s []int, index int, value int) ([]int, bool) {
	n := len(s)
	if index < 0 || index > n {
		return nil, false
	}
	s = append(s[:index], append([]int{value}, s[index:]...)...)
	fmt.Println(s)
	return s, true
}

func Delete(s []int, index int) ([]int, bool) {
	n := len(s)
	if index < 0 || index > n-1 {
		return nil, false
	}
	s = append(s[:index], s[index+1:]...)
	fmt.Println(s)
	return s, true
}
