package main

import (
	"strings"
)

//	func main() {
//		str := "[[[]"
//		fmt.Println(isValid(str))
//	}

func isValid(s string) bool {
	stack := make([]string, 0)
	if s == "" {
		return true
	}

	str := strings.Split(s, "")
	if len(str)%2 != 0 {
		return false
	}

	if strings.Contains(s, "(") && !strings.Contains(s, ")") || strings.Contains(s, ")") && !strings.Contains(s, "(") {
		return false
	}
	if strings.Contains(s, "{") && !strings.Contains(s, "}") || strings.Contains(s, "}") && !strings.Contains(s, "{") {
		return false
	}
	if strings.Contains(s, "[") && !strings.Contains(s, "]") || strings.Contains(s, "]") && !strings.Contains(s, "[") {
		return false
	}

	for i := 0; i < len(str); i++ {
		if str[i] == "(" {
			str1 := ")"
			stack = append(stack, str1)
		} else if str[i] == "{" {
			str1 := "}"
			stack = append(stack, str1)
		} else if str[i] == "[" {
			str1 := "]"
			stack = append(stack, str1)
		} else {
			if len(stack) == 0 {
				return false
			}
			if str[i] != stack[len(stack)-1] {
				return false
			} else {
				stack = stack[:len(stack)-1]
			}
		}
	}
	return len(stack) == 0
}
