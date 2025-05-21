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
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Bytes returns an approximate (within 5%) human-readable representations of a
// byte value in IEC units.
//
// Examples: "1.5 MiB", "21 GiB", "3 B".
func Bytes[T Integer](bytes T) SafeString {
	if bytes < 0 {
		// Note: uint64(-bytes) doesn't work correctly when bytes is the minimum
		// value for a smaller type.
		return "-" + bytesUint64(-uint64(bytes), false)
	}
	return bytesUint64(uint64(bytes), false)
}

// BytesCompact returns an approximate (within 5%) human-readable representations of a
// byte value. It is similar to Bytes but omits the space and the "i" in the
// units. The units are still base-1024.
// Examples: "1.5MB", "21GB", "3B".
func BytesCompact[T Integer](bytes T) SafeString {
	if bytes < 0 {
		return "-" + bytesUint64(-uint64(bytes), true)
	}
	return bytesUint64(uint64(bytes), true)
}

func bytesUint64(bytes uint64, compact bool) SafeString {
	n, scaled := iecUnit(bytes)
	digits := 0
	if scaled < 10 {
		digits = 1
	}
	if compact {
		return SafeString(fmt.Sprintf("%s%sB", Float(scaled, digits), siUnits[n]))
	}
	return SafeString(fmt.Sprintf("%s %sB", Float(scaled, digits), iecUnits[n]))
}

// BytesExact is similar to Bytes, but the result is exact and can be parsed
// back into the original value. It separates groups of digits in large numbers
// with commas for readability.
//
// It is guaranteed that ParseBytes[T](BytesExact[T](x)) == x for all x.
//
// An example when this should be used instead of Bytes is when we are
// marshaling a configuration value.
//
// Examples: "1,234 KiB", "21,000 GiB", "1,000,000 B".
func BytesExact[T Integer](bytes T) SafeString {
	if bytes < 0 {
		return "-" + bytesExactUint64(-uint64(bytes))
	}
	return bytesExactUint64(uint64(bytes))
}

func bytesExactUint64(bytes uint64) SafeString {
	i := 0
	if bytes != 0 {
		for ; i < len(iecUnits)-1 && bytes%1024 == 0; i++ {
			bytes /= 1024
		}
	}
	valStr := strconv.FormatUint(bytes, 10)
	var buf strings.Builder
	buf.Grow(len(valStr)*4/3 + len(iecUnits[i]) + 2)

	// Add commas to make the number more readable.
	n := 1 + (len(valStr)-1)%3 // length of the first digit group.
	buf.WriteString(valStr[:n])
	for i := n; i < len(valStr); i += 3 {
		buf.WriteByte(',')
		buf.WriteString(valStr[i : i+3])
	}
	buf.WriteByte(' ')
	buf.WriteString(iecUnits[i])
	buf.WriteByte('B')

	return SafeString(buf.String())
}

func ParseBytes[T Integer](s string) (T, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("cannot parse bytes from %q", s)
	}
	unsignedType := T(0)-T(1) > T(0)
	minusSign := false
	unsignedPart := s
	if s[0] == '-' {
		// Type is unsigned.
		if unsignedType {
			return T(0), fmt.Errorf("cannot parse non-negative bytes value from %q", s)
		}
		minusSign = true
		unsignedPart = s[1:]
	}
	val, err := parseBytesUint64(unsignedPart)
	if err != nil {
		return 0, err
	}

	// Apply negation and convert to T, checking for numeric overflow.
	result, ok := func() (T, bool) {
		if minusSign {
			if val > -math.MinInt64 {
				return T(0), false
			}
			x := -int64(val)
			result := T(x)
			if int64(result) != x {
				return T(0), false
			}
			return result, true
		}

		result := T(val)
		if unsignedType {
			if uint64(result) != val {
				return T(0), false
			}
		} else if val > math.MaxInt64 || result < 0 || int64(result) != int64(val) {
			return T(0), false
		}
		return result, true
	}()

	if !ok {
		return T(0), fmt.Errorf("cannot parse bytes value from %q (numeric overflow)", s)
	}
	return result, nil
}

func parseBytesUint64(s string) (uint64, error) {
	numStr := s
	for i, r := range s {
		if (r >= '0' && r <= '9') || r == '.' || r == ',' {
			continue
		}
		numStr = s[:i]
		break
	}
	suffix := strings.TrimSpace(s[len(numStr):])
	suffix = strings.ToUpper(suffix)
	// Tolerate but don't require ending with B.
	suffix = strings.TrimSuffix(suffix, "B")
	unitIdx, _, ok := parseUnit(suffix)
	if !ok {
		return 0, fmt.Errorf("cannot parse bytes from %q", s)
	}
	scale := uint64(1) << (10 * unitIdx)
	numStr = strings.ReplaceAll(numStr, ",", "")

	// We want to guarantee exact parsing of integer values, even those too large
	// to be accurately represented in a float64.
	if !strings.Contains(numStr, ".") {
		value, err := strconv.ParseUint(numStr, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("cannot parse bytes from %q", s)
		}
		if value != 0 && scale > math.MaxUint64/value {
			return 0, fmt.Errorf("cannot parse bytes from %q (numeric overflow)", s)
		}
		return scale * value, nil
	}
	value, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return 0, fmt.Errorf("cannot parse bytes from %q", s)
	}
	value = math.Round(value * float64(scale))
	if value > math.MaxUint64 {
		return 0, fmt.Errorf("cannot parse bytes from %q (numeric overflow)", s)
	}
	return uint64(value), nil
}
