package main

import "fmt"

// 定长滑动窗口
func maxVowels(s string, k int) (ans int) {
	vowel := 0
	for i, in := range s {
		if in == 'a' || in == 'e' || in == 'i' || in == 'o' || in == 'u' {
			vowel++
		}

		fmt.Println("走到这里", i)
		if i < k-1 {
			continue
		}

		ans = max(ans, vowel)

		out := s[i-k+1]
		fmt.Println(string(out), i)
		if out == 'a' || out == 'e' || out == 'i' || out == 'o' || out == 'u' {
			vowel--
		}
	}
	return ans
}
