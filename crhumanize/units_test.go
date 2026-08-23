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

// TestUnitTables verifies that the SI and IEC tables are parallel: the prefixes
// correspond to each other one-to-one (differing only by the i indicator) and
// the corresponding scales are within 20% of each other. parseUnit relies on
// this correspondence to map an SI prefix to the equivalent IEC scale.
func TestUnitTables(t *testing.T) {
	if len(SI.prefixes) != len(IEC.prefixes) {
		t.Fatalf("SI and IEC have a different number of prefixes: %d vs %d", len(SI.prefixes), len(IEC.prefixes))
	}
	for i := range SI.prefixes {
		if IEC.prefixes[i] != SI.prefixes[i]+"i" {
			t.Errorf("%d: IEC prefix %q does not correspond to SI prefix %q", i, IEC.prefixes[i], SI.prefixes[i])
		}
		if ratio := float64(IEC.scales[i]) / float64(SI.scales[i]); math.Abs(ratio-1) > 0.2 {
			t.Errorf("%d: IEC scale %d and SI scale %d differ by more than 20%%", i, IEC.scales[i], SI.scales[i])
		}
	}
}

func TestParseUnit(t *testing.T) {
	testCases := []struct {
		unit        string
		expected    uint64
		expectedIEC uint64
	}{
		{"", 1, 1},

		{"K", 1000, 1024},
		{"m", 1000 * 1000, 1024 * 1024},
		{"G", 1000 * 1000 * 1000, 1024 * 1024 * 1024},
		{"T", 1000 * 1000 * 1000 * 1000, 1024 * 1024 * 1024 * 1024},
		{"p", 1000 * 1000 * 1000 * 1000 * 1000, 1024 * 1024 * 1024 * 1024 * 1024},
		{"E", 1000 * 1000 * 1000 * 1000 * 1000 * 1000, 1024 * 1024 * 1024 * 1024 * 1024 * 1024},

		{"Ki", 1024, 1024},
		{"MI", 1024 * 1024, 1024 * 1024},
		{"gI", 1024 * 1024 * 1024, 1024 * 1024 * 1024},
		{"Ti", 1024 * 1024 * 1024 * 1024, 1024 * 1024 * 1024 * 1024},
		{"pI", 1024 * 1024 * 1024 * 1024 * 1024, 1024 * 1024 * 1024 * 1024 * 1024},
		{"eI", 1024 * 1024 * 1024 * 1024 * 1024 * 1024, 1024 * 1024 * 1024 * 1024 * 1024 * 1024},
	}
	for i, tc := range testCases {
		for _, forceIEC := range []bool{false, true} {
			expected := tc.expected
			if forceIEC {
				expected = tc.expectedIEC
			}
			scale, err := parseUnit(tc.unit, forceIEC)
			if err != nil {
				t.Errorf("%d: unexpected error for unit %s (forceIEC=%t): %v", i, tc.unit, forceIEC, err)
				continue
			}
			if scale != expected {
				t.Errorf("%d: expected %d for unit %s (forceIEC=%t), got %d", i, expected, tc.unit, forceIEC, scale)
			}
		}
	}

	// Unknown prefixes must be rejected, with or without forceIEC.
	invalid := []string{"X", "Xi", "i", "I", "KK", "KiB", "1", "K i", " K"}
	for _, unit := range invalid {
		for _, forceIEC := range []bool{false, true} {
			if scale, err := parseUnit(unit, forceIEC); err == nil {
				t.Errorf("expected error for unit %q (forceIEC=%t), got %d", unit, forceIEC, scale)
			}
		}
	}
}
