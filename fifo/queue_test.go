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

package fifo

import (
	"math/rand"
	"testing"

	"github.com/cockroachdb/crlib/testutils/require"
)

var pool = MakeQueueBackingPool[int]()

func TestQueue(t *testing.T) {
	q := MakeQueue[int](&pool)
	require.Equal(t, q.PeekFront(), nil)
	require.Equal(t, q.PeekBack(), nil)
	require.Equal(t, q.Len(), 0)
	q.PushBack(1)
	require.Equal(t, *q.PeekBack(), 1)
	q.PushBack(2)
	require.Equal(t, *q.PeekBack(), 2)
	q.PushBack(3)
	require.Equal(t, *q.PeekBack(), 3)
	require.Equal(t, q.Len(), 3)
	require.Equal(t, collect(&q), []int{1, 2, 3})
	*q.PeekBack() = 4
	require.Equal(t, collect(&q), []int{1, 2, 4})
	require.Equal(t, *q.PeekFront(), 1)
	q.PopFront()
	require.Equal(t, *q.PeekFront(), 2)
	q.PopFront()
	require.Equal(t, *q.PeekFront(), 4)
	require.Equal(t, *q.PeekBack(), 4)
	q.PopFront()
	require.Equal(t, q.PeekFront(), nil)
	require.Equal(t, q.PeekBack(), nil)
	require.Equal(t, len(collect(&q)), 0)

	for i := 1; i <= 1000; i++ {
		q.PushBack(i)
		require.Equal(t, q.Len(), i)
	}
	for i := 1; i <= 1000; i++ {
		require.Equal(t, *q.PeekFront(), i)
		q.PopFront()
		require.Equal(t, q.Len(), 1000-i)
	}
}

func TestQueueRand(t *testing.T) {
	q := MakeQueue[int](&pool)
	l, r := 0, 0
	for iteration := 0; iteration < 100; iteration++ {
		for n := rand.Intn(100); n > 0; n-- {
			r++
			q.PushBack(r)
			require.Equal(t, q.Len(), r-l)
			require.Equal(t, *q.PeekBack(), r)
		}
		for n := rand.Intn(q.Len() + 1); n > 0; n-- {
			l++
			require.Equal(t, *q.PeekFront(), l)
			q.PopFront()
			require.Equal(t, q.Len(), r-l)
		}
		if q.Len() == 0 {
			require.Equal(t, q.PeekBack(), nil)
		} else {
			require.Equal(t, *q.PeekBack(), r)
		}
		elems := collect(&q)
		require.Equal(t, len(elems), r-l)
		for i, e := range elems {
			require.Equal(t, e, l+1+i)
		}
		// Stopping the iteration early yields a prefix.
		var prefix []int
		for e := range q.All() {
			if len(prefix) == 3 {
				break
			}
			prefix = append(prefix, *e)
		}
		require.Equal(t, prefix, elems[:min(3, len(elems))])
	}
}

func collect(q *Queue[int]) []int {
	var res []int
	for e := range q.All() {
		res = append(res, *e)
	}
	return res
}
