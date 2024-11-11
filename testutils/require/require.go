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

package require

import "fmt"

// TB is an interface common to *testing.T and *testing.B.
type TB interface {
	Error(args ...any)
	Errorf(format string, args ...any)
	Fatal(args ...any)
	Fatalf(format string, args ...any)
	Helper()
	Log(args ...any)
	Logf(format string, args ...any)
}

// withMsg implements TB and prepends some information to all logs or error
// messages.
type withMsg struct {
	TB

	msg string
}

func (w *withMsg) Error(args ...any) {
	w.TB.Helper()
	w.TB.Errorf("%s: %s", w.msg, fmt.Sprint(args...))
}
func (w *withMsg) Errorf(format string, args ...any) {
	w.TB.Helper()
	w.TB.Errorf("%s: %s", w.msg, fmt.Sprintf(format, args...))
}

func (w *withMsg) Fatal(args ...any) {
	w.TB.Helper()
	w.TB.Fatalf("%s: %s", w.msg, fmt.Sprint(args...))
}

func (w *withMsg) Fatalf(format string, args ...any) {
	w.TB.Helper()
	w.TB.Fatalf("%s: %s", w.msg, fmt.Sprintf(format, args...))
}

func (w *withMsg) Log(args ...any) {
	w.TB.Helper()
	w.TB.Logf("%s: %s", w.msg, fmt.Sprint(args...))
}

func (w *withMsg) Logf(format string, args ...any) {
	w.TB.Helper()
	w.TB.Logf("%s: %s", w.msg, fmt.Sprintf(format, args...))
}

// WithMsg returns a TB that can be used with assertions and logs which
// prepends a message to any log or error message.
//
// Example:
//
//	{
//	  t := require.WithMsg(t, "n=", n)
//	  require.Equal(t, a, b)
//	  require.LT(t, c, d)
//	}
//
// A failure message would look like:
//
//	n=5: expected 6 == 7
func WithMsg(tb TB, args ...any) TB {
	return &withMsg{TB: tb, msg: fmt.Sprint(args...)}
}

// WithMsgf returns a TB that can be used with assertions and logs which
// prepends a message to any log or error message.
//
// Example:
//
//	{
//	  t := require.WithMsgf(t, "n=%d", n)
//	  require.Equal(t, a, b)
//	  require.LT(t, c, d)
//	}
//
// A failure message would look like:
//
//	n=5: expected 6 == 7
func WithMsgf(tb TB, format string, args ...any) TB {
	return &withMsg{TB: tb, msg: fmt.Sprintf(format, args...)}
}
