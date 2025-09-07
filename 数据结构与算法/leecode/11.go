package main

import "fmt"

func maxArea(height []int) int {
	right := len(height) - 1
	left := 0

	var maxArea int

	for left < right {

		if height[left] < height[right] {
			maxArea = max(maxArea, height[left]*(right-left))
			left++
		} else {
			maxArea = max(maxArea, height[right]*(right-left))
			right--
		}
	}

	fmt.Println(maxArea)
	return maxArea
}
