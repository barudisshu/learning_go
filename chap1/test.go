package main

import (
	"errors"
	"fmt"
)

func div(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}

	return a / b, nil
}

func main() {
	a, b := 10, 2
	c, err := div(a, b)
	fmt.Println(c, err)
}
