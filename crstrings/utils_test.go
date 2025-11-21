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

	"github.com/cockroachdb/crlib/testutils/require"
)

type num int

func (n num) String() string {
	return fmt.Sprintf("%03d", int(n))
}

func TestJoinStringers(t *testing.T) {
	nums := []num{0, 1, 2, 3}
	require.Equal(t, "", JoinStringers(", ", nums[:0]...))
	require.Equal(t, "000", JoinStringers(", ", nums[0]))
	require.Equal(t, "000, 001", JoinStringers(", ", nums[0], nums[1]))
	require.Equal(t, "000, 001, 002, 003", JoinStringers(", ", nums...))
}

func TestMapAndJoin(t *testing.T) {
	nums := []int{0, 1, 2, 3}
	fn := func(n int) string {
		return fmt.Sprintf("%d", n)
	}
	require.Equal(t, "", MapAndJoin(fn, ", ", nums[:0]...))
	require.Equal(t, "0", MapAndJoin(fn, ", ", nums[0]))
	require.Equal(t, "0, 1", MapAndJoin(fn, ", ", nums[0], nums[1]))
	require.Equal(t, "0, 1, 2, 3", MapAndJoin(fn, ", ", nums...))
}

func TestIf(t *testing.T) {
	require.Equal(t, "", If(false, "true"))
	require.Equal(t, "true", If(true, "true"))
}

func TestIfElse(t *testing.T) {
	require.Equal(t, "false", IfElse(false, "true", "false"))
	require.Equal(t, "true", IfElse(true, "true", "false"))
}

func TestWithSep(t *testing.T) {
	require.Equal(t, "a,b", WithSep("a", ",", "b"))
	require.Equal(t, "a", WithSep("a", ",", ""))
	require.Equal(t, "b", WithSep("", ",", "b"))
}

func TestFilterEmpty(t *testing.T) {
	s := []string{"a", "", "b", "", "c", ""}
	require.Equal(t, "a,b,c", strings.Join(FilterEmpty(s), ","))
}

func TestLines(t *testing.T) {
	require.Equal(t, `["a" "b" "c"]`, fmt.Sprintf("%q", Lines("a\nb\nc")))
	require.Equal(t, `["a" "b" "c"]`, fmt.Sprintf("%q", Lines("a\nb\nc\n")))
	require.Equal(t, `["a" "b" "c" ""]`, fmt.Sprintf("%q", Lines("a\nb\nc\n\n")))
	require.Equal(t, `["" "a" "b" "c"]`, fmt.Sprintf("%q", Lines("\na\nb\nc\n")))
	require.Equal(t, `[]`, fmt.Sprintf("%q", Lines("")))
	require.Equal(t, `[]`, fmt.Sprintf("%q", Lines("\n")))
}

func TestIndent(t *testing.T) {
	testCases := [][2]string{
		{"", ""},
		{"foo", "--foo"},
		{"foo\n", "--foo\n"},
		{"foo\n\n", "--foo\n--\n"},
		{"foo\nbar", "--foo\n--bar"},
		{"foo\nbar\n", "--foo\n--bar\n"},
		{"foo\n\nbar\n", "--foo\n--\n--bar\n"},
	}
	for _, tc := range testCases {
		require.Equal(t, tc[1], Indent("--", tc[0]))
	}
}

func TestUnwrapText(t *testing.T) {
	expected := "This is a single line string. It looks fine."

	require.Equal(t, expected, UnwrapText(`This
is a single line string.
It looks fine.`))

	require.Equal(t, expected, UnwrapText(`


		This
		is a single line string.
		It looks
		fine.

	`))

	expected = "This is a paragraph that is wrapped on multiple lines.\n\nThis is another paragraph."
	require.Equal(t, expected, UnwrapText(`This is a paragraph that
		is wrapped on multiple lines.
	
		This is another
		paragraph.`))

	require.Equal(t, expected, UnwrapText(`
	
		This is a paragraph that
		is wrapped on multiple lines.

		This is another paragraph.
	
	`))
}
