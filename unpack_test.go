package main

import (
	"os"
	"testing"
	"time"
)

func BenchmarkUnpack_1Thread(b *testing.B) {

	for i := 0; i < b.N; i++ {
		f, _ := os.Open("braid.dat")
		defer f.Close()
		s := time.Now()
		Unpack(f, 1, false)
		b.ReportMetric(float64(time.Since(s).Seconds()), "sec")
	}
}

func BenchmarkUnpack_8Threads(b *testing.B) {

	for i := 0; i < b.N; i++ {
		f, _ := os.Open("braid.dat")
		defer f.Close()
		s := time.Now()
		Unpack(f, 8, false)
		b.ReportMetric(float64(time.Since(s).Seconds()), "sec")
	}
}

func BenchmarkUnpack_16Threads(b *testing.B) {

	for i := 0; i < b.N; i++ {
		f, _ := os.Open("braid.dat")
		defer f.Close()
		s := time.Now()
		Unpack(f, 16, false)
		b.ReportMetric(float64(time.Since(s).Seconds()), "sec")
	}
}

func BenchmarkUnpack_32Threads(b *testing.B) {

	for i := 0; i < b.N; i++ {
		f, _ := os.Open("braid.dat")
		defer f.Close()
		s := time.Now()
		Unpack(f, 32, false)
		b.ReportMetric(float64(time.Since(s).Seconds()), "sec")
	}
}

func BenchmarkUnpack_64Threads(b *testing.B) {

	for i := 0; i < b.N; i++ {
		f, _ := os.Open("braid.dat")
		defer f.Close()
		s := time.Now()
		Unpack(f, 64, false)
		b.ReportMetric(float64(time.Since(s).Seconds()), "sec")
	}
}

func BenchmarkUnpack_100Threads(b *testing.B) {

	for i := 0; i < b.N; i++ {
		f, _ := os.Open("braid.dat")
		defer f.Close()
		s := time.Now()
		Unpack(f, 100, false)
		b.ReportMetric(float64(time.Since(s).Seconds()), "sec")
	}
}
