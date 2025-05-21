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

func Count[T Integer](count T) SafeString {
	if count < 0 {
		return "-" + Count[T](-count)
	}

	n, scaled := siUnit(uint64(count))
	digits := 0
	if scaled < 10 {
		digits = 1
	}
	return SafeString(fmt.Sprintf("%s%s", Float(scaled, digits), siUnits[n]))
}
