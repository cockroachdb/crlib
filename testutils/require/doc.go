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

/*
Package require implements convenience wrappers around checking conditions and
failing tests.

The interface is inspired from `github.com/stretchr/testify/require` but the
implementation is simpler and uses generics. The benefit of generics is that we
don't have to add casts to make the types match, e.g.
[require.Equal](t, uint32Var, 2).

Failed assertions result in a t.Fatal() call.

# Equality

  - [require.Equal]
  - [require.NotEqual]
  - [require.True]
  - [require.False]

# Comparisons

  - [require.LT]
  - [require.LE]
  - [require.GT]
  - [require.GE]

# Channels

  - [require.Recv], [require.RecvWithin]
  - [require.NoRecv], [require.NoRecvWithin]

# Errors
  - [require.NoError]
  - [require.NoError1], [require.NoError2]

# Including info in error messages
  - [require.WithMsg], [require.WithMsgf]
*/
package require
