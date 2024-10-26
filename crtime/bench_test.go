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
	"fmt"
	"io"
	"testing"
	"time"
)

// Sample results on an Apple M1 (darwin/arm64):
//
//	Mono/NowMono-10   41.8ns ± 1%
//	Mono/time.Now-10  66.8ns ± 3%
//
// On Intel(R) Xeon(R) CPU @ 2.80GHz (linux/amd64):
//
//	Mono/NowMono-24   49.0ns ± 0%
//	Mono/time.Now-24  65.1ns ± 0%
func BenchmarkMono(b *testing.B) {
	b.Run("NowMono", func(b *testing.B) {
		var d time.Duration
		for i := 0; i < b.N; i++ {
			m := NowMono()
			d += m.Elapsed()
		}
		fmt.Fprint(io.Discard, d)
	})
	b.Run("time.Now", func(b *testing.B) {
		var d time.Duration
		for i := 0; i < b.N; i++ {
			t := time.Now()
			d += time.Since(t)
		}
		fmt.Fprint(io.Discard, d)
	})
}
