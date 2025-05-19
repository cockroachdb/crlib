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
	"math"
	"math/bits"
	"strings"
)

// Integer is a constraint that permits any integer type.
type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// SafeString represents a human readable representation of a value. It
// implements a `SafeValue()` marker method (implementing the
// github.com/cockroachdb/redact.SafeValue interface) to signal that it
// represents a string that does not need to be redacted.
type SafeString string

// SafeValue implements cockroachdb/redact.SafeValue.
func (fs SafeString) SafeValue() {}

// String implements fmt.Stringer.
func (fs SafeString) String() string { return string(fs) }

var iecUnits = []string{"", "Ki", "Mi", "Gi", "Ti", "Pi", "Ei"}
var siUnits = []string{"", "K", "M", "G", "T", "P", "E"}

func iecUnit(value uint64) (index int, scaled float64) {
	n := (max(0, bits.Len64(value)-1)) / 10
	return n, float64(value) / float64(uint64(1)<<(10*n))
}

func siUnit(value uint64) (index int, scaled float64) {
	if value < 10 {
		return 0, float64(value)
	}
	n := int(math.Floor(math.Log10(float64(value)))) / 3
	return n, float64(value) / math.Pow10(3*n)
}

func parseUnit(s string) (index int, iec bool, ok bool) {
	s = strings.ToUpper(s)
	s, iec = strings.CutSuffix(s, "I")
	for i := range siUnits {
		if s == siUnits[i] {
			if i == 0 && iec {
				// Just an "i" is not ok.
				return 0, false, false
			}
			return i, iec, true
		}
	}
	return 0, false, false
}
