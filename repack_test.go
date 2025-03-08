package main

import (
	"os"
	"testing"
	"time"
)

func BenchmarkRepack_1Thread(b *testing.B) {

	for i := 0; i < b.N; i++ {
		f, _ := os.Open("braid.dat")
		defer f.Close()
		s := time.Now()
		Repack(4, 1)
		b.ReportMetric(float64(time.Since(s).Seconds()), "sec")
	}
}

func BenchmarkRepack_8Threads(b *testing.B) {

	for i := 0; i < b.N; i++ {
		f, _ := os.Open("braid.dat")
		defer f.Close()
		s := time.Now()
		Repack(4, 8)
		b.ReportMetric(float64(time.Since(s).Seconds()), "sec")
	}
}

func BenchmarkRepack_16Threads(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f, _ := os.Open("braid.dat")
		defer f.Close()
		s := time.Now()
		Repack(4, 16)
		b.ReportMetric(float64(time.Since(s).Seconds()), "sec")
	}
}

func BenchmarkRepack_32Threads(b *testing.B) {

	for i := 0; i < b.N; i++ {
		f, _ := os.Open("braid.dat")
		defer f.Close()
		s := time.Now()
		Repack(4, 32)
		b.ReportMetric(float64(time.Since(s).Seconds()), "sec")
	}
}

func BenchmarkRepack_64Threads(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f, _ := os.Open("braid.dat")
		defer f.Close()
		s := time.Now()
		Repack(4, 64)
		b.ReportMetric(float64(time.Since(s).Seconds()), "sec")
	}
}

func BenchmarkRepack_100Threads(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f, _ := os.Open("braid.dat")
		defer f.Close()
		s := time.Now()
		Repack(4, 100)
		b.ReportMetric(float64(time.Since(s).Seconds()), "sec")
	}
}
