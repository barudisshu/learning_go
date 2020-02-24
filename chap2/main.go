package main

import "strconv"

func main() {

	x, _ := strconv.Atoi("12")		// 忽略Atoi的err返回值
	println(x)
}
