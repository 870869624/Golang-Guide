package main

func main() {
	maxSum([]int{1, 2, 2}, 2, 2)

}
func maxSum(nums []int, m int, k int) int64 {
	var sum int64
	var maxSum int64
	cnt := map[int]int{}

	for i, v := range nums {
		sum += int64(nums[i])

		cnt[v]++

		if i < k-1 {
			continue
		}

		if i == 0 {
			maxSum = sum
		} else if len(cnt) >= m {
			maxSum = max(maxSum, sum)
		}

		cnt[nums[i-k+1]]--

		if cnt[nums[i-k+1]] == 0 {
			delete(cnt, nums[i-k+1])
		}

		sum -= int64(nums[i-k+1])
	}
	return maxSum
}
