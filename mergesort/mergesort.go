package mergesort

import (
	"strings"
)

// Merge
// pre: left is ordered, right is ordered
// post: output is the ordered combination of left and right
// ******************************************
func merge(left []string, right []string) []string {
	output := make([]string, 0)
	l := 0
	r := 0
	for l < len(left) && r < len(right) {
		leftLower := strings.ToLower(left[l])
		rightLower := strings.ToLower(right[r])
		if strings.Compare(leftLower, rightLower) < 0 {
			output = append(output, left[l])
			l++
		} else {
			output = append(output, right[r])
			r++
		}
	}
	output = append(output, left[l:]...)
	output = append(output, right[r:]...)
	return output
}

func mergeSort(input []string, result chan []string) {
	if len(input) == 1 {
		result <- input
		return
	}
	leftChan := make(chan []string)
	rightChan := make(chan []string)
	m := len(input) / 2
	go mergeSort(input[:m], leftChan)
	go mergeSort(input[m:], rightChan)
	leftResult := <-leftChan
	rightResult := <-rightChan
	close(leftChan)
	close(rightChan)
	result <- merge(leftResult, rightResult)
}
func Sort(input []string) []string {
	result := make(chan []string)
	go mergeSort(input, result)
	r := <-result
	return r
}
