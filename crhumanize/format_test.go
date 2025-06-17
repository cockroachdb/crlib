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
	"strings"
	"testing"

	"github.com/cockroachdb/crlib/crstrings"
	"github.com/cockroachdb/crlib/internal/datadriven"
)

func TestFormat(t *testing.T) {
	datadriven.Walk(t, "testdata/format", func(t *testing.T, path string) {
		datadriven.RunTest(t, path, func(t *testing.T, td *datadriven.TestData) string {
			switch td.Cmd {
			case "uint64":
				return runFormatTest[uint64](t, td)
			case "int64":
				return runFormatTest[int64](t, td)
			case "uint32":
				return runFormatTest[uint32](t, td)
			case "int32":
				return runFormatTest[int32](t, td)
			case "uint16":
				return runFormatTest[uint16](t, td)
			case "int16":
				return runFormatTest[int16](t, td)
			case "uint8":
				return runFormatTest[uint8](t, td)
			case "int8":
				return runFormatTest[int8](t, td)
			default:
				td.Fatalf(t, "unknown command: %s", td.Cmd)
				return ""
			}
		})
	})
}

func runFormatTest[T Integer](t *testing.T, td *datadriven.TestData) string {
	var unit string
	td.MaybeScanArgs(t, "unit", &unit)

	table := [][5]string{{"Value", "IEC", "IEC Exact", "SI", "SI Exact"}}

	for _, l := range crstrings.Lines(td.Input) {
		var val T
		if _, err := fmt.Sscanf(l, "%d", &val); err != nil {
			td.Fatalf(t, "error parsing %q: %v", l, err)
		}
		table = append(table, [5]string{
			fmt.Sprint(val),
			string(Format(val, IEC, unit)),
			string(Format(val, IEC, unit, Exact)),
			string(Format(val, SI, unit)),
			string(Format(val, SI, unit, Exact)),
		})
	}
	var colLens [5]int
	for i := range table {
		for j := range colLens {
			colLens[j] = max(colLens[j], len(table[i][j]))
		}
	}
	var buf strings.Builder
	for i := range table {
		fmt.Fprintf(&buf, "%*s  %*s  %*s  %*s  %*s\n",
			colLens[0], table[i][0],
			colLens[1], table[i][1],
			colLens[2], table[i][2],
			colLens[3], table[i][3],
			colLens[4], table[i][4],
		)
	}
	return buf.String()
}
