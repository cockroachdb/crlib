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

package crhumanize

import "fmt"

func ExampleBytes() {
	fmt.Println(Bytes(120_000))
	fmt.Println(Bytes(1024*10000, Exact))
	fmt.Println(Bytes(950_000, Compact, OmitI))
	// Output:
	// 117 KiB
	// 10,000 KiB
	// 928KB
}

func ExampleBytesPerSec() {
	fmt.Println(BytesPerSec(120_000))
	fmt.Println(BytesPerSec(1024*10000, Exact))
	fmt.Println(BytesPerSec(950_000, Compact, OmitI))
	// Output:
	// 117 KiB/s
	// 10,000 KiB/s
	// 928KB/s
}

func ExampleCount() {
	fmt.Println(Count(120_000))
	fmt.Println(Count(1024*10000, Exact))
	fmt.Println(Count(950_000, Compact, OmitI))
	// Output:
	// 120 K
	// 10,240 K
	// 950K
}
