package main

func findMaxAverage(nums []int, k int) float64 {
	var sum int
	var ave float64
	var maxAve float64
	var flag int
	for i, v := range nums {
		sum += v
		if i < k-1 {
			continue
		}
		flag++
		ave = float64(sum) / float64(k)
		maxAve = max(ave, maxAve)

		if flag == 1 {
			maxAve = ave
		}

		sum -= nums[i-k+1]
	}
	return maxAve
}
