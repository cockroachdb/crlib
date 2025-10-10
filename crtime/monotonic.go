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

package crtime

import (
	"time"

	"github.com/cockroachdb/crlib/crsync"
)

// Mono represents a moment in time in terms of a monotonic clock. Its value is
// the duration since the start of the process.
//
// Note that if the system doesn't support a monotonic clock, the wall clock is
// used.
type Mono time.Duration

type ClockReading[T any] interface {
	// Add returns a ClockReading time.Duration in the future from a given
	// ClockReading.
	Add(time.Duration) T
	// Sub returns the duration elapsed between two clock readings.
	Sub(T) time.Duration
}

// NowMono returns a moment in time in terms of a monotonic clock. It is faster
// than time.Now which also consults the wall clock.
func NowMono() Mono {
	// Note: time.Since reads only the monotonic clock (if it is available).
	return Mono(time.Since(startTime))
}

// Sub returns the duration that elapsed between two moments.
func (m Mono) Sub(other Mono) time.Duration {
	return time.Duration(m - other)
}

// Add returns the time.Mono m+d.
func (m Mono) Add(d time.Duration) Mono {
	return Mono(time.Duration(m) + d)
}

// Elapsed returns the duration that elapsed since m.
func (m Mono) Elapsed() time.Duration {
	return time.Duration(NowMono() - m)
}

// ToUTC returns the UTC time corresponding to the monotonic time.
//
// The time is derived from the current wall clock, adjusted by the difference
// in the monotonic clock values. Note that if the wall clock has been changed
// since the Mono value was obtained, the result does not reflect the wall clock
// at that point in time.
func (m Mono) ToUTC() time.Time {
	now := time.Now()
	adjustment := time.Duration(m) - now.Sub(startTime)
	return now.UTC().Add(adjustment)
}

// MonoFromTime converts a time.Time to a Mono value. If the time has a
// monotonic component, it is used.
func MonoFromTime(t time.Time) Mono {
	return Mono(t.Sub(startTime))
}

// AtomicMono provides atomic access to a Mono value.
type AtomicMono = crsync.TypedAtomicInt64[Mono]

// We use startTime as a reference point against which we can call
// time.Since(). This solution is suggested by the Go runtime code:
// https://github.com/golang/go/blob/889abb17e125bb0f5d8de61bb80ef15fbe2a130d/src/runtime/time_nofake.go#L19
var startTime = time.Now()
