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
	"math/rand/v2"
	"reflect"
	"strings"
	"testing"

	"github.com/cockroachdb/crlib/crstrings"
	"github.com/cockroachdb/crlib/internal/datadriven"
)

func TestParse(t *testing.T) {
	datadriven.Walk(t, "testdata/parse", func(t *testing.T, path string) {
		datadriven.RunTest(t, path, func(t *testing.T, td *datadriven.TestData) string {
			switch td.Cmd {
			case "uint64":
				return runParseTest[uint64](t, td)
			case "int64":
				return runParseTest[int64](t, td)
			case "uint32":
				return runParseTest[uint32](t, td)
			case "int32":
				return runParseTest[int32](t, td)
			case "uint16":
				return runParseTest[uint16](t, td)
			case "int16":
				return runParseTest[int16](t, td)
			case "uint8":
				return runParseTest[uint8](t, td)
			case "int8":
				return runParseTest[int8](t, td)
			default:
				td.Fatalf(t, "unknown command: %s", td.Cmd)
				return ""
			}
		})
	})
}

func runParseTest[T Integer](t *testing.T, td *datadriven.TestData) string {
	var unit string
	td.MaybeScanArgs(t, "unit", &unit)

	var buf strings.Builder
	for _, l := range crstrings.Lines(td.Input) {
		result, err := Parse[T](l, unit)
		if err != nil {
			fmt.Fprintf(&buf, "%s: %s\n", l, err)
		} else {
			fmt.Fprintf(&buf, "%s: %d\n", l, result)
		}
	}
	return buf.String()
}

func TestRoundtrip(t *testing.T) {
	for i := 0; i < 10000; i++ {
		n := rand.Uint64()
		if i%4 == 0 {
			// Test values around powers of two.
			n = uint64(1) << (n % 65)
			n = n + rand.Uint64N(10) - 5
		}
		testRoundtrip[int8](t, int8(n))
		testRoundtrip[int8](t, -int8(n))
		testRoundtrip[uint8](t, uint8(n))

		testRoundtrip[int16](t, int16(n))
		testRoundtrip[int16](t, -int16(n))
		testRoundtrip[uint16](t, uint16(n))

		testRoundtrip[int32](t, int32(n))
		testRoundtrip[int32](t, -int32(n))
		testRoundtrip[uint32](t, uint32(n))

		testRoundtrip[int64](t, int64(n))
		testRoundtrip[int64](t, -int64(n))
		testRoundtrip[uint64](t, n)
	}
}

func testRoundtrip[T Integer](t *testing.T, value T) {
	t.Helper()
	var flags []FmtFlag
	if rand.IntN(2) == 0 {
		flags = append(flags, Compact)
	}

	formatted := Bytes(value, flags...)
	parsed, err := ParseBytes[T](string(formatted))
	if err != nil {
		t.Fatalf("could not parse Bytes[%s](%d)=%q: %v", reflect.TypeOf(value), value, formatted, err)
	}
	if delta := math.Abs(float64(parsed) - float64(value)); delta > 0.05*math.Abs(float64(value)) {
		t.Fatalf("invalid parse result on Bytes[%s](%d)=%q: %d (delta %.1f%%)", reflect.TypeOf(value), value, formatted, parsed, delta*100/float64(value))
	}

	flags = append(flags, Exact)
	formatted = Bytes(value, flags...)
	parsed, err = ParseBytes[T](string(formatted))
	if err != nil {
		t.Fatalf("could not parse Bytes[%s](%d)=%q: %v", reflect.TypeOf(value), value, formatted, err)
	}
	if parsed != value {
		t.Fatalf("invalid parse result on Bytes[%s](%d, Exact)=%q: %d", reflect.TypeOf(value), value, formatted, parsed)
	}
}
