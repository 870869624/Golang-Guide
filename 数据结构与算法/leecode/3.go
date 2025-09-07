package main

import "strings"

func lengthOfLongestSubstring(s string) int {
	str := strings.Split(s, "")

	var len1 int
	var maxLen int
	var m = make(map[string]int)

	for i := range str {
		m[str[i]]++

		len1 = len(m)
		if m[str[i]] > 1 {
			len1 = 0
		}

		if i == 0 {
			maxLen = len1
		} else {
			maxLen = max(maxLen, len1)
		}

	}
	return maxLen
}
