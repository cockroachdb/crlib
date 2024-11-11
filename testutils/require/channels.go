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

import "time"

// Recv asserts that a value is received on the channel within the specified
// duration within 1 second and returns that value.
func Recv[T any](tb TB, ch chan T) T {
	select {
	case v := <-ch:
		return v
	case <-time.After(1 * time.Second):
		tb.Helper()
		tb.Fatal("did not receive on channel")
		panic("unreachable")
	}
}

// RecvWithin asserts that a value is received on the channel within the specified
// duration, and returns that value.
func RecvWithin[T any](tb TB, ch chan T, within time.Duration) T {
	select {
	case v := <-ch:
		return v
	case <-time.After(within):
		tb.Helper()
		tb.Fatal("did not receive on channel")
		panic("unreachable")
	}
}

// NoRecv asserts that no value is received on the channel within 10ms.
func NoRecv[T any](tb TB, ch chan T) {
	select {
	case <-ch:
		tb.Helper()
		tb.Fatal("received unexpected notification")
	case <-time.After(10 * time.Millisecond):
	}
}

// NoRecvWithin asserts that no value is received on the channel within the
// specified duration.
func NoRecvWithin[T any](tb TB, ch chan T, within time.Duration) {
	select {
	case <-ch:
		tb.Helper()
		tb.Fatal("received unexpected notification")
	case <-time.After(within):
	}
}
