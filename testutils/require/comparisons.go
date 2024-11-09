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

package require

// LT asserts that a < b.
func LT[T ordered](tb TB, a, b T) {
	if !(a < b) {
		tb.Helper()
		tb.Fatalf("expected %v < %v", a, b)
	}
}

// LE asserts that a <= b.
func LE[T ordered](tb TB, a, b T) {
	if !(a <= b) {
		tb.Helper()
		tb.Fatalf("expected %v <= %v", a, b)
	}
}

// GT asserts that a > b.
func GT[T ordered](tb TB, a, b T) {
	if !(a > b) {
		tb.Helper()
		tb.Fatalf("expected %v > %v", a, b)
	}
}

// GE asserts that a >= b.
func GE[T ordered](tb TB, a, b T) {
	if !(a >= b) {
		tb.Helper()
		tb.Fatalf("expected %v >= %v", a, b)
	}
}

type ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr | ~float32 | ~float64 | ~string
}
