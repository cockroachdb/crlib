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

package crstrings

import (
	"fmt"
	"strings"
	"testing"
)

type num int

func (n num) String() string {
	return fmt.Sprintf("%03d", int(n))
}

func TestJoinStringers(t *testing.T) {
	nums := []num{0, 1, 2, 3}
	expect(t, "", JoinStringers(", ", nums[:0]...))
	expect(t, "000", JoinStringers(", ", nums[0]))
	expect(t, "000, 001", JoinStringers(", ", nums[0], nums[1]))
	expect(t, "000, 001, 002, 003", JoinStringers(", ", nums...))
}

func TestMapAndJoin(t *testing.T) {
	nums := []int{0, 1, 2, 3}
	fn := func(n int) string {
		return fmt.Sprintf("%d", n)
	}
	expect(t, "", MapAndJoin(fn, ", ", nums[:0]...))
	expect(t, "0", MapAndJoin(fn, ", ", nums[0]))
	expect(t, "0, 1", MapAndJoin(fn, ", ", nums[0], nums[1]))
	expect(t, "0, 1, 2, 3", MapAndJoin(fn, ", ", nums...))
}

func expect(t *testing.T, expected, actual string) {
	t.Helper()
	if actual != expected {
		t.Errorf("expected %q got %q", expected, actual)
	}
}

func TestIf(t *testing.T) {
	expect(t, "", If(false, "true"))
	expect(t, "true", If(true, "true"))
}

func TestIfElse(t *testing.T) {
	expect(t, "false", IfElse(false, "true", "false"))
	expect(t, "true", IfElse(true, "true", "false"))
}

func TestPrependIfNotEmpty(t *testing.T) {
	expect(t, "<prefix>a", PrependIfNotEmpty("<prefix>", "a"))
	expect(t, "", PrependIfNotEmpty("<prefix>", ""))
}

func TestAppendIfNotEmpty(t *testing.T) {
	expect(t, "a<suffix>", AppendIfNotEmpty("a", "<suffix>"))
	expect(t, "", AppendIfNotEmpty("", "<suffix>"))
}

func TestFilterEmpty(t *testing.T) {
	s := []string{"a", "", "b", "", "c", ""}
	expect(t, "a,b,c", strings.Join(FilterEmpty(s), ","))
}

func TestLines(t *testing.T) {
	expect(t, `["a" "b" "c"]`, fmt.Sprintf("%q", Lines("a\nb\nc")))
	expect(t, `["a" "b" "c"]`, fmt.Sprintf("%q", Lines("a\nb\nc\n")))
	expect(t, `["a" "b" "c" ""]`, fmt.Sprintf("%q", Lines("a\nb\nc\n\n")))
	expect(t, `["" "a" "b" "c"]`, fmt.Sprintf("%q", Lines("\na\nb\nc\n")))
	expect(t, `[]`, fmt.Sprintf("%q", Lines("")))
	expect(t, `[]`, fmt.Sprintf("%q", Lines("\n")))
}
