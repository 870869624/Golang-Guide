package main

func maxScore(cardPoints []int, k int) int {
	var sum int
	var minSum int
	var maxSum int

	for _, v := range cardPoints {
		maxSum += v
	}

	if k == len(cardPoints) {
		return maxSum
	}

	n := len(cardPoints)
	winSize := n - k

	for i := range cardPoints {

		sum += cardPoints[i]
		if i < winSize-1 {
			continue
		}

		if i == winSize-1 {
			minSum = sum
		}

		minSum = min(minSum, sum)
		sum -= cardPoints[i-winSize+1]
	}

	return maxSum - minSum
}
