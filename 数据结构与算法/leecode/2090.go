package main

func getAverages(nums []int, k int) []int {

	var sum int64
	var ave int64

	res := make([]int, len(nums))

	for i := range res {
		res[i] = -1
	}

	for i, v := range nums {

		sum += int64(v)

		if i < 2*k {
			continue
		}

		ave = sum / (2*int64(k) + 1)

		res[i-k] = int(ave)

		sum -= int64(nums[i-2*k])
	}

	return res
}
