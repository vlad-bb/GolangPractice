package main

import (
	"hw_03/documentstore"
	"testing"
)

func BenchmarkPut(b *testing.B) {
	for i := 0; i < b.N; i++ {
		documentstore.Put(firstDoc)
	}
}

func BenchmarkGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		documentstore.Get("first")
	}
}

func BenchmarkList(b *testing.B) {
	for i := 0; i < b.N; i++ {
		documentstore.List()
	}
}

func BenchmarkDelete(b *testing.B) {
	for i := 0; i < b.N; i++ {
		documentstore.Delete("first")
	}
}
