package main

func numOfSubarrays(arr []int, k int, threshold int) int {
	var sum int
	var ave int
	var num int
	for i, v := range arr {
		sum += v
		if i < k-1 {
			continue
		}

		ave = sum / k

		if ave >= threshold {
			num++
		}

		sum -= arr[i-k+1]
	}
	return num
}
