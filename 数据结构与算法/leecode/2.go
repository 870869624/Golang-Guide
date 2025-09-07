package main

type ListNode struct {
	Val  int
	Next *ListNode
}

func addTwoNumbers(l1 *ListNode, l2 *ListNode) *ListNode {
	if l1 == nil {
		return l2
	} else if l2 == nil {
		return l1
	}

	var result = &ListNode{}
	result.Val = result.Val + l1.Val + l2.Val
	if result.Val >= 10 {
		result.Val -= 10
		if l1.Next != nil && l2.Next == nil {
			l1.Next.Val++
			l2.Next = &ListNode{Val: 0, Next: nil}
		} else if l1.Next == nil && l2.Next != nil {
			l2.Next.Val++
			l1.Next = &ListNode{Val: 0, Next: nil}
		} else if l1.Next == nil && l2.Next == nil {
			l1.Next = &ListNode{Val: 1, Next: nil}
		} else {
			l1.Next.Val++
		}
	}
	result.Next = addTwoNumbers(l1.Next, l2.Next)
	return result
}
