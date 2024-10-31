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
	"math/bits"
	"unsafe"
)

// UvarintLen32 returns the number of bytes necessary for the Go
// encoding/binary.Uvarint encoding.
//
// It is always equivalent to len(binary.AppendUvarint(nil, uint32(x))) but faster.
func UvarintLen32(x uint32) int {
	// We |1 to avoid the special case of x=0.
	b := uint32(bits.Len32(x|1)) + 6
	// The result is b / 7. Instead of dividing by 7, we multiply by 37 which is
	// approximately 2^8/7 and then divide by 2^8. This approximation is exact for
	// small values in the range we care about.
	return int((b * 37) >> 8)
}

// UvarintLen64 returns the number of bytes necessary for the Go
// encoding/binary.Uvarint encoding.
//
// It is always equivalent to len(binary.AppendUvarint(nil, x)) but faster.
func UvarintLen64(x uint64) int {
	// We |1 to avoid the special case of x=0.
	b := uint32(bits.Len64(x|1)) + 6
	// The result is b / 7. Instead of dividing by 7, we multiply by 37 which is
	// approximately 2^8/7 and then divide by 2^8. This approximation is exact for
	// small values in the range we care about.
	return int((b * 37) >> 8)
}

// UvarintLenSum32 returns the total number of bytes to encode all the given
// integers using the Go encoding/binary.Uvarint encoding.
func UvarintLenSum32(values []uint32) int {
	p := unsafe.Pointer(unsafe.SliceData(values))
	// Note: it's illegal for an unsafe.Pointer to point outside our memory, so we
	// point at the last byte.
	end := unsafe.Pointer(uintptr(p) + uintptr(len(values))*4 - 1)
	var res int
	if len(values) >= 4 {
		end4 := unsafe.Pointer(uintptr(end) - 3*4)
		for uintptr(p) < uintptr(end4) {
			v := (*[4]uint32)(p)
			res += UvarintLen32(v[0])
			res += UvarintLen32(v[1])
			res += UvarintLen32(v[2])
			res += UvarintLen32(v[3])
			p = unsafe.Pointer(uintptr(p) + 4*4)
		}
	}
	for uintptr(p) < uintptr(end) {
		res += UvarintLen32(*(*uint32)(p))
		p = unsafe.Pointer(uintptr(p) + 4)
	}
	return res
}

// UvarintLenSum64 returns the total number of bytes to encode all the given
// integers using the Go encoding/binary.Uvarint encoding.
func UvarintLenSum64(values []uint64) int {
	p := unsafe.Pointer(unsafe.SliceData(values))
	// Note: it's illegal for an unsafe.Pointer to point outside our memory, so we
	// point at the last byte.
	end := unsafe.Pointer(uintptr(p) + uintptr(len(values))*8 - 1)
	var res int
	if len(values) >= 4 {
		end4 := unsafe.Pointer(uintptr(end) - 3*8)
		for uintptr(p) < uintptr(end4) {
			v := (*[4]uint64)(p)
			res += UvarintLen64(v[0])
			res += UvarintLen64(v[1])
			res += UvarintLen64(v[2])
			res += UvarintLen64(v[3])
			p = unsafe.Pointer(uintptr(p) + 4*8)
		}
	}
	for uintptr(p) < uintptr(end) {
		res += UvarintLen64(*(*uint64)(p))
		p = unsafe.Pointer(uintptr(p) + 8)
	}
	return res
}
