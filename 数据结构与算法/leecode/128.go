package main

import (
	"fmt"
)

func longestConsecutive(nums []int) int {
	numSet := map[int]bool{}

	for _, v := range nums {
		numSet[v] = true
	}

	var maxLen int

	for i := range numSet {
		if !numSet[i-1] {
			currentNum := i
			currentLen := 1

			for numSet[currentNum+1] {
				currentNum++
				currentLen++
			}

			maxLen = max(maxLen, currentLen)
		}
	}
	fmt.Println(maxLen)
	return maxLen
}
