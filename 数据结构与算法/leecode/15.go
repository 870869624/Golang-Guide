package main

import (
	"sort"
)

// 三重循环，利用双指针来解题
// 固定第一个数，双指针分别指向第二个数和第三个数，然后利用双指针来寻找满足条件的数
// 第二个数变大，第三个数肯定会变小
func threeSum(nums []int) [][]int {
	n := len(nums)
	ans := make([][]int, 0)
	sort.Ints(nums)

	for first := 0; first < n; first++ {
		third := n - 1
		target := -1 * nums[first]

		for second := first + 1; second < n; second++ {
			if second > first+1 && nums[second] == nums[second-1] {
				continue
			}
			for second < third && nums[second]+nums[third] > target {
				third--
			}
			if second == third {
				break
			}
			if nums[second]+nums[third] == target {
				ans = append(ans, []int{nums[first], nums[second], nums[third]})
			}
		}
	}
	return ans
}
