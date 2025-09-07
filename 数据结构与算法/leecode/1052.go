package main

func maxSatisfied(customers []int, grumpy []int, minutes int) int {
	var sum int
	len := len(customers)
	for i := 0; i < len; i++ {
		if grumpy[i] == 0 { //如果不生气
			sum += customers[i]
			customers[i] = 0
		}
	}

	var pre int
	var maxPre int
	for i, v := range customers {
		pre += v

		if i < minutes-1 {
			continue
		}
		if i == minutes-1 {
			maxPre = pre
		}

		maxPre = max(maxPre, pre)
		pre -= customers[i-minutes+1]
	}

	return sum + maxPre
}
