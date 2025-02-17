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

package crtime

import (
	"testing"
	"time"

	"github.com/cockroachdb/crlib/testutils/require"
)

func TestMono(t *testing.T) {
	a := NowMono()
	time.Sleep(10 * time.Millisecond)
	b := NowMono()
	require.GE(t, b.Sub(a), 9*time.Millisecond)
	c := MonoFromTime(time.Now())
	d := NowMono()
	require.LE(t, b, c)
	require.LE(t, c, d)

	t.Run("ToUTC", func(t *testing.T) {
		const d = 50 * time.Millisecond
		const tolerance = time.Millisecond

		start := NowMono()
		expected := time.Now().UnixNano()
		time.Sleep(d)
		actual := start.ToUTC().UnixNano()
		if actual < expected-tolerance.Nanoseconds() || actual > expected+tolerance.Nanoseconds() {
			t.Fatalf("actual - expected = %s", time.Duration(actual-expected))
		}
	})
}
