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

package cralloc

import (
	"math/rand/v2"
	"testing"

	"github.com/cockroachdb/crlib/testutils/require"
)

func TestScratchBuffer(t *testing.T) {
	s := (*ScratchBuffer)(nil).AllocUnsafe(100)
	require.Equal(t, len(s), 100)
	var sb ScratchBuffer
	s = sb.AllocUnsafe(100)
	require.Equal(t, len(s), 100)
	c := cap(s)
	s = sb.AllocUnsafe(50)
	require.Equal(t, len(s), 50)
	require.Equal(t, cap(s), c)
	s = sb.AllocUnsafe(101)
	require.Equal(t, len(s), 101)
	require.GT(t, cap(s), 101)

	t.Run("AllocZero", func(t *testing.T) {
		for range 100 {
			var sb ScratchBuffer
			maxN := 1 + rand.IntN(1000)
			for range 20 {
				n := rand.IntN(maxN)
				b := sb.AllocZeroUnsafe(n)
				for i := range b {
					require.Equal(t, b[i], 0)
				}
				// Trash the entire buffer.
				b = b[:cap(b)]
				for i := range b {
					b[i] = 0xcc
				}
			}
		}
	})

	t.Run("Append", func(t *testing.T) {
		var sb ScratchBuffer
		b := sb.AllocUnsafe(100)
		b = sb.Append(b, make([]byte, 1000)...)
		require.Equal(t, len(b), 1100)
		// Ensure the capacity has grown.
		require.GE(t, sb.Capacity(), 1100)

		// Append an unrelated slice.
		b = sb.Append(make([]byte, 1100), make([]byte, 10000)...)
		require.Equal(t, len(b), 11100)
		// Ensure the capacity did not grow.
		require.LT(t, sb.Capacity(), 10000)
	})
}
