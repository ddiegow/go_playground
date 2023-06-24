package pointers

import (
	"fmt"
)

type data struct {
	num   int
	left  *data
	right *data
}

func Run() {
	var test = make([]data, 3)
	test[0].num = 0
	test[0].left = nil
	test[0].right = &test[1]
	test[1].num = 1
	test[1].left = &test[0]
	test[1].right = &test[2]
	test[2].num = 2
	test[2].left = &test[1]
	test[2].right = &test[1]
	fmt.Println(test[0].right.num) // 1
	fmt.Println(test[1].left.num)  // 0
	fmt.Println(test[1].right.num) // 2
}
