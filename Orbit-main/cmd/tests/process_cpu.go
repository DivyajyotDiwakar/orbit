package main

import (
	"runtime"
	"syscall"
	"time"
)

type ProcessCPU struct {
	User   time.Duration
	System time.Duration
}

func currentProcessCPU() ProcessCPU {
	var usage syscall.Rusage
	if err := syscall.Getrusage(syscall.RUSAGE_SELF, &usage); err != nil {
		return ProcessCPU{}
	}

	return ProcessCPU{
		User:   timevalDuration(usage.Utime),
		System: timevalDuration(usage.Stime),
	}
}

func timevalDuration(value syscall.Timeval) time.Duration {
	return time.Duration(value.Sec)*time.Second + time.Duration(value.Usec)*time.Microsecond
}

func logicalCPUCount() int {
	return runtime.NumCPU()
}
