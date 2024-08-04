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

package lint

import (
	"go/build"
	"os/exec"
	"strings"
	"testing"
)

const (
	cmdGo       = "go"
	staticcheck = "honnef.co/go/tools/cmd/staticcheck"
	crlfmt      = "github.com/cockroachdb/crlfmt"
)

func dirCmd(t *testing.T, dir string, name string, args ...string) string {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("error running %s %s: %s\n%s\n", name, strings.Join(args, "\n"), err, string(out))
	}
	return strings.TrimSpace(string(out))
}

func installTool(t *testing.T, path string) {
	cmd := exec.Command(cmdGo, "install", "-C", "../devtools", path)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("cannot install %q: %v\n%s\n", path, err, out)
	}
}

func TestLint(t *testing.T) {
	const root = "github.com/cockroachdb/crlib"

	pkg, err := build.Import(root, "../..", 0)
	if err != nil {
		t.Fatal(err)
	}

	pkgs := strings.Split(
		dirCmd(t, pkg.Dir, "go", "list", "./..."), "\n",
	)

	// TestGoVet is the fastest check that verifies that all files build, so we
	// want to run it first (and not in parallel).
	t.Run("TestGoVet", func(t *testing.T) {
		out := dirCmd(t, pkg.Dir, "go", "vet", "-all", "./...")
		for _, l := range out {
			t.Error(l)
		}
	})

	t.Run("TestStaticcheck", func(t *testing.T) {
		installTool(t, staticcheck)
		t.Parallel()

		out := dirCmd(t, pkg.Dir, "staticcheck", pkgs...)
		if out != "" {
			t.Errorf("staticcheck:\n%s\n", out)
		}
	})

	t.Run("TestCrlfmt", func(t *testing.T) {
		installTool(t, crlfmt)
		t.Parallel()

		args := []string{"-fast", "-tab", "2", "."}
		out := dirCmd(t, pkg.Dir, "crlfmt", args...)
		if out != "" {
			t.Errorf("crlfmt:\n%s\n", out)
		}

		if t.Failed() {
			reWriteCmd := []string{"crlfmt", "-w"}
			reWriteCmd = append(reWriteCmd, args...)
			t.Logf("run the following to fix your formatting:\n"+
				"\n%s\n\n"+
				"Don't forget to add amend the result to the correct commits.",
				strings.Join(reWriteCmd, " "),
			)
		}
	})
}
