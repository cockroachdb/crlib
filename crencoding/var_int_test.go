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

package crencoding

import (
	"encoding/binary"
	"math"
	"math/rand/v2"
	"testing"

	"github.com/cockroachdb/crlib/testutils/require"
)

// TestUvarintLen tests UvarintLen32 and UvarintLen64.
func TestUvarintLen(t *testing.T) {
	check := func(n uint64) {
		t := require.WithMsgf(t, "n=%d", n)
		res64 := UvarintLen64(n)
		require.Equal(t, res64, len(binary.AppendUvarint(nil, n)))

		res32 := UvarintLen32(uint32(n))
		require.Equal(t, res32, len(binary.AppendUvarint(nil, uint64(uint32(n)))))
	}
	check(0)
	check(math.MaxUint64)
	for i := uint64(0); i < 64; i++ {
		check(1<<i - 1)
		check(1 << i)
		check(1<<i + 1)
	}
	for i := 0; i < 100000; i++ {
		check(rand.Uint64() >> rand.UintN(64))
	}
}
