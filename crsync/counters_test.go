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

package crsync

import (
	"fmt"
	"math/rand/v2"
	"reflect"
	"runtime"
	"slices"
	"sync"
	"sync/atomic"
	"testing"
)

func TestCounter(t *testing.T) {
	c := MakeCounter()
	expect := func(expected int64) {
		if actual := c.Get(); actual != expected {
			t.Helper()
			t.Fatalf("expected %d, got %d", expected, actual)
		}
	}
	expect(0)
	c.Add(10)
	expect(10)
	c.Add(20)
	c.Add(-5)
	expect(25)
}

func TestCountersAll(t *testing.T) {
	c := MakeCounters(4)
	c.Add(0, 10)
	c.Add(1, 20)
	c.Add(0, 100)
	c.Add(2, 30)
	c.Add(1, 200)
	c.Add(3, 40)
	expected := []int64{110, 220, 30, 40}
	actual := slices.Collect(c.All())
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected %v, got %v", expected, actual)
	}
}

func TestCountersRand(t *testing.T) {
	numCounters := 1 + rand.IntN(100)
	c := MakeCounters(numCounters)
	numWorkers := 1 + rand.IntN(runtime.GOMAXPROCS(0)*10)
	allVals := make([][]int64, numWorkers)
	var wg sync.WaitGroup
	for i := range numWorkers {
		i := i
		vals := make([]int64, numCounters)
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range rand.IntN(1000) {
				k := rand.IntN(numCounters)
				v := rand.Int64N(1000)
				c.Add(k, v)
				vals[k] += v
			}
			allVals[i] = vals
		}()
	}
	wg.Wait()
	actual := slices.Collect(c.All())
	for k := range numCounters {
		var expected int64
		for i := range numWorkers {
			expected += allVals[i][k]
		}
		if actual[k] != expected {
			t.Fatalf("expected %d, got %d", expected, actual[k])
		}
	}
}

// BenchmarkCounters compares the performance of Counters against simple atomic
// counters, and against sharded counters with 4*P shards and random shard
// choice. There are two Counters versions (crsync and crync-cr), depending on
// whether the CockroachDB Go runtime (and cockroach_go tag) is used.
//
// # Benchmark results
//
// ## Apple M1 Pro (10 core)
//
// benchmark      simple       randshards   crsync        crsync-cr
// c=1/p=1-10     6.96ns ± 0%  9.69ns ± 0%  12.02ns ± 0%  7.34ns ± 1%
// c=1/p=4-10      169ns ±10%    53ns ± 4%     21ns ±39%    13ns ±26%
// c=1/p=10-10     752ns ± 3%   125ns ± 0%     51ns ±20%    51ns ± 7% *
// c=1/p=40-10    3.04µs ± 2%  0.67µs ±13%   0.26µs ±31%  0.29µs ±16%
// c=10/p=1-10    4.49ns ± 1%  9.72ns ± 0%  12.39ns ± 0%  4.95ns ± 0%
// c=10/p=4-10     147ns ± 5%    49ns ± 3%     27ns ±33%     6ns ± 1%
// c=10/p=10-10    790ns ± 6%   106ns ± 0%     47ns ±20%     8ns ± 4% *
// c=10/p=40-10   3.24µs ± 2%  0.61µs ±11%   0.24µs ± 9%  0.11µs ±28%
// c=100/p=1-10   4.33ns ± 0%  9.76ns ± 0%  12.41ns ± 0%  4.82ns ± 0%
// c=100/p=4-10   73.9ns ± 4%  46.0ns ± 5%   21.9ns ±22%   6.2ns ± 6%
// c=100/p=10-10   197ns ± 1%    94ns ±10%     53ns ±17%    11ns ± 1% *
// c=100/p=40-10   893ns ± 6%   524ns ± 7%    249ns ±19%   125ns ± 8%
// .                                              * one worker per core
//
// ## Intel(R) Xeon(R) CPU @ 2.80GH (24 core, n2-custom-24-32768 on GCE)
//
// benchmark      simple       randshards   crsync       crsync-cr
// c=1/p=1-24     14.1ns ± 0%  21.0ns ± 1%  37.7ns ± 0%  13.5ns ± 0%
// c=1/p=4-24     92.9ns ± 1%  50.1ns ± 1%  63.8ns ±29%  13.3ns ± 0%
// c=1/p=24-24     487ns ±18%  178ns ±105%   144ns ±39%    57ns ±60% *
// c=1/p=96-24    1.84µs ± 2%  0.59µs ± 3%  0.52µs ± 6%  0.29µs ± 7%
// c=10/p=1-24    13.8ns ± 0%  21.2ns ± 1%  38.0ns ± 1%  14.1ns ± 3%
// c=10/p=4-24    91.0ns ± 3%  48.6ns ± 1%  63.8ns ±16%  14.0ns ± 2%
// c=10/p=24-24    461ns ± 8%   176ns ±53%   146ns ±36%   110ns ±84% *
// c=10/p=96-24   1.79µs ± 1%  0.55µs ± 8%  0.52µs ± 6%  0.31µs ± 5%
// c=100/p=1-24   13.7ns ± 0%  22.0ns ± 2%  38.0ns ± 0%  14.1ns ±10%
// c=100/p=4-24   63.5ns ± 1%  46.4ns ± 2%  66.7ns ±30%  14.2ns ± 5%
// c=100/p=24-24   295ns ±27%    87ns ± 1%   121ns ±24%    44ns ±71% *
// c=100/p=96-24  1.11µs ± 2%  0.53µs ± 4%  0.52µs ± 8%  0.31µs ± 5%
// .                                             * one worker per core
func BenchmarkCounters(b *testing.B) {
	forEach := func(b *testing.B, fn func(b *testing.B, c, p int)) {
		for _, c := range []int{1, 10, 100} {
			for _, p := range []int{1, 4, runtime.GOMAXPROCS(0), 4 * runtime.GOMAXPROCS(0)} {
				b.Run(fmt.Sprintf("c=%d/p=%d", c, p), func(b *testing.B) {
					fn(b, c, p)
				})
			}
		}
	}

	// simple uses non-sharded atomic counters.
	b.Run("simple", func(b *testing.B) {
		forEach(b, func(b *testing.B, c, p int) {
			counters := make([]atomic.Int64, c)
			incCounter := func(counter int) {
				counters[counter].Add(1)
			}
			runCountersBenchmark(b, c, p, incCounter)
		})
	})

	// randshards uses a 4*N shards with random shard choice.
	b.Run("randshards", func(b *testing.B) {
		forEach(b, func(b *testing.B, c, p int) {
			counters := makeCounters(runtime.GOMAXPROCS(0)*4, c)
			incCounter := func(counter int) {
				shard := rand.Uint32N(counters.numShards)
				counters.counters[shard*counters.shardSize+uint32(counter)].Add(1)
			}
			runCountersBenchmark(b, c, p, incCounter)
		})
	})

	name := "crsync"
	if UsingCockroachGo {
		name += "-cr"
	}
	b.Run(name, func(b *testing.B) {
		forEach(b, func(b *testing.B, c, p int) {
			counters := MakeCounters(c)
			incCounter := func(counter int) {
				counters.Add(counter, 1)
			}
			runCountersBenchmark(b, c, p, incCounter)
		})
	})
}

func runCountersBenchmark(
	b *testing.B, numCounters, parallelism int, incCounter func(counter int),
) {
	const batchSize = 1000
	// Each element of ch corresponds to a batch of operations to be performed.
	ch := make(chan int, 1+b.N/batchSize)

	var wg sync.WaitGroup
	for range parallelism {
		wg.Add(1)
		go func() {
			defer wg.Done()

			rng := rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))
			for numOps := range ch {
				for range numOps {
					incCounter(rng.IntN(numCounters))
				}
			}
		}()
	}

	numOps := int64(b.N) * int64(parallelism)
	for i := int64(0); i < numOps; i += batchSize {
		ch <- int(min(batchSize, numOps-i))
	}
	close(ch)
	wg.Wait()
}
