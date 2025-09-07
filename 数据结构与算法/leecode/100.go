package main

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

var a1 = &TreeNode{Val: 1, Left: a2, Right: a3}
var a2 = &TreeNode{Val: 2, Left: nil, Right: nil}
var a3 = &TreeNode{Val: 3, Left: nil, Right: nil}

var b1 = &TreeNode{Val: 1, Left: b2, Right: b3}
var b2 = &TreeNode{Val: 2, Left: nil, Right: nil}
var b3 = &TreeNode{Val: 3, Left: nil, Right: nil}

func isSameTree(p *TreeNode, q *TreeNode) bool {
	if p == nil && q == nil {
		return true
	} else if p == nil || q == nil {
		return false
	}

	if p.Val == q.Val {
		return isSameTree(p.Left, q.Left) && isSameTree(p.Right, q.Right)
	} else {
		return false
	}
}
