package main

import "fmt"

func main() {
	req := map[string]interface{}{"type": "1", "delay": "delay", "reason": "reason", "target": []string{"1", "2", "3"}}
	fmt.Println(req)
}
