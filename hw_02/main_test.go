package main

import "testing"

func BenchmarkFibonacciIterative(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FibonacciIterative(6)
	}
}

func BenchmarkFibonacciRecursive(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FibonacciRecursive(6)
	}
}

func BenchmarkIsPrime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsPrime(6)
	}
}

func BenchmarkIsBinaryPalindrome(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsBinaryPalindrome(5)
	}
}

func BenchmarkValidParentheses(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ValidParentheses("func() { return fmt.Println(len([]int{1,2,3}))}")
	}
}

func BenchmarkIncrement(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Increment("111")
	}
}
