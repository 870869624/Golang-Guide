package main

func maxFreeTime(eventTime int, k int, startTime []int, endTime []int) int {
	var ans int
	n := len(startTime) //有多少个开会时间断

	get := func(i int) int {
		if i == 0 {
			return startTime[0]
		}

		if i == n {
			return eventTime - endTime[n-1]
		}
		return startTime[i] - endTime[i-1]
	}

	s := 0

	for i := range n + 1 {
		s += get(i)
		if i < k {
			continue
		}

		ans = max(ans, s)
		s -= get(i - k)
	}

	return ans

}
