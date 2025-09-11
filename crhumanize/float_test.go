// Copyright 2025 The Cockroach Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.

package crhumanize

import (
	"fmt"
	"testing"
)

func TestFloat(t *testing.T) {
	tests := []struct {
		value         float64
		decimalDigits int
		expected      string
	}{
		{0, 0, "0"},
		{0, 1, "0"},
		{0, 1, "0"},
		{0.1, 1, "0.1"},
		{0.1, 2, "0.1"},
		{0.01, 1, "0"},
		{0.01, 2, "0.01"},
		{0.01, 4, "0.01"},
		{-1.23456789, 2, "-1.23"},
		{1.23456789, 2, "1.23"},
		{1.23456789, 3, "1.235"},
		{1.23456789, 3, "1.235"},
		{123456.7777, 1, "123456.8"},
		{123456.7777, 2, "123456.78"},
		{123456.1010, 4, "123456.101"},
		{123456.1010, 2, "123456.1"},
		{-123456.1010, 1, "-123456.1"},
		{123456.1010, 1, "123456.1"},
		{-123456.1010, 0, "-123456"},
		{123456.1010, 0, "123456"},
	}

	for _, test := range tests {
		result := string(Float(test.value, test.decimalDigits))
		if result != test.expected {
			t.Errorf("Float(%f, %d) = %s; expected %s", test.value, test.decimalDigits, result, test.expected)
		}
	}
}

func TestPercent(t *testing.T) {
	tests := []struct {
		a, b     float64
		expected string
	}{
		{a: 0, b: 0, expected: ""},
		{a: 0, b: 100, expected: "0%"},
		{a: 0.0001, b: 100.0, expected: "~0%"},
		{a: 0.044, b: 100, expected: "~0%"},
		{a: 0.05, b: 100, expected: "0.1%"},
		{a: 0.1234, b: 100.0, expected: "0.1%"},
		{a: -0.1234, b: 100.0, expected: "-0.1%"},
		{a: 0.05, b: 100.0, expected: "0.1%"},
		{a: 9.95, b: 100.0, expected: "10%"},
		{a: 9.94, b: 100.0, expected: "9.9%"},
		{a: -9.95, b: 100.0, expected: "-10%"},
		{a: -9.94, b: 100.0, expected: "-9.9%"},
		{a: 10.52345, b: 100.0, expected: "11%"},
	}

	for _, test := range tests {
		result := string(Percent(test.a, test.b))
		if result != test.expected {
			t.Errorf("Percent(%f,%f) = %s; expected %s", test.a, test.b, result, test.expected)
		}
	}
}

func ExampleFloat() {
	fmt.Println(Float(100.1234, 3))
	fmt.Println(Float(100.12, 3))
	fmt.Println(Float(100.1, 3))
	fmt.Println(Float(100, 3))
	// Output:
	// 100.123
	// 100.12
	// 100.1
	// 100
}

func ExamplePercent() {
	fmt.Println(Percent(uint64(0), uint64(10000)))
	fmt.Println(Percent(uint64(1), uint64(10000)))
	fmt.Println(Percent(uint64(12), uint64(10000)))
	fmt.Println(Percent(uint64(123), uint64(10000)))
	fmt.Println(Percent(uint64(1234), uint64(10000)))
	// Output:
	// 0%
	// ~0%
	// 0.1%
	// 1.2%
	// 12%
}
