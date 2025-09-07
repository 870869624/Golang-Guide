package main

func trap(height []int) int {

	sumWhite := 0
	max1 := 0
	for i := 0; i < len(height); i++ {
		max1 = max(max1, height[i])
		sumWhite += height[i]
	}

	n := len(height)
	if n == 0 {
		return 0
	}

	left := 0
	right := n - 1
	h := 1

	sumAll := 0
	for h <= max1 {
		for left <= right && height[left] < h {
			left++
		}
		for left <= right && height[right] < h {
			right--
		}
		if left == right {
			sumAll += h
		} else {
			sumAll += (right - left) * h
		}
		h++
	}

	return sumAll - sumWhite
}
