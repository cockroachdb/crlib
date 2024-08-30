GO := go
PKG := ./...
GOFLAGS :=
STRESSFLAGS :=
TAGS := crlib_invariants
TESTS := .

.PHONY: gen-bazel
gen-bazel:
	@echo "Generating WORKSPACE"
	@echo 'workspace(name = "com_github_cockroachdb_crlib")' > WORKSPACE
	@echo 'Running gazelle...'
	${GO} run github.com/bazelbuild/bazel-gazelle/cmd/gazelle@v0.37.0 update --go_prefix=github.com/cockroachdb/crlib --repo_root=.
	@echo 'You should now be able to build Cockroach using:'
	@echo '  ./dev build short -- --override_repository=com_github_cockroachdb_crlib=${CURDIR}'

.PHONY: clean-bazel
clean-bazel:
	git clean -dxf WORKSPACE BUILD.bazel '**/BUILD.bazel'

.PHONY: test
test:
	${GO} test -tags '$(TAGS)' ${testflags} -run ${TESTS} ${PKG}

.PHONY: lint
lint:
	${GO} test -tags '$(TAGS)' -run ${TESTS} ./internal/lint

.PHONY: stress
stress:
	${GO} test -tags '$(TAGS)' -exec 'stress -p 2 -maxruns 1000' -v -run ${TESTS} ${PKG}
