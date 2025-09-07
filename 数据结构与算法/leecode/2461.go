package main

func maximumSubarraySum(nums []int, k int) int64 {
	var sum int64
	var maxSum int64

	num := map[int]int{}
	for i, v := range nums {
		sum += int64(v)
		num[v]++

		if i < k-1 {
			continue
		} else if i == k-1 {
			maxSum = sum
		}

		if len(num) < k {
			maxSum = max(maxSum, 0)
		} else {
			maxSum = max(maxSum, sum)
		}

		num[nums[i-k+1]]--
		if num[nums[i-k+1]] == 0 {
			delete(num, nums[i-k+1])
		}

		sum -= int64(num[i-k+1])
	}
	return maxSum
}
