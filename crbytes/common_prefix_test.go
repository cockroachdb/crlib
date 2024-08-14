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

package crbytes

import (
	"bytes"
	"math/rand"
	"testing"
)

func TestCommonPrefixAllLengths(t *testing.T) {
	// Construct cases with each length up to a certain size.
	for l := 0; l <= 256; l++ {
		for k := 0; k <= l; k++ {
			a := bytes.Repeat([]byte("x"), l)
			b := bytes.Repeat([]byte("x"), l)
			if k < l {
				b[k] = '0'
			}
			if res := CommonPrefix(a, b); res != k {
				t.Errorf("length=%d expected=%d result=%d\n", l, k, res)
			}
			// Always test the generic implementation too.
			if res := commonPrefixGeneric(a, b); res != k {
				t.Errorf("length=%d expected=%d result=%d\n", l, k, res)
			}
		}
	}
}

func TestCommonPrefixRand(t *testing.T) {
	for _, tc := range []struct {
		maxLen   int
		alphabet int
	}{
		{maxLen: 4, alphabet: 2},
		{maxLen: 100, alphabet: 2},
		{maxLen: 200, alphabet: 2},
		{maxLen: 10, alphabet: 4},
		{maxLen: 500, alphabet: 4},
		{maxLen: 10, alphabet: 26},
		{maxLen: 500, alphabet: 26},
	} {
		for n := 0; n < 1000; n++ {
			a := genBytes(rand.Intn(tc.maxLen+1), tc.alphabet)
			b := genBytes(rand.Intn(tc.maxLen+1), tc.alphabet)
			expected := commonPrefixNaive(a, b)
			if res := CommonPrefix(a, b); res != expected {
				t.Errorf("%q %q expected=%d result=%d\n", a, b, expected, res)
			}
			// Always test the generic implementation too.
			if res := commonPrefixGeneric(a, b); res != expected {
				t.Errorf("%q %q expected=%d result=%d\n", a, b, expected, res)
			}
		}
	}
}

func commonPrefixNaive(a, b []byte) int {
	n := min(len(a), len(b))
	i := 0
	for i < n && a[i] == b[i] {
		i++
	}
	return i
}

func genBytes(length int, alphabet int) []byte {
	a := make([]byte, length)
	for i := range a {
		a[i] = 'a' + byte(rand.Intn(alphabet))
	}
	return a
}
