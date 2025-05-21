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

import "testing"

func TestCount(t *testing.T) {
	tests := []struct {
		value    int64
		expected string
	}{
		{0, "0"},
		{1, "1"},
		{-1, "-1"},
		{2, "2"},
		{100, "100"},
		{-100, "-100"},
		{900, "900"},
		{1000, "1K"},
		{1040, "1K"},
		{1050, "1.1K"},
		{1900, "1.9K"},
		{1951, "2K"},
		{9900, "9.9K"},
		{9951, "10K"},
		{10200, "10K"},
		{10600, "11K"},
		{12_000_000, "12M"},
		{1_200_000_000, "1.2G"},
		{12_000_000_000, "12G"},
		{1_200_000_000_000, "1.2T"},
		{12_000_000_000_000, "12T"},
		{1_200_000_000_000_000, "1.2P"},
		{12_000_000_000_000_000, "12P"},
		{1_200_000_000_000_000_000, "1.2E"},
		{2_000_000_000_000_000_000, "2E"},
		{-2_000_000_000_000_000_000, "-2E"},
	}
	for _, test := range tests {
		result := string(Count(test.value))
		if result != test.expected {
			t.Errorf("Count(%d) = %s; expected %s", test.value, result, test.expected)
		}
	}
}
