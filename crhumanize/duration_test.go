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
	"strings"
	"testing"
	"time"

	"github.com/cockroachdb/crlib/crstrings"
	"github.com/cockroachdb/crlib/internal/datadriven"
)

func TestDuration(t *testing.T) {
	datadriven.RunTest(t, "testdata/duration", func(t *testing.T, td *datadriven.TestData) string {
		if td.Cmd != "duration" {
			td.Fatalf(t, "unknown command: %q", td.Cmd)
		}
		var buf strings.Builder
		for _, l := range crstrings.Lines(td.Input) {
			d, err := time.ParseDuration(l)
			if err != nil {
				td.Fatalf(t, "could not parse duration %q: %v", l, err)
			}
			fmt.Fprintf(&buf, "%s -> %s\n", d, Duration(d))
		}
		return buf.String()
	})
}

func TestDurationError(t *testing.T) {
	for _, v := range []time.Duration{time.Microsecond, time.Second, time.Minute, time.Hour, 100 * time.Hour, 10000 * time.Hour} {
		for i := 0; i < 1000; i++ {
			d := time.Duration(rand.Int64N(int64(v)))
			s := string(Duration(d))
			d1, err := time.ParseDuration(s)
			if err != nil {
				t.Fatalf("%s: could not parse duration %q: %v", d, s, err)
			}
			if relativeErr := math.Abs(float64(d1-d)) / float64(d); relativeErr > 0.05 {
				t.Fatalf("%s -> %s -> %s error is too large: %f\n", d, s, d1, relativeErr)
			}
		}
	}
}
