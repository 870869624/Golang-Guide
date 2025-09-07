package main

func decrypt(code []int, k int) []int {
	n := len(code)
	res := make([]int, n)
	var sum int

	if k == 0 {
		return res
	}

	code = append(code, code...)

	l, r := 1, k
	if k < 0 {
		l, r = n+k, n-1
	}

	for _, v := range code[l : r+1] {
		sum += v
	}

	for i := range res {
		res[i] = sum
		sum -= code[l]
		sum += code[r+1]
		l, r = l+1, r+1
	}
	return res
}
