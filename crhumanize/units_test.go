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
	"testing"
)

func TestRound(t *testing.T) {
	testCases := []struct {
		units    Units
		value    uint64
		expected string
	}{
		{SI, 0, "0"},
		{SI, 1, "1"},
		{SI, 10, "10"},
		{SI, 100, "100"},
		{SI, 900, "900"},
		{SI, 123_400_000, "123.4M"},
		{SI, 678_900_000, "678.9M"},
		{SI, math.MaxUint64, "18.45E"},

		{IEC, 0, "0"},
		{IEC, 1, "1"},
		{IEC, 10, "10"},
		{IEC, 100, "100"},
		{IEC, 900, "900"},
		{IEC, 1000, "0.98Ki"},
		{IEC, 123_400_000, "117.68Mi"},
		{IEC, 678_900_000, "647.45Mi"},
		{IEC, math.MaxUint64, "16Ei"},
	}
	for i, tc := range testCases {
		scaled, suffix := tc.units.Round(tc.value)
		res := string(Float(scaled, 2)) + suffix
		if res != tc.expected {
			t.Errorf("%d: expected %s, got %s for value %d", i, tc.expected, res, tc.value)
		}
	}
}

func TestParseUnit(t *testing.T) {
	testCases := []struct {
		unit     string
		expected uint64
	}{
		{"", 1},

		{"K", 1000},
		{"m", 1000 * 1000},
		{"G", 1000 * 1000 * 1000},
		{"T", 1000 * 1000 * 1000 * 1000},
		{"p", 1000 * 1000 * 1000 * 1000 * 1000},
		{"E", 1000 * 1000 * 1000 * 1000 * 1000 * 1000},

		{"Ki", 1024},
		{"MI", 1024 * 1024},
		{"gI", 1024 * 1024 * 1024},
		{"Ti", 1024 * 1024 * 1024 * 1024},
		{"pI", 1024 * 1024 * 1024 * 1024 * 1024},
		{"eI", 1024 * 1024 * 1024 * 1024 * 1024 * 1024},
	}
	for i, tc := range testCases {
		scale, err := parseUnit(tc.unit)
		if err != nil {
			t.Errorf("%d: unexpected error for unit %s: %v", i, tc.unit, err)
			continue
		}
		if scale != tc.expected {
			t.Errorf("%d: expected %d for unit %s, got %d", i, tc.expected, tc.unit, scale)
		}
	}
}
