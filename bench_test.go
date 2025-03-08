package main

import (
	"os"
	"testing"
)

func BenchmarkUnpack(b *testing.B) {

	for i := 0; i < b.N; i++ {
		f, _ := os.Open("braid.dat")
		defer f.Close()
		Unpack(f)
	}
}
