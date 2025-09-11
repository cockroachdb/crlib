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

import "math/bits"

// MakePerm64 constructs a new Mixer from a 64-bit seed, providing a
// deterministic, pseduorandom, bijective mapping of 64-bit values X to 64-bit
// values Y.
func MakePerm64(seed uint64) Perm64 {
	// derive 4 x 32-bit round keys from the 64-bit seed using only ARX ops.
	const c0 = 0x9E3779B97F4A7C15 // golden ratio (used here as XOR salt)
	const c1 = 0xC2B2AE3D27D4EB4F // a constant

	var m Perm64
	s0 := seed
	s1 := bits.RotateLeft64(seed^c0, 13)
	s2 := bits.RotateLeft64(seed^c1, 37)
	s3 := bits.RotateLeft64(seed^(c0^c1), 53)

	m.seed[0] = uint32(s0)
	m.seed[1] = uint32(s1 >> 32)
	m.seed[2] = uint32(s2)
	m.seed[3] = uint32(s3 >> 32)
	return m
}

// A Perm64 provides a deterministic, pseduorandom permutation of 64-bit values.
type Perm64 struct {
	seed [4]uint32
}

// Nth returns the nth value in the permutation of the 64-bit values. The return
// value may be passed to Index to recover n. The permutation is pseduorandom.
func (p Perm64) Nth(n uint64) uint64 {
	// Use a simple Feistel network with 4 rounds to shuffle data.

	L := uint32(n >> 32)
	R := uint32(n)
	for r := range p.seed {
		t := f(R^p.seed[r], p.seed[(r+1)&3])
		L, R = R, L^t
	}
	return (uint64(L) << 32) | uint64(R)
}

// Index inverts the permutation, returning the index of the provided value in
// the permutation. If y was produced by Nth(x), then Index(y) returns x.
func (p Perm64) Index(y uint64) uint64 {
	L := uint32(y >> 32)
	R := uint32(y)
	for r := 3; r >= 0; r-- {
		// reverse of: L, R = R, L ^ F(R^k[r], k[(r+1)&3])
		prevR := L
		prevL := R ^ f(prevR^p.seed[r], p.seed[(r+1)&3])
		L, R = prevL, prevR
	}
	return (uint64(L) << 32) | uint64(R)
}

// ARX-only round function.
func f(x, k uint32) uint32 {
	x ^= k
	x += bits.RotateLeft32(x, 5)
	x ^= bits.RotateLeft32(x, 7)
	x += bits.RotateLeft32(x, 16)
	return x
}
