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

package crmath

import (
	"math"
	"math/big"
	"math/rand/v2"
	"testing"
)

func TestScaleUint64(t *testing.T) {
	expect := func(x, a, b, expected uint64) {
		if got := ScaleUint64(x, a, b); got != expected {
			t.Helper()
			t.Fatalf("ScaleUint64(%d, %d, %d) = %d; want %d", x, a, b, got, expected)
		}
	}

	expect(0, 1, 1, 0)

	expect(1, 1, 123456, 1)
	expect(1, 1, math.MaxUint64, 1)

	expect(1000, 1, 10, 100)

	expect(1<<52, 1<<40, 1<<60, 1<<32)
	expect(1<<52, 1<<60, 1<<50, 1<<62)
	expect(1<<52, 1<<60, 1<<40, math.MaxUint64)

	for range 100 {
		n := rand.Uint64()
		expect(n, 1, 1, n)
		expect(n, math.MaxUint64, math.MaxUint64, n)
		m := rand.Uint64()
		expect(n, m, m, n)
	}

	// calc is an alternative implementation using big.Int.
	calc := func(xx, aa, bb uint64) uint64 {
		var x, a, b, numerator, res big.Int
		x.SetUint64(xx)
		a.SetUint64(aa)
		b.SetUint64(bb)
		numerator.Mul(&a, &x)
		numerator.Add(&numerator, &b)
		numerator.Add(&numerator, big.NewInt(-1))
		res.Quo(&numerator, &b)
		if res.BitLen() > 64 {
			return math.MaxUint64
		}
		return res.Uint64()
	}

	checkDivByZero := func(x, a uint64) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatalf("ScaleUint64(%d, %d, 0) did not panic", x, a)
			}
		}()
		_ = ScaleUint64(x, a, 0)
	}

	for range 1000 {
		r := func() uint64 {
			if rand.IntN(2) == 0 {
				// Return values close to the fast path cutoff.
				return math.MaxUint32 + rand.Uint64N(5) - 2
			}
			return rand.Uint64()
		}
		x, a, b := r(), r(), r()
		expect(x, a, b, calc(x, a, b))

		checkDivByZero(x, a)
	}

}
