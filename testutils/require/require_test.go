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

package require_test

import (
	"testing"

	"github.com/cockroachdb/crlib/testutils/require"
)

func TestWithMsg(t *testing.T) {
	t2 := require.WithMsg(t, "foo")
	// foo: hello1
	t2.Logf("hello%d", 1)

	// 1.2: hello2
	t2 = require.WithMsgf(t, "%d.%d", 1, 2)
	t2.Log("hello2")

	// 1.2: bar: hello3
	t3 := require.WithMsgf(t2, "bar")
	t3.Log("hello3")
}
