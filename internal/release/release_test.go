package release_test

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"hastcode.com/site/internal/release"
)

// piece writes a minimal Piece file under dir (content/explainers/<vertical>/
// or content/briefs/) with the given published timestamp and draft flag.
func piece(t *testing.T, dir, rel, published string, draft bool) string {
	t.Helper()
	fm := "---\ntitle: \"T " + filepath.Base(rel) + "\"\n"
	if strings.HasPrefix(rel, "explainers") {
		fm += "kind: concept\n"
	}
	fm += "published: \"" + published + "\"\nauthor: ada\n"
	if strings.HasPrefix(rel, "explainers") {
		fm += "reviewed: \"" + published + "\"\n"
	}
	if draft {
		fm += "draft: true\n"
	}
	fm += "---\n\nBody.\n"
	p := filepath.Join(dir, rel)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(fm), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func clock(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}

func TestGapDays_CountsLiveAndDraftedPiecesPerPublicationDayInsideTheLookBack(t *testing.T) {
	dir := t.TempDir()
	piece(t, dir, "explainers/distributed-systems/a.md", "2026-11-12T09:00:00Z", false)
	piece(t, dir, "briefs/b.md", "2026-11-12T10:00:00Z", false)
	piece(t, dir, "explainers/databases-storage/c.md", "2026-11-13T23:30:00-05:00", false) // 13 Nov in its own offset
	piece(t, dir, "explainers/apis-networking/d.md", "2026-11-14T09:00:00Z", true)         // a draft already fills one slot
	piece(t, dir, "explainers/apis-networking/old.md", "2026-11-10T09:00:00Z", false)      // outside the look-back

	got, err := release.GapDays(dir, clock("2026-11-15T10:00:00Z"), 3)
	if err != nil {
		t.Fatal(err)
	}
	want := []release.Gap{{Day: "2026-11-13", Slots: 1}, {Day: "2026-11-14", Slots: 1}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GapDays = %+v, want %+v", got, want)
	}
}

func TestGapDays_NeverBeforeTheCapStartedAndNeverToday(t *testing.T) {
	dir := t.TempDir()
	got, err := release.GapDays(dir, clock("2026-10-31T08:00:00Z"), 3)
	if err != nil {
		t.Fatal(err)
	}
	want := []release.Gap{{Day: "2026-10-30", Slots: 2}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GapDays = %+v, want only the cap's first day %+v", got, want)
	}
}

func TestRelease_ReleasesEveryDraftKeepingItsPublicationDay(t *testing.T) {
	dir := t.TempDir()
	e := piece(t, dir, "explainers/distributed-systems/a.md", "2026-11-13T09:17:00Z", true)
	b := piece(t, dir, "briefs/b.md", "2026-11-15T06:00:00+05:30", true)
	piece(t, dir, "explainers/databases-storage/live.md", "2026-11-12T09:00:00Z", false)
	now := clock("2026-11-15T20:00:00Z") // already 16 Nov in +05:30

	plan, err := release.NewPlan(dir, nil, now)
	if err != nil {
		t.Fatal(err)
	}
	var paths []string
	for _, p := range plan.Releases {
		paths = append(paths, p.Path)
	}
	if want := []string{"briefs/b.md", "explainers/distributed-systems/a.md"}; !reflect.DeepEqual(paths, want) {
		t.Fatalf("plan releases %v, want %v", paths, want)
	}
	if err := plan.Apply(); err != nil {
		t.Fatal(err)
	}

	explainer, _ := os.ReadFile(e)
	for _, want := range []string{"published: \"2026-11-13T09:17:00Z\"\nverified: 2026-11-15\n", "reviewed: \"2026-11-15T20:00:00Z\"\n", "author: ada\n", "---\n\nBody.\n"} {
		if !strings.Contains(string(explainer), want) {
			t.Errorf("released Explainer lacks %q:\n%s", want, explainer)
		}
	}
	if strings.Contains(string(explainer), "draft: true") {
		t.Error("released Explainer still carries draft: true")
	}
	brief, _ := os.ReadFile(b)
	for _, want := range []string{"published: \"2026-11-15T06:00:00+05:30\"\nverified: 2026-11-16\n"} {
		if !strings.Contains(string(brief), want) {
			t.Errorf("released Brief lacks %q (verified is the day in the Piece's own offset):\n%s", want, brief)
		}
	}
	if strings.Contains(string(brief), "draft: true") || strings.Contains(string(brief), "reviewed:") {
		t.Errorf("released Brief carries draft or reviewed:\n%s", brief)
	}
}

func TestRelease_RefusesADayOverTheCapNamingTheFiles(t *testing.T) {
	dir := t.TempDir()
	piece(t, dir, "explainers/distributed-systems/x.md", "2026-11-13T09:00:00Z", false)
	piece(t, dir, "briefs/y.md", "2026-11-13T10:00:00Z", false)
	piece(t, dir, "explainers/databases-storage/z.md", "2026-11-13T11:00:00Z", true)

	_, err := release.NewPlan(dir, nil, clock("2026-11-15T10:00:00Z"))
	if err == nil || !strings.Contains(err.Error(), "2026-11-13") || !strings.Contains(err.Error(), "explainers/databases-storage/z.md") {
		t.Fatalf("Plan error = %v, want a refusal naming 2026-11-13 and z.md", err)
	}
}

func TestRelease_FilesLimitsTheReleaseAndMustBeDrafts(t *testing.T) {
	dir := t.TempDir()
	a := piece(t, dir, "explainers/distributed-systems/a.md", "2026-11-13T09:00:00Z", true)
	b := piece(t, dir, "briefs/b.md", "2026-11-14T09:00:00Z", true)
	live := piece(t, dir, "briefs/live.md", "2026-11-12T09:00:00Z", false)
	now := clock("2026-11-15T10:00:00Z")

	plan, err := release.NewPlan(dir, []string{b}, now)
	if err != nil {
		t.Fatal(err)
	}
	if err := plan.Apply(); err != nil {
		t.Fatal(err)
	}
	if data, _ := os.ReadFile(a); !strings.Contains(string(data), "draft: true") {
		t.Error("a.md was released though FILES named only b.md")
	}
	if data, _ := os.ReadFile(b); strings.Contains(string(data), "draft: true") {
		t.Error("b.md was not released")
	}
	if _, err := release.NewPlan(dir, []string{live}, now); err == nil || !strings.Contains(err.Error(), "briefs/live.md") {
		t.Errorf("Plan with a live file = %v, want an error naming it", err)
	}
}

func TestGapDays_TodayIsTheClocksOwnDayNotUTC(t *testing.T) {
	dir := t.TempDir()
	piece(t, dir, "explainers/distributed-systems/a.md", "2026-11-15T09:00:00+05:30", false)
	piece(t, dir, "briefs/b.md", "2026-11-15T10:00:00+05:30", false)

	// 03:00 on 16 Nov in +05:30 is still 15 Nov in UTC; the look-back
	// must count from 16 Nov, so 15 Nov (full) is checked and 12 Nov is not.
	got, err := release.GapDays(dir, clock("2026-11-16T03:00:00+05:30"), 3)
	if err != nil {
		t.Fatal(err)
	}
	want := []release.Gap{{Day: "2026-11-13", Slots: 2}, {Day: "2026-11-14", Slots: 2}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GapDays = %+v, want %+v", got, want)
	}
}

func TestRelease_ADraftAlreadyCarryingVerifiedGetsExactlyOneVerifiedLine(t *testing.T) {
	dir := t.TempDir()
	path := piece(t, dir, "explainers/distributed-systems/a.md", "2026-11-13T09:00:00Z", true)
	data, _ := os.ReadFile(path)
	data = []byte(strings.Replace(string(data), "published: \"2026-11-13T09:00:00Z\"\n", "published: \"2026-11-13T09:00:00Z\"\nverified: 2026-11-13\n", 1))
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}

	plan, err := release.NewPlan(dir, []string{"explainers/distributed-systems/a.md"}, clock("2026-11-15T10:00:00Z")) // content-relative form
	if err != nil {
		t.Fatal(err)
	}
	if err := plan.Apply(); err != nil {
		t.Fatal(err)
	}
	got, _ := os.ReadFile(path)
	if n := strings.Count(string(got), "verified:"); n != 1 || !strings.Contains(string(got), "verified: 2026-11-15\n") {
		t.Errorf("released file has %d verified lines, want one reading 2026-11-15:\n%s", n, got)
	}
}

func TestRedate_MovesAnExplainersDayKeepingTheClockTimeAndRefusesABrief(t *testing.T) {
	dir := t.TempDir()
	e := piece(t, dir, "explainers/distributed-systems/a.md", "2026-11-15T07:05:49+05:30", true)
	b := piece(t, dir, "briefs/b.md", "2026-11-15T07:05:49Z", true)

	if err := release.Redate(dir, e, "2026-11-13"); err != nil {
		t.Fatal(err)
	}
	got, _ := os.ReadFile(e)
	for _, want := range []string{"published: \"2026-11-13T07:05:49+05:30\"\n", "reviewed: \"2026-11-13T07:05:49+05:30\"\n", "draft: true\n"} {
		if !strings.Contains(string(got), want) {
			t.Errorf("re-dated Explainer lacks %q:\n%s", want, got)
		}
	}
	if err := release.Redate(dir, b, "2026-11-13"); err == nil || !strings.Contains(err.Error(), "Brief") {
		t.Errorf("Redate on a Brief = %v, want a refusal naming Briefs", err)
	}
}
