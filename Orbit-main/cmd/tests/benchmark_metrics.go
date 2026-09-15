package main

import (
	"fmt"
	"runtime"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

type BenchmarkMetrics struct {
	Success       int64
	SoldOut       int64
	AlreadyBooked int64
	Errors        int64

	Start time.Time
	End   time.Time

	LoadGeneratorCPUStart ProcessCPU
	LoadGeneratorCPUEnd   ProcessCPU
	LogicalCPUs           int

	mu           sync.Mutex
	Latencies    []time.Duration
	ErrorSamples []string
}

var Metrics BenchmarkMetrics

func BenchmarkStart() {
	Metrics = BenchmarkMetrics{}
	Metrics.LoadGeneratorCPUStart = currentProcessCPU()
	Metrics.LogicalCPUs = logicalCPUCount()
	Metrics.Start = time.Now()
}

func BenchmarkFinish() {
	Metrics.End = time.Now()
	Metrics.LoadGeneratorCPUEnd = currentProcessCPU()
}

func AddLatency(d time.Duration) {
	Metrics.mu.Lock()
	Metrics.Latencies = append(Metrics.Latencies, d)
	Metrics.mu.Unlock()
}

func IncSuccess() {
	atomic.AddInt64(&Metrics.Success, 1)
}

func IncSoldOut() {
	atomic.AddInt64(&Metrics.SoldOut, 1)
}

func IncAlreadyBooked() {
	atomic.AddInt64(&Metrics.AlreadyBooked, 1)
}

func IncError() {
	atomic.AddInt64(&Metrics.Errors, 1)
}

func AddErrorSample(message string) {
	Metrics.mu.Lock()
	defer Metrics.mu.Unlock()
	if len(Metrics.ErrorSamples) < 5 {
		Metrics.ErrorSamples = append(Metrics.ErrorSamples, message)
	}
}

func percentile(sorted []time.Duration, p float64) time.Duration {
	if len(sorted) == 0 {
		return 0
	}
	idx := int(p * float64(len(sorted)))
	if idx >= len(sorted) {
		idx = len(sorted) - 1
	}
	return sorted[idx]
}

func PrintBenchmark() {

	success := atomic.LoadInt64(&Metrics.Success)
	soldOut := atomic.LoadInt64(&Metrics.SoldOut)
	alreadyBooked := atomic.LoadInt64(&Metrics.AlreadyBooked)
	errCount := atomic.LoadInt64(&Metrics.Errors)

	total := success + soldOut + alreadyBooked + errCount

	duration := Metrics.End.Sub(Metrics.Start)

	var avgLatency, p95Latency, p99Latency time.Duration

	Metrics.mu.Lock()
	latencies := make([]time.Duration, len(Metrics.Latencies))
	copy(latencies, Metrics.Latencies)
	Metrics.mu.Unlock()

	if len(latencies) > 0 {
		var sum time.Duration
		for _, d := range latencies {
			sum += d
		}
		avgLatency = sum / time.Duration(len(latencies))

		sort.Slice(latencies, func(i, j int) bool { return latencies[i] < latencies[j] })
		p95Latency = percentile(latencies, 0.95)
		p99Latency = percentile(latencies, 0.99)
	}

	var throughput float64
	if duration.Seconds() > 0 {
		throughput = float64(total) / duration.Seconds()
	}

	var successPct, errorPct float64
	if total > 0 {
		successPct = float64(success) / float64(total) * 100
		errorPct = float64(errCount) / float64(total) * 100
	}

	loadGeneratorCPU := (Metrics.LoadGeneratorCPUEnd.User - Metrics.LoadGeneratorCPUStart.User) +
		(Metrics.LoadGeneratorCPUEnd.System - Metrics.LoadGeneratorCPUStart.System)
	loadGeneratorCores := 0.0
	loadGeneratorMachinePct := 0.0
	if duration > 0 {
		loadGeneratorCores = float64(loadGeneratorCPU) / float64(duration)
		if Metrics.LogicalCPUs > 0 {
			loadGeneratorMachinePct = loadGeneratorCores / float64(Metrics.LogicalCPUs) * 100
		}
	}

	fmt.Println()
	fmt.Println("========== DEBUG ==========")
	fmt.Println("Buyer Count :", BuyerCount)
	fmt.Println("Total Count :", success+soldOut+alreadyBooked+errCount)
	fmt.Println("===========================")

	fmt.Println()
	fmt.Println("========================================")
	fmt.Println("BENCHMARK REPORT")
	fmt.Println("========================================")
	fmt.Printf("Concurrent Buyers : %d\n", BuyerCount)
	fmt.Printf("Total Requests    : %d\n", total)
	fmt.Printf("Success           : %d (%.2f%%)\n", success, successPct)
	fmt.Printf("Sold Out          : %d\n", soldOut)
	fmt.Printf("Already Booked    : %d\n", alreadyBooked)
	fmt.Printf("Errors            : %d (%.2f%%)\n", errCount, errorPct)
	fmt.Println("----------------------------------------")
	fmt.Printf("Duration          : %v\n", duration)
	fmt.Printf("Throughput        : %.2f req/sec\n", throughput)
	fmt.Printf("Average Latency   : %v\n", avgLatency)
	fmt.Printf("p95 Latency       : %v\n", p95Latency)
	fmt.Printf("p99 Latency       : %v\n", p99Latency)
	fmt.Println("----------------------------------------")
	fmt.Printf("Load Generator CPU: %v (%.2f core avg)\n", loadGeneratorCPU, loadGeneratorCores)
	fmt.Printf("Generator CPU Share: %.2f%% of %d logical CPUs\n", loadGeneratorMachinePct, Metrics.LogicalCPUs)
	fmt.Printf("Go Routines (end) : %d\n", runtime.NumGoroutine())
	Metrics.mu.Lock()
	errorSamples := append([]string(nil), Metrics.ErrorSamples...)
	Metrics.mu.Unlock()
	if len(errorSamples) > 0 {
		fmt.Println("Unexpected response samples:")
		for _, sample := range errorSamples {
			fmt.Printf("  - %s\n", sample)
		}
	}
	fmt.Println("========================================")
}
