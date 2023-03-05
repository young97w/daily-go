package main

import "fmt"

func main() {
	f := float64(3.14159)
	fmt.Println(printNumWith2(f))

	data := []byte("golang")
	fmt.Println(printBytes(data))
}

// 输出两位小数
func printNumWith2(float642 float64) string {
	return fmt.Sprintf("%.2f", float642)
}

func printBytes(data []byte) string {
	return fmt.Sprintf("%x", data)
}
