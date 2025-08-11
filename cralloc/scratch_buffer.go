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

import "unsafe"

// ScratchBuffer is a helper for the common pattern of reusing a byte buffer to
// reduce slice allocations. To use, replace `make([]byte, n)` with
// `sb.Alloc(n)`.
type ScratchBuffer struct {
	p        unsafe.Pointer
	capacity int
}

// AllocUnsafe returns a byte slice of length n and arbitrary capacity which can
// be used until the next call to AllocUnsafe/AllocZero/Append.
//
// WARNING: the slice contains arbitrary data.
//
// This method is marked Unsafe because the allowed lifetime of the returned
// slice is limited.
//
// If the receiver is nil, always allocates a new slice.
func (sb *ScratchBuffer) AllocUnsafe(n int) []byte {
	if sb == nil {
		return make([]byte, n)
	}
	s := unsafe.Slice((*byte)(sb.p), sb.capacity)
	if sb.capacity >= n {
		return s[:n]
	}
	// Adapted from slices.Grow().
	s = append(s[:0], make([]byte, n)...)
	sb.p = unsafe.Pointer(&s[0])
	sb.capacity = cap(s)
	return s
}

// AllocZeroUnsafe returns a byte slice of length n and arbitrary capacity which
// can be used until the next call to AllocUnsafe/AllocZero/Append. The slice is
// zeroed out.
//
// WARNING: the slice contains arbitrary data between the length and the
// capacity.
//
// This method is marked Unsafe because the allowed lifetime of the returned
// slice is limited.
//
// If the receiver is nil, always allocates a new slice.
func (sb *ScratchBuffer) AllocZeroUnsafe(n int) []byte {
	if sb == nil {
		return make([]byte, n)
	}
	s := unsafe.Slice((*byte)(sb.p), sb.capacity)
	if sb.capacity >= n {
		s = s[:n]
		clear(s)
		return s
	}
	// Adapted from slices.Grow(). We do not want to simply use make([]byte, n)
	// because we want the scratch buffer to grow according to the append()
	// heuristics. Otherwise, an allocation pattern of slowly increasing sizes
	// would cause an allocation each time.
	s = append(s[:0], make([]byte, n)...)
	sb.p = unsafe.Pointer(&s[0])
	sb.capacity = cap(s)
	return s
}

// Append is like the built-in append(), but it also updates the scratch buffer
// so that any newly allocated buffer can be reused.
//
// Append can be used with buffers not allocated through the scratch buffer (in
// which case the scratch buffer is not updated).
func (sb *ScratchBuffer) Append(buf []byte, values ...byte) []byte {
	res := append(buf, values...)
	if sb != nil && unsafe.SliceData(buf) == (*byte)(sb.p) && unsafe.SliceData(res) != (*byte)(sb.p) {
		sb.p = unsafe.Pointer(unsafe.SliceData(res))
		sb.capacity = cap(res)
	}
	return res
}

// Capacity returns the current capacity.
func (sb *ScratchBuffer) Capacity() int {
	if sb == nil {
		return 0
	}
	return sb.capacity
}

// Reset clears the buffer. This can be useful if we want to avoid retaining a
// very large buffer.
func (sb *ScratchBuffer) Reset() {
	if sb != nil {
		*sb = ScratchBuffer{}
	}
}
