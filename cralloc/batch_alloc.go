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

package cralloc

import "sync"

// BatchAllocator is used to allocate small objects in batches, reducing the
// number of individual allocations.
//
// The tradeoff is that the lifetime of the objects in a batch are tied
// together, which can potentially result in higher memory usage. In addition,
// there can be O(GOMAXPROCS) extra instantiated batches at any one time.
// BatchAllocator should be used when T is small and it does not contain
// references to large objects.
//
// Sample usage:
//
//	var someTypeBatchAlloc = MakeBatchAllocator[SomeType]()  // global
//		...
//		x := someTypeBatchAlloc.Alloc()
type BatchAllocator[T any] struct {
	// We use a sync.Pool as an approximation to maintaining one batch per CPU.
	// This is more efficient than using a mutex and provides good memory
	// locality.
	pool sync.Pool
}

// MakeBatchAllocator initializes a BatchAllocator.
func MakeBatchAllocator[T any]() BatchAllocator[T] {
	return BatchAllocator[T]{
		pool: sync.Pool{
			New: func() any {
				return &batch[T]{}
			},
		},
	}
}

const batchSize = 8

// Init must be called before the batch allocator can be used.
func (ba *BatchAllocator[T]) Init() {
	ba.pool.New = func() any {
		return &batch[T]{}
	}
}

// Alloc returns a new zeroed out instance of T.
func (ba *BatchAllocator[T]) Alloc() *T {
	b := ba.pool.Get().(*batch[T])
	// If Init() was not called, the first Alloc() will panic here.
	t := &b.buf[b.used]
	b.used++
	if b.used < batchSize {
		// Batch has more objects available, put it back into the pool.
		ba.pool.Put(b)
	}
	return t
}

type batch[T any] struct {
	// elements buf[:used] have been returned via Alloc. The rest are unused and
	// zero.
	buf  [batchSize]T
	used int8
}
