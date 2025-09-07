package main

func getBig(a int, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}
func maxSales(sales []int) int {
	maxSum := 0
	maxLineSum := sales[0]
	for i := 0; i < len(sales); i++ {
		maxSum = getBig(sales[i]+maxSum, sales[i])
		maxLineSum = getBig(maxLineSum, maxSum)
	}
	return maxLineSum
}
