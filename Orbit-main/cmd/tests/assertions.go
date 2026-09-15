package main

import (
	"fmt"
	"os"
)

var (
	totalTests  int
	passedTests int
)

func Pass(name string) {
	totalTests++
	passedTests++
	fmt.Printf("PASS : %s\n", name)
}

func Fail(name string, reason string) {
	totalTests++
	fmt.Printf("FAIL : %s\n", name)

	if reason != "" {
		fmt.Printf("   ↳ %s\n", reason)
	}

	fmt.Println()
	os.Exit(1)
}

func Assert(err error) {
	if err != nil {
		Fail("Unexpected Error", err.Error())
	}
}

func AssertStatus(expected int, actual int) {
	if expected != actual {
		Fail(
			"HTTP Status",
			fmt.Sprintf("expected %d got %d", expected, actual),
		)
	}
}

func AssertTrue(ok bool, testName string, reason string) {
	if !ok {
		Fail(testName, reason)
	}

	Pass(testName)
}

func AssertEqual[T comparable](expected T, actual T, testName string) {

	if expected != actual {
		Fail(
			testName,
			fmt.Sprintf("expected=%v actual=%v", expected, actual),
		)
	}

	Pass(testName)
}

func PrintSummary() {

	fmt.Println()
	fmt.Println("========================================")
	fmt.Println("Integration Test Summary")
	fmt.Println("========================================")

	fmt.Printf("Passed : %d/%d\n", passedTests, totalTests)

	if passedTests == totalTests {
		fmt.Println("ALL TESTS PASSED")
	} else {
		fmt.Println("SOME TESTS FAILED")
	}

	fmt.Println("========================================")
}
