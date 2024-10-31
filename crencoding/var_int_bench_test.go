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
	"fmt"
	"io"
	"math/rand/v2"
	"testing"
)

// Sample benchmark results:
//
// Apple M1 (arm64/darwin):
// name                                           time/op
// UvarintLen32/range=100/simple-10                 0.64ns ± 1%
// UvarintLen32/range=100/crlib-10                  0.86ns ± 2%
// UvarintLen32/range=1000/simple-10                1.11ns ± 1%
// UvarintLen32/range=1000/crlib-10                 0.90ns ±12%
// UvarintLen32/range=100000/simple-10              1.89ns ± 1%
// UvarintLen32/range=100000/crlib-10               0.84ns ± 0%
// UvarintLen32/range=1000000000/simple-10          3.05ns ± 1%
// UvarintLen32/range=1000000000/crlib-10           0.84ns ± 0%
// UvarintLen32/range=4000000000/simple-10          2.83ns ± 2%
// UvarintLen32/range=4000000000/crlib-10           0.85ns ± 1%
//
// UvarintLen64/range=100/simple-10                 0.63ns ± 1%
// UvarintLen64/range=100/crlib-10                  0.84ns ± 0%
// UvarintLen64/range=1000/simple-10                1.10ns ± 1%
// UvarintLen64/range=1000/crlib-10                 0.84ns ± 0%
// UvarintLen64/range=100000/simple-10              1.94ns ± 1%
// UvarintLen64/range=100000/crlib-10               0.84ns ± 0%
// UvarintLen64/range=1000000000/simple-10          2.83ns ± 3%
// UvarintLen64/range=1000000000/crlib-10           0.84ns ± 0%
// UvarintLen64/range=1000000000000/simple-10       3.02ns ± 0%
// UvarintLen64/range=1000000000000/crlib-10        0.84ns ± 0%
// UvarintLen64/range=1000000000000000/simple-10    4.88ns ± 1%
// UvarintLen64/range=1000000000000000/crlib-10     0.84ns ± 0%
//
// UvarintLenSum32/loop-10                          0.67ns ± 0%
// UvarintLenSum32/crlib-10                         0.49ns ± 0%
//
// UvarintLenSum64/loop-10                          0.67ns ± 0%
// UvarintLenSum64/crlib-10                         0.49ns ± 0%
//
// name                                           speed
// UvarintLenSum32/loop-10                        6.11TB/s ± 0%
// UvarintLenSum32/crlib-10                       8.41TB/s ± 0%
//
// UvarintLenSum64/loop-10                        6.12TB/s ± 0%
// UvarintLenSum64/crlib-10                       8.41TB/s ± 0%

func BenchmarkUvarintLen32(b *testing.B) {
	for _, valRange := range []uint32{100, 1000, 100_000, 1_000_000_000, 4_000_000_000} {
		b.Run(fmt.Sprintf("range=%d", valRange), func(b *testing.B) {
			const numValues = 1024
			values := make([]uint32, numValues)
			for i := range values {
				values[i] = rand.Uint32N(valRange)
			}

			b.Run("simple", func(b *testing.B) {
				var x int
				for i := 0; i < b.N; i++ {
					x ^= simpleVarUint32Len(values[i&(numValues-1)])
				}
				fmt.Fprint(io.Discard, x)
			})

			b.Run("crlib", func(b *testing.B) {
				var x int
				for i := 0; i < b.N; i++ {
					x ^= UvarintLen32(values[i&(numValues-1)])
				}
				fmt.Fprint(io.Discard, x)
			})
		})
	}
}

func BenchmarkUvarintLen64(b *testing.B) {
	for _, valRange := range []uint64{100, 1000, 100_000, 1_000_000_000, 1_000_000_000_000, 1_000_000_000_000_000} {
		b.Run(fmt.Sprintf("range=%d", valRange), func(b *testing.B) {
			const numValues = 1024
			values := make([]uint64, numValues)
			for i := range values {
				values[i] = rand.Uint64N(valRange)
			}

			b.Run("simple", func(b *testing.B) {
				var x int
				for i := 0; i < b.N; i++ {
					x ^= simpleVarUint64Len(values[i&(numValues-1)])
				}
				fmt.Fprint(io.Discard, x)
			})

			b.Run("crlib", func(b *testing.B) {
				var x int
				for i := 0; i < b.N; i++ {
					x ^= UvarintLen64(values[i&(numValues-1)])
				}
				fmt.Fprint(io.Discard, x)
			})
		})
	}
}

func BenchmarkUvarintLenSum32(b *testing.B) {
	const numValues = 1024
	values := make([]uint32, numValues)
	for i := range values {
		values[i] = rand.Uint32()
	}

	b.Run("loop", func(b *testing.B) {
		b.SetBytes(numValues * 4)
		var x int
		n := b.N / numValues
		for i := 0; i < n; i++ {
			var s int
			for _, v := range values {
				s += UvarintLen32(v)
			}
			x ^= s
		}
		fmt.Fprint(io.Discard, x)
	})

	b.Run("crlib", func(b *testing.B) {
		b.SetBytes(numValues * 4)
		var x int
		n := b.N / numValues
		for i := 0; i < n; i++ {
			x ^= UvarintLenSum32(values)
		}
		fmt.Fprint(io.Discard, x)
	})
}

func BenchmarkUvarintLenSum64(b *testing.B) {
	const numValues = 1024
	values := make([]uint64, numValues)
	for i := range values {
		values[i] = rand.Uint64()
	}

	b.Run("loop", func(b *testing.B) {
		b.SetBytes(numValues * 4)
		var x int
		n := b.N / numValues
		for i := 0; i < n; i++ {
			var s int
			for _, v := range values {
				s += UvarintLen64(v)
			}
			x ^= s
		}
		fmt.Fprint(io.Discard, x)
	})

	b.Run("crlib", func(b *testing.B) {
		b.SetBytes(numValues * 4)
		var x int
		n := b.N / numValues
		for i := 0; i < n; i++ {
			x ^= UvarintLenSum64(values)
		}
		fmt.Fprint(io.Discard, x)
	})
}

func simpleVarUint32Len(n uint32) int {
	r := 1
	for n > 0x80 {
		r++
		n >>= 7
	}
	return r
}

func simpleVarUint64Len(n uint64) int {
	r := 1
	for n > 0x80 {
		r++
		n >>= 7
	}
	return r
}
