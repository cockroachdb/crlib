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
	"math/rand/v2"
	"reflect"
	"strings"
	"testing"
)

func TestBytes(t *testing.T) {
	tests := []struct {
		value    int64
		expected string
	}{
		{0, "0 B"},
		{1, "1 B"},
		{-1, "-1 B"},
		{2, "2 B"},
		{100, "100 B"},
		{-100, "-100 B"},
		{900, "900 B"},
		{1000, "1000 B"},
		{1023, "1023 B"},
		{1024, "1 KiB"},
		{1090, "1.1 KiB"},
		{2900, "2.8 KiB"},
		{10100, "9.9 KiB"},
		{10200, "10 KiB"},
		{10240, "10 KiB"},
		{10300, "10 KiB"},
		{10800, "11 KiB"},
		{12 << 20, "12 MiB"},
		{1200 << 20, "1.2 GiB"},
		{12 << 30, "12 GiB"},
		{1200 << 30, "1.2 TiB"},
		{12 << 40, "12 TiB"},
		{1200 << 40, "1.2 PiB"},
		{12 << 50, "12 PiB"},
		{1200 << 50, "1.2 EiB"},
		{2 << 60, "2 EiB"},
		{-(2 << 60), "-2 EiB"},
	}
	for _, test := range tests {
		testBytes(t, test.value, test.expected)
	}
	// Verify extreme values of various types.
	testBytes[int16](t, math.MinInt16, "-32 KiB")
	testBytes[int16](t, math.MaxInt16, "32 KiB", skipParse)

	testBytes[uint16](t, math.MaxUint16, "64 KiB", skipParse)

	testBytes[int32](t, math.MinInt32, "-2 GiB")
	testBytes[int32](t, math.MaxInt32, "2 GiB", skipParse)

	testBytes[uint32](t, math.MaxUint32, "4 GiB", skipParse)

	testBytes[int64](t, math.MinInt64, "-8 EiB")
	testBytes[int64](t, math.MaxInt64, "8 EiB", skipParse)

	testBytes[uint64](t, math.MaxUint64, "16 EiB", skipParse)
}

const skipParse = true

func testBytes[T Integer](t *testing.T, value T, expected string, skipParse ...bool) {
	t.Helper()
	result := string(Bytes(value))
	if result != expected {
		t.Errorf("Bytes(%d) = %s; expected %s", value, result, expected)
	}
	// Compact string should be the same as the expected string but without spaces and "i".
	expectedCompact := strings.ReplaceAll(expected, " ", "")
	expectedCompact = strings.ReplaceAll(expectedCompact, "i", "")
	resultCompact := string(BytesCompact(value))
	if resultCompact != expectedCompact {
		t.Errorf("BytesCompact(%d) = %s; expected %s", value, resultCompact, expectedCompact)
	}

	if len(skipParse) > 0 {
		// Special case: this value cannot be parsed back.
		return
	}

	parsed, err := ParseBytes[T](result)
	if err != nil {
		t.Fatalf("could not parse Bytes[%s](%d)=%q: %v", reflect.TypeOf(value), value, result, err)
	}
	delta := math.Abs(float64(parsed) - float64(value))
	if delta > 0.05*math.Abs(float64(value)) {
		t.Fatalf("invalid parse result on Bytes[%s](%d)=%q: %d", reflect.TypeOf(value), value, result, parsed)
	}
	parsed2, err := ParseBytes[T](resultCompact)
	if err != nil {
		t.Fatalf("could not parse BytesCompact[%s](%d)=%q: %v", reflect.TypeOf(value), value, resultCompact, err)
	}
	if parsed2 != parsed {
		t.Fatalf("unexpected parse result on BytesCompact[%s](%d)=%q: %d", reflect.TypeOf(value), value, resultCompact, parsed)
	}
}

func TestBytesExact(t *testing.T) {
	tests := []struct {
		value    int64
		expected string
	}{
		{0, "0 B"},
		{1, "1 B"},
		{-1, "-1 B"},
		{2, "2 B"},
		{100, "100 B"},
		{-100, "-100 B"},
		{900, "900 B"},
		{1000, "1,000 B"},
		{1023, "1,023 B"},
		{1024, "1 KiB"},
		{1090, "1,090 B"},
		{2900, "2,900 B"},
		{10100, "10,100 B"},
		{10240, "10 KiB"},
		{1_000_000, "1,000,000 B"},
		{12 << 20, "12 MiB"},
		{12<<20 + 1, "12,582,913 B"},
		{12<<20 + 1024, "12,289 KiB"},
		{1200 << 20, "1,200 MiB"},
		{12 << 30, "12 GiB"},
		{1200 << 30, "1,200 GiB"},
		{12 << 40, "12 TiB"},
		{1200 << 40, "1,200 TiB"},
		{12 << 50, "12 PiB"},
		{12<<50 + 1, "13,510,798,882,111,489 B"},
		{1200 << 50, "1,200 PiB"},
		{2 << 60, "2 EiB"},
		{-(2 << 60), "-2 EiB"},
		{math.MaxInt64, "9,223,372,036,854,775,807 B"},
		{math.MinInt64, "-8 EiB"},
	}
	for _, test := range tests {
		testBytesExact(t, test.value, test.expected)
	}

	// Verify extreme values of various types.
	testBytesExact[int16](t, math.MinInt16, "-32 KiB")
	testBytesExact[int16](t, math.MaxInt16, "32,767 B")

	testBytesExact[uint16](t, math.MaxUint16, "65,535 B")

	testBytesExact[int32](t, math.MinInt32, "-2 GiB")
	testBytesExact[int32](t, math.MaxInt32, "2,147,483,647 B")

	testBytesExact[uint32](t, math.MaxUint32, "4,294,967,295 B")

	testBytesExact[int64](t, math.MinInt64, "-8 EiB")
	testBytesExact[int64](t, math.MaxInt64, "9,223,372,036,854,775,807 B")

	testBytesExact[uint64](t, math.MaxUint64, "18,446,744,073,709,551,615 B")

	for i := 0; i < 10000; i++ {
		n := rand.Uint64()
		if i%4 == 0 {
			// Test values around powers of two.
			n = uint64(1) << (n % 63)
			n = n + rand.Uint64N(10) - 5
		}
		testBytesExactRoundtrip[int8](t, n)
		testBytesExactRoundtrip[uint8](t, n)
		testBytesExactRoundtrip[int16](t, n)
		testBytesExactRoundtrip[uint16](t, n)
		testBytesExactRoundtrip[int32](t, n)
		testBytesExactRoundtrip[uint32](t, n)
		testBytesExactRoundtrip[int64](t, n)
		testBytesExactRoundtrip[uint64](t, n)
		testBytesExactRoundtrip[int](t, n)
		testBytesExactRoundtrip[uint](t, n)
	}
}

func testBytesExact[T Integer](t *testing.T, value T, expected string) {
	t.Helper()
	result := string(BytesExact(value))
	if result != expected {
		t.Errorf("BytesExact(%d) = %s; expected %s", value, result, expected)
	}
	parsed, err := ParseBytes[T](result)
	if err != nil {
		t.Fatalf("could not parse BytesExact(%d)=%q: %v", value, result, err)
	}
	if parsed != value {
		t.Fatalf("invalid parse result on BytesExact(%d)=%q: %d", value, result, parsed)
	}
}

func testBytesExactRoundtrip[T Integer](t *testing.T, valueUint uint64) {
	value := T(valueUint)
	t.Helper()
	result := string(BytesExact[T](value))
	parsed, err := ParseBytes[T](result)
	if err != nil {
		t.Fatalf("could not parse BytesExact(%d)=%q: %v", value, result, err)
	}
	if parsed != value {
		t.Fatalf("invalid parse result on BytesExact(%d)=%q: %d", value, result, parsed)
	}
}

func TestParseBytesErrors(t *testing.T) {
	expectParseErr[uint64](t, "18,446,744,073,709,551,616 B")
	expectParseErr[uint64](t, "18 EiB")
	expectParseErr[uint64](t, "123.45 EiB")
	expectParseErr[uint64](t, "100,000,000,000,000,000,000 B")
	expectParseErr[int64](t, "9,223,372,036,854,775,808")
	expectParseErr[int64](t, "-9,223,372,036,854,775,809")
	expectParseErr[int64](t, "9.1EB")

	expectParseErr[uint32](t, "4,294,967,296")
	expectParseErr[int32](t, "4,294,967,296")
	expectParseErr[int32](t, "2GB")
	expectParseErr[int32](t, "-2.1GB")
	expectParseErr[int32](t, "-2,147,483,649 B")

	expectParseErr[uint16](t, "65536")
	expectParseErr[uint16](t, "655361234")
	expectParseErr[int16](t, "32768")
	expectParseErr[int16](t, "32KiB")
	expectParseErr[int16](t, "-32.1KB")
	expectParseErr[int16](t, "-32769 B")

	expectParseErr[uint8](t, "256")
	expectParseErr[int8](t, "128")
	expectParseErr[int8](t, "-129")
}

func expectParseErr[T Integer](t *testing.T, s string) {
	if _, err := ParseBytes[T](s); err == nil {
		t.Helper()
		t.Errorf("expected error parsing %q", s)
	}
}
