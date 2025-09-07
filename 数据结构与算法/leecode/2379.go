package main

import "strings"

func minimumRecolors(blocks string, k int) int {

	var withe int
	var minWithe int
	strBlocks := strings.Split(blocks, "")

	for i, v := range strBlocks {
		if v == "W" {
			withe++
		}

		if i < k-1 {
			continue
		} else if i == k-1 {
			minWithe = withe
		}

		minWithe = min(minWithe, withe)
		if strBlocks[i-k+1] == "W" {
			withe--
		}
	}
	return minWithe
}
