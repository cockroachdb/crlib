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

package crencoding

import (
	"encoding/binary"
	"math"
	"math/rand/v2"
	"testing"
)

// TestUvarintLen tests UvarintLen32 and UvarintLen64.
func TestUvarintLen(t *testing.T) {
	check := func(n uint64) {
		res64 := UvarintLen64(n)
		if expected := len(binary.AppendUvarint(nil, n)); res64 != expected {
			t.Fatalf("invalid result for %d: %d instead of %d", n, res64, expected)
		}
		res32 := UvarintLen32(uint32(n))
		if expected := len(binary.AppendUvarint(nil, uint64(uint32(n)))); res32 != expected {
			t.Fatalf("invalid result for %d: %d instead of %d", n, res32, expected)
		}
	}
	check(0)
	check(math.MaxUint64)
	for i := uint64(0); i < 64; i++ {
		check(1<<i - 1)
		check(1 << i)
		check(1<<i + 1)
	}
	for i := 0; i < 100000; i++ {
		check(rand.Uint64() >> rand.UintN(64))
	}
}

// TestUvarintLenSum tests UvarintLenSum32 and UvarintLenSum64.
func TestUvarintLenSum(t *testing.T) {
	var lengths []int
	for i := 0; i < 100; i++ {
		lengths = append(lengths, i)
	}
	for i := 0; i < 100; i++ {
		lengths = append(lengths, rand.IntN(1000))
	}
	for _, l := range lengths {
		vals64 := make([]uint64, l)
		vals32 := make([]uint32, l)
		for shift := 0; shift < 64; shift++ {
			var exp32, exp64 int
			for i := range vals64 {
				vals64[i] = rand.Uint64() >> shift
				vals32[i] = uint32(vals64[i])
				exp64 += UvarintLen64(vals64[i])
				exp32 += UvarintLen32(vals32[i])
			}
			if actual := UvarintLenSum64(vals64); actual != exp64 {
				t.Log(vals64)
				t.Fatalf("UvarintLenSum64: expected %d, got %d", exp64, actual)
			}
			if actual := UvarintLenSum32(vals32); actual != exp32 {
				t.Fatalf("UvarintLenSum32: expected %d, got %d", exp32, actual)
			}
		}
	}
}
