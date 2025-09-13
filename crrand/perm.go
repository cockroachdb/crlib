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

// Package crrand implements functionality related to pseudorandom number
// generation.
package crrand

import (
	"math/bits"
	"math/rand/v2"
)

// MakePerm64 constructs a new Perm64 from a 64-bit seed, providing a
// deterministic, pseudorandom, bijective mapping of 64-bit values X to 64-bit
// values Y.
func MakePerm64(seed uint64) Perm64 {
	prng := rand.New(rand.NewPCG(seed, seed))
	return Perm64{
		seed: [4]uint32{
			prng.Uint32(),
			prng.Uint32(),
			prng.Uint32(),
			prng.Uint32(),
		},
	}
}

// A Perm64 provides a deterministic, pseudorandom permutation of 64-bit values.
type Perm64 struct {
	seed [4]uint32
}

// At returns the nth value in the permutation of the 64-bit values. The return
// value may be passed to Index to recover n. The permutation is pseudorandom.
func (p Perm64) At(n uint64) uint64 {
	// Use a simple Feistel network with 4 rounds to shuffle data.
	L := uint32(n >> 32)
	R := uint32(n)
	for r := range p.seed {
		t := arx(R^p.seed[r], p.seed[(r+1)&3])
		L, R = R, L^t
	}
	return (uint64(L) << 32) | uint64(R)
}

// IndexOf inverts the permutation, returning the index of the provided value in
// the permutation. If y was produced by At(x), then IndexOf(y) returns x.
func (p Perm64) IndexOf(y uint64) uint64 {
	L := uint32(y >> 32)
	R := uint32(y)
	for r := 3; r >= 0; r-- {
		// reverse of: L, R = R, L ^ arx(R^k[r], k[(r+1)&3])
		prevR := L
		prevL := R ^ arx(prevR^p.seed[r], p.seed[(r+1)&3])
		L, R = prevL, prevR
	}
	return (uint64(L) << 32) | uint64(R)
}

// ARX-only round function.
func arx(x, k uint32) uint32 {
	x ^= k
	x += bits.RotateLeft32(x, 5)
	x ^= bits.RotateLeft32(x, 7)
	x += bits.RotateLeft32(x, 16)
	return x
}
