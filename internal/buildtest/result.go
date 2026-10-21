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

// RequireErrorRule fails the test unless res has an error attributed to file
// for the named rule whose message contains want.
func RequireErrorRule(tb testing.TB, res builder.Result, file, rule, want string) {
	tb.Helper()
	if !hasRuleProblem(res.Errors, file, rule, want) {
		tb.Fatalf("want a %q error for %s containing %q, got:\n%s", rule, file, want, res)
	}
}

// RequireWarningRule is RequireErrorRule for warnings.
func RequireWarningRule(tb testing.TB, res builder.Result, file, rule, want string) {
	tb.Helper()
	if !hasRuleProblem(res.Warnings, file, rule, want) {
		tb.Fatalf("want a %q warning for %s containing %q, got:\n%s", rule, file, want, res)
	}
}

// RequireNoProblemRule fails the test if res has an error or a warning for
// file under the named rule.
func RequireNoProblemRule(tb testing.TB, res builder.Result, file, rule string) {
	tb.Helper()
	if hasRuleProblem(res.Errors, file, rule, "") || hasRuleProblem(res.Warnings, file, rule, "") {
		tb.Fatalf("want no %q problem for %s, got:\n%s", rule, file, res)
	}
}

func hasRuleProblem(ps []builder.Problem, file, rule, want string) bool {
	for _, p := range ps {
		if p.File == file && p.Rule == rule && strings.Contains(p.Message, want) {
			return true
		}
	}
	return false
}
