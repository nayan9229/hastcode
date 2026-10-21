package buildtest

import (
	"strings"
	"testing"

	"hastcode.com/site/builder"
)

// RequireError fails the test unless res has an error attributed to file
// (content-relative, slash-separated) whose text contains want.
func RequireError(tb testing.TB, res builder.Result, file, want string) {
	tb.Helper()
	if !hasProblem(res.Errors, file, want) {
		tb.Fatalf("want an error for %s containing %q, got:\n%s", file, want, res)
	}
}

// RequireWarning is RequireError for warnings.
func RequireWarning(tb testing.TB, res builder.Result, file, want string) {
	tb.Helper()
	if !hasProblem(res.Warnings, file, want) {
		tb.Fatalf("want a warning for %s containing %q, got:\n%s", file, want, res)
	}
}

func hasProblem(ps []builder.Problem, file, want string) bool {
	for _, p := range ps {
		if p.File == file && strings.Contains(p.Error(), want) {
			return true
		}
	}
	return false
}
