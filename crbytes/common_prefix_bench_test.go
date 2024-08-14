// Copyright 2024 The Cockroach Authors.
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

package crbytes

import (
	"bytes"
	"math/rand"
	"slices"
	"testing"
)

// Sample benchmark results:
//
// linux/amd64, Intel(R) Xeon(R) CPU @ 2.80GHz:
//
//	CommonPrefix/small/crbytes-24   8.54ns ± 1%
//	CommonPrefix/small/generic-24   11.0ns ± 1%
//	CommonPrefix/small/naive-24     13.6ns ± 1%
//	CommonPrefix/medium/crbytes-24  13.8ns ± 2%
//	CommonPrefix/medium/generic-24  26.3ns ± 2%
//	CommonPrefix/medium/naive-24    31.7ns ± 2%
//	CommonPrefix/large/crbytes-24    153ns ± 2%
//	CommonPrefix/large/generic-24    362ns ± 2%
//	CommonPrefix/large/naive-24      755ns ± 1%
//
// darwin/arm64, Apple M1:
//
//	CommonPrefix/small/crbytes-10   5.48ns ± 1%
//	CommonPrefix/small/generic-10   7.02ns ± 7%
//	CommonPrefix/small/naive-10     9.58ns ± 2%
//	CommonPrefix/medium/crbytes-10  7.42ns ± 8%
//	CommonPrefix/medium/generic-10  15.6ns ± 5%
//	CommonPrefix/medium/naive-10    23.5ns ± 7%
//	CommonPrefix/large/crbytes-10    125ns ± 4%
//	CommonPrefix/large/generic-10    249ns ±11%
//	CommonPrefix/large/naive-10      698ns ± 0%
func BenchmarkCommonPrefix(b *testing.B) {
	small := lexicographicSet(4, 16)
	medium := lexicographicSet(10, 100)
	large := lexicographicSet(1000, 10000)
	b.Run("small", func(b *testing.B) {
		runBenchComparison(b, small)
	})
	b.Run("medium", func(b *testing.B) {
		runBenchComparison(b, medium)
	})
	b.Run("large", func(b *testing.B) {
		runBenchComparison(b, large)
	})
}

func runBenchComparison(b *testing.B, input [][]byte) {
	b.Run("crbytes", func(b *testing.B) {
		runBench(b, input, CommonPrefix)
	})
	b.Run("generic", func(b *testing.B) {
		runBench(b, input, commonPrefixGeneric)
	})
	b.Run("naive", func(b *testing.B) {
		runBench(b, input, commonPrefixNaive)
	})
}

func runBench(b *testing.B, input [][]byte, impl func(a, b []byte) int) {
	n := len(input)
	j := 0
	var sum int
	for i := 0; i < b.N; i++ {
		next := j + 1
		if next >= n {
			next = 0
		}
		sum += impl(input[j], input[next])
		j = next
	}
	b.Logf("average result: %d\n", sum/b.N)
}

// lexicographicSet returns a lexicographically ordered list of byte slices
// which all have a common prefix of length minLength, with random bytes (with
// alphabet size 2) following up to maxLength.
func lexicographicSet(minLength, maxLength int) [][]byte {
	const n = 10_000
	const alphabet = 2
	prefix := genBytes(minLength, alphabet)

	result := make([][]byte, n)
	for i := range result {
		result[i] = slices.Concat(prefix, genBytes(rand.Intn(maxLength-minLength+1), alphabet))
	}
	slices.SortFunc(result, bytes.Compare)
	return result
}
