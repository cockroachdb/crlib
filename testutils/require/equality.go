// Copyright 2024 The Cockroach Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.

package require

import (
	"fmt"
	"reflect"
)

// Equal asserts that a and b are deeply equal.
func Equal[T any](tb TB, a, b T) {
	if !reflect.DeepEqual(a, b) {
		tb.Helper()
		aStr := fmt.Sprint(a)
		bStr := fmt.Sprint(b)
		if len(aStr)+len(bStr) > 80 {
			tb.Fatalf("expected equality:\n  a: %s\n  b: %s", aStr, bStr)
		} else {
			tb.Fatalf("expected %s == %s", aStr, bStr)
		}
	}
}

// NotEqual asserts that a and b are deeply equal.
func NotEqual[T any](tb TB, a, b T) {
	if reflect.DeepEqual(a, b) {
		tb.Helper()
		tb.Fatalf("expected %v != %v", a, b)
	}
}

// True asserts that the value is true.
func True[T ~bool](tb TB, a T) {
	if !a {
		tb.Helper()
		tb.Fatalf("expected true")
	}
}

// False asserts that the value is false.
func False[T ~bool](tb TB, a T) {
	if a {
		tb.Helper()
		tb.Fatalf("expected false")
	}
}
