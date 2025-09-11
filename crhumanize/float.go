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
	"strconv"
	"strings"
)

// Float formats the given float with the specified number of decimal digits.
// Trailing 0 decimals are stripped.
func Float(value float64, decimalDigits int) SafeString {
	s := strconv.FormatFloat(value, 'f', decimalDigits, 64)
	s = stripTrailingZeroDecimals(s)
	return SafeString(s)
}

// CompactFloat formats the given float with at most one decimal digit.
// Specifically: we show a decimal digit only when the integer part is a single
// digit.
func CompactFloat(value float64) SafeString {
	decimalDigits := 0
	if math.Abs(value) < 9.95 {
		decimalDigits = 1
	}
	return Float(value, decimalDigits)
}

func stripTrailingZeroDecimals(s string) string {
	if !strings.ContainsRune(s, '.') {
		return s
	}
	for s[len(s)-1] == '0' {
		s = s[:len(s)-1]
	}
	if s[len(s)-1] == '.' {
		s = s[:len(s)-1]
	}
	return s
}
