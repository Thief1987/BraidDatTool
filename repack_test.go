package main

import (
	"testing"
	"time"
)

func BenchmarkRepack_1Thread(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := time.Now()
		arcName := "braid.dat_new"
		Repack(4, 1, arcName, false)
		b.ReportMetric(float64(time.Since(s).Seconds()), "sec")
	}
}

func BenchmarkRepack_8Threads(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := time.Now()
		arcName := "braid.dat_new"
		Repack(4, 8, arcName, false)
		b.ReportMetric(float64(time.Since(s).Seconds()), "sec")
	}
}

func BenchmarkRepack_16Threads(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := time.Now()
		arcName := "braid.dat_new"
		Repack(4, 16, arcName, false)
		b.ReportMetric(float64(time.Since(s).Seconds()), "sec")
	}
}

func BenchmarkRepack_32Threads(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := time.Now()
		arcName := "braid.dat_new"
		Repack(4, 32, arcName, false)
		b.ReportMetric(float64(time.Since(s).Seconds()), "sec")
	}
}

func BenchmarkRepack_64Threads(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := time.Now()
		arcName := "braid.dat_new"
		Repack(4, 64, arcName, false)
		b.ReportMetric(float64(time.Since(s).Seconds()), "sec")
	}
}

func BenchmarkRepack_100Threads(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := time.Now()
		arcName := "braid.dat_new"
		Repack(4, 100, arcName, false)
		b.ReportMetric(float64(time.Since(s).Seconds()), "sec")
	}
}
