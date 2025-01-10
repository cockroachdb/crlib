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
	"context"
	"errors"
	"math/rand"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/cockroachdb/crlib/testutils/require"
)

func TestSemaphoreAPI(t *testing.T) {
	s := NewSemaphore(10)
	require.Equal(t, s.TryAcquire(5), true)
	require.Equal(t, s.TryAcquire(10), false)
	require.Equal(t, "capacity: 10, outstanding: 5, num-had-to-wait: 0", s.Stats().String())

	ch := make(chan struct{}, 10)
	go func() {
		if err := s.Acquire(context.Background(), 8); err != nil {
			t.Error(err)
		}
		ch <- struct{}{}
		if err := s.Acquire(context.Background(), 1); err != nil {
			t.Error(err)
		}
		ch <- struct{}{}
		if err := s.Acquire(context.Background(), 5); err != nil {
			t.Error(err)
		}
		ch <- struct{}{}
	}()
	require.NoRecv(t, ch)
	s.Release(5)
	require.Recv(t, ch)
	require.Recv(t, ch)
	require.NoRecv(t, ch)
	s.Release(1)
	require.NoRecv(t, ch)
	s.Release(8)
	require.Recv(t, ch)

	require.True(t, strings.Contains(s.Stats().String(), "capacity: 10, outstanding: 5"))
	// Test UpdateCapacity.
	go func() {
		if err := s.Acquire(context.Background(), 8); err != nil {
			t.Error(err)
		}
		ch <- struct{}{}
		if err := s.Acquire(context.Background(), 1); err != nil {
			t.Error(err)
		}
		ch <- struct{}{}
	}()
	require.NoRecv(t, ch)
	s.UpdateCapacity(15)
	require.Recv(t, ch)
	require.Recv(t, ch)
	s.UpdateCapacity(2)
	go func() {
		// Request more than the capacity.
		if err := s.Acquire(context.Background(), 5); err != nil {
			t.Error(err)
		}
		ch <- struct{}{}
	}()
	require.NoRecv(t, ch)
	s.Release(5)
	require.NoRecv(t, ch)
	s.Release(8)
	require.NoRecv(t, ch)
	s.Release(1)
	// Last request should now be allowed, despite being larger than the capacity.
	require.Recv(t, ch)
}

// TestSemaphoreBasic is a test with multiple goroutines acquiring a unit and
// releasing it right after.
func TestSemaphoreBasic(t *testing.T) {
	capacities := []int64{1, 5, 10, 50, 100}
	goroutineCounts := []int{1, 10, 100}

	for _, capacity := range capacities {
		for _, numGoroutines := range goroutineCounts {
			s := NewSemaphore(capacity)
			ctx := context.Background()
			resCh := make(chan error, numGoroutines)

			for i := 0; i < numGoroutines; i++ {
				go func() {
					err := s.Acquire(ctx, 1)
					if err != nil {
						resCh <- err
						return
					}
					s.Release(1)
					resCh <- nil
				}()
			}

			for i := 0; i < numGoroutines; i++ {
				if err := require.Recv(t, resCh); err != nil {
					t.Fatal(err)
				}
			}

			if stats := s.Stats(); stats.Outstanding != 0 {
				t.Fatalf("expected nothing outstanding; got %s", stats)
			}
		}
	}
}

// TestSemaphoreContextCancellation tests the behavior that for an ongoing
// blocked acquisition, if the context passed in gets canceled the acquisition
// gets canceled too with an error indicating so.
func TestSemaphoreContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	s := NewSemaphore(1)
	if err := s.Acquire(ctx, 1); err != nil {
		t.Fatal(err)
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- s.Acquire(ctx, 1)
	}()

	cancel()

	err := require.Recv(t, errCh)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context cancellation error, got %v", err)
	}

	stats := s.Stats()
	require.Equal(t, stats.Capacity, 1)
	require.Equal(t, stats.Outstanding, 1)
}

// TestSemaphoreCanceledAcquisitions tests the behavior where we enqueue
// multiple acquisitions with canceled contexts and expect any subsequent
// acquisition with a valid context to proceed without error.
func TestSemaphoreCanceledAcquisitions(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	s := NewSemaphore(1)
	if err := s.Acquire(ctx, 1); err != nil {
		t.Fatal(err)
	}

	cancel()
	const numGoroutines = 5

	errCh := make(chan error)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			errCh <- s.Acquire(ctx, 1)
		}()
	}

	for i := 0; i < numGoroutines; i++ {
		if err := require.Recv(t, errCh); !errors.Is(err, context.Canceled) {
			t.Fatalf("expected context cancellation error, got %v", err)
		}
	}
	s.Release(1)

	go func() {
		errCh <- s.Acquire(context.Background(), 1)
	}()

	if err := require.Recv(t, errCh); err != nil {
		t.Fatal(err)
	}
}

// TestSemaphoreNumHadToWait checks Stats().NumHadToWait.
func TestSemaphoreNumHadToWait(t *testing.T) {
	s := NewSemaphore(1)
	ctx := context.Background()
	doneCh := make(chan struct{}, 10)
	doAcquire := func(ctx context.Context) {
		err := s.Acquire(ctx, 1)
		if ctx.Err() == nil {
			if err != nil {
				t.Error(err)
			}
			doneCh <- struct{}{}
		}
	}

	assertNumWaitersSoon := func(exp int64) {
		for i := 0; ; i++ {
			got := s.Stats().NumHadToWait
			if got == exp {
				return
			}
			if i >= 20 {
				t.Fatalf("expected num-had-to-wait to be %d, got %d", got, exp)
			}
			time.Sleep(50 * time.Millisecond)
		}
	}
	// Initially s should have no waiters.
	require.Equal(t, s.Stats().NumHadToWait, 0)
	if err := s.Acquire(ctx, 1); err != nil {
		t.Fatal(err)
	}
	// Still no waiters.
	require.Equal(t, s.Stats().NumHadToWait, 0)
	for i := 0; i < 10; i++ {
		go doAcquire(ctx)
	}
	assertNumWaitersSoon(10)
	s.Release(1)
	require.Recv(t, doneCh)
	go doAcquire(ctx)
	assertNumWaitersSoon(11)
	for i := 0; i < 10; i++ {
		s.Release(1)
		require.Recv(t, doneCh)
	}
	require.Equal(t, s.Stats().NumHadToWait, 11)
}

func TestConcurrentUpdatesAndAcquisitions(t *testing.T) {
	ctx := context.Background()
	var wg sync.WaitGroup
	const maxCap = 100
	s := NewSemaphore(maxCap)
	const N = 100
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			runtime.Gosched()
			newCap := rand.Int63n(maxCap-1) + 1
			s.UpdateCapacity(newCap)
		}()
	}
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			runtime.Gosched()
			n := rand.Int63n(maxCap)
			err := s.Acquire(ctx, n)
			runtime.Gosched()
			if err == nil {
				s.Release(n)
			}
		}()
	}
	wg.Wait()
	s.UpdateCapacity(maxCap)
	stats := s.Stats()
	require.Equal(t, stats.Capacity, 100)
	require.Equal(t, stats.Outstanding, 0)
}
