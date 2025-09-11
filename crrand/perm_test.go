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

package crrand

import (
	"math"
	"math/rand/v2"
	"testing"
	"time"
)

var interestingUint64s = []uint64{
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 32, 33, 63, 64, 65, 129, 3050, 29356,
	297532935, 2973539791203, 0x9E3779B97F4A7C15, math.MaxUint64 - 1,
	math.MaxUint64,
}

func TestPerm64(t *testing.T) {
	for _, seed := range interestingUint64s {
		p := MakePerm64(seed)
		for _, x := range interestingUint64s {
			y := p.At(x)
			x2 := p.IndexOf(y)
			if x != x2 {
				t.Errorf(".At(%d) = %d, .IndexOf(%d) = %d, want %d", x, y, y, x2, x)
			}
		}
	}
}

func TestPerm64Random(t *testing.T) {
	seed := uint64(time.Now().UnixNano())
	defer func() {
		if t.Failed() {
			t.Logf("seed: %d", seed)
		}
	}()
	rng := rand.New(rand.NewPCG(seed, seed))
	p := MakePerm64(rng.Uint64())
	for i := 0; i < 1000; i++ {
		x := rng.Uint64()
		y := p.At(x)
		x2 := p.IndexOf(y)
		if x != x2 {
			t.Errorf("p.At(%d) = %d, p.IndexOf(%d) = %d, want %d", x, y, y, x2, x)
		}
	}
}
