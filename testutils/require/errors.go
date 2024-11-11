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

import "fmt"

// NoError asserts that err is nil.
func NoError(tb TB, err error) {
	if err != nil {
		tb.Helper()
		tb.Fatalf("unexpected error: %v", err)
	}
}

// NoError1 is passed an arbitrary value and an error and panics if the error is
// not-nil, otherwise returns the value. It can be used to get the return value
// of a fallible function that must succeed.
//
// Instead of:
//
//	v, err := SomeFunc()
//	if err != nil {
//	  t.Fatal(err)
//	}
//
// We can use:
//
//	v := require.NoError1(SomeFunc())
func NoError1[T any](a T, err error) T {
	if err != nil {
		panic(fmt.Sprintf("unexpected error: %+v", err))
	}
	return a
}

// NoError2 is passed two arbitrary values and an error and panics if the error
// is not-nil, otherwise returns the values. It can be used to get the return
// values of a fallible function that must succeed.
//
// Instead of:
//
//	v, w, err := SomeFunc()
//	if err != nil {
//	  t.Fatal(err)
//	}
//
// We can use:
//
//	v, w := require.NoError2(SomeFunc())
func NoError2[T any, U any](a T, b U, err error) (T, U) {
	if err != nil {
		panic(fmt.Sprintf("unexpected error: %+v", err))
	}
	return a, b
}
