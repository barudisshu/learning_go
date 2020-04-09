package main

var x = 0x100

const y = 0x200

func main() {
	println(&x, x)
	println(&y, y) // 错误：cannot take the address of y
}
