package views_test

import (
	"bytes"
	"image"
	"image/png"
	"strings"
	"testing"

	"hastcode.com/site/internal/views"
)

// The hero card renderer. The builder tests prove the card goes through
// the hero pipeline; these check what is drawn, through the renderer's
// text hook (HeroCardText) and the PNG itself.

func iggy() views.HeroCard {
	return views.HeroCard{
		Repo:        "apache/iggy",
		Description: "Persistent message streaming platform written in Rust",
		Language:    "Rust",
		Stars:       2468,
		License:     "Apache-2.0",
		Release:     "server-0.5.0",
	}
}

func TestHeroCard_DrawsTheWordmarkRepositoryDescriptionLanguageStarsLicenseAndRelease(t *testing.T) {
	lines := views.HeroCardText(iggy())

	if len(lines) < 4 {
		t.Fatalf("card lines = %q, want the wordmark, repository, description and facts", lines)
	}
	if lines[0] != "hastcode" {
		t.Errorf("first line = %q, want the wordmark", lines[0])
	}
	if lines[1] != "apache/iggy" {
		t.Errorf("second line = %q, want the repository's full name", lines[1])
	}
	if got := strings.Join(lines[2:len(lines)-1], " "); got != "Persistent message streaming platform written in Rust" {
		t.Errorf("description lines = %q, want the whole description (it fits)", got)
	}
	if got := lines[len(lines)-1]; got != "Rust · 2,468 stars · Apache-2.0 · server-0.5.0" {
		t.Errorf("facts line = %q, want language, stars, license and latest release", got)
	}
}

func TestHeroCard_LongDescriptionIsCutWithAnEllipsisAfterTwoLines(t *testing.T) {
	c := iggy()
	c.Description = "Kubernetes Cluster API Provider AWS provides consistent deployment and day 2 operations of self-managed and EKS Kubernetes clusters on AWS, with more words than any card could hold."

	lines := views.HeroCardText(c)

	desc := lines[2 : len(lines)-1]
	if len(desc) != 2 {
		t.Fatalf("description takes %d lines %q, want 2", len(desc), desc)
	}
	if !strings.HasPrefix(desc[0], "Kubernetes Cluster API Provider AWS") || !strings.HasSuffix(desc[1], "…") {
		t.Errorf("description lines = %q, want the start of the description ending in an ellipsis", desc)
	}
	if strings.Contains(strings.Join(desc, " "), "any card could hold") {
		t.Errorf("description lines = %q, want the end cut", desc)
	}
}

func TestHeroCard_LeavesOutAbsentFactsAndRunesTheFontCannotDraw(t *testing.T) {
	c := views.HeroCard{Repo: "php/frankenphp", Description: "🚀 The modern PHP app server ★", Stars: 9001}

	lines := views.HeroCardText(c)

	want := []string{"hastcode", "php/frankenphp", "The modern PHP app server", "9,001 stars"}
	if strings.Join(lines, "|") != strings.Join(want, "|") {
		t.Errorf("card lines = %q, want %q (no language, license or release; no emoji)", lines, want)
	}
	if lines := views.HeroCardText(views.HeroCard{Repo: "php/frankenphp"}); strings.Join(lines, "|") != "hastcode|php/frankenphp" {
		t.Errorf("card lines with no facts = %q, want only the wordmark and repository", lines)
	}
}

func TestHeroCard_AltIsTheRepositoryAndTheFullDescription(t *testing.T) {
	if got := iggy().Alt(); got != "apache/iggy: Persistent message streaming platform written in Rust" {
		t.Errorf("Alt() = %q", got)
	}
	if got := (views.HeroCard{Repo: "apache/iggy"}).Alt(); got != "apache/iggy" {
		t.Errorf("Alt() without a description = %q, want the repository alone", got)
	}
}

func TestHeroCard_PNGIs1600x900AndTheSameBytesEveryTime(t *testing.T) {
	data := views.HeroCardPNG(iggy())

	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	if b := img.Bounds(); b.Dx() != 1600 || b.Dy() != 900 {
		t.Errorf("card is %dx%d, want 1600x900", b.Dx(), b.Dy())
	}
	if !bytes.Equal(data, views.HeroCardPNG(iggy())) {
		t.Error("two renders of the same card differ")
	}
}

// Spec: the hero pipeline cuts 16:9, 4:3, 1:1 and 1200x630 crops from the
// card. The 1:1 crop keeps the middle 900 px and the OG crop trims 30 px
// top and bottom, so everything but the bottom band must sit inside that.
func TestHeroCard_KeepsEveryLineInsideTheSquareAndOpenGraphCrops(t *testing.T) {
	long := views.HeroCard{
		Repo:        "kubernetes-sigs/cluster-api-provider-aws",
		Description: "Kubernetes Cluster API Provider AWS provides consistent deployment and day 2 operations of self-managed and EKS Kubernetes clusters on AWS.",
		Language:    "Jupyter Notebook",
		Stars:       1234567,
		License:     "GPL-3.0-or-later",
		Release:     "v2.10.0-beta.1",
	}
	for name, c := range map[string]views.HeroCard{"short": iggy(), "long": long} {
		t.Run(name, func(t *testing.T) {
			img, err := png.Decode(bytes.NewReader(views.HeroCardPNG(c)))
			if err != nil {
				t.Fatal(err)
			}
			minX, minY, maxX, maxY := 1600, 900, -1, -1
			for y := 0; y < 860; y++ { // above the accent band
				for x := range 1600 {
					if !isBackground(img, x, y) {
						minX, minY, maxX, maxY = min(minX, x), min(minY, y), max(maxX, x), max(maxY, y)
					}
				}
			}
			if maxX < 0 {
				t.Fatal("nothing is drawn above the band")
			}
			if minX < 350 || maxX >= 1250 {
				t.Errorf("ink spans x %d..%d, want within 350..1249 (the 1:1 crop)", minX, maxX)
			}
			if minY < 30 || maxY >= 870 {
				t.Errorf("ink spans y %d..%d, want within 30..869 (the 1200x630 crop)", minY, maxY)
			}
		})
	}
}

func isBackground(img image.Image, x, y int) bool {
	r, g, b, _ := img.At(x, y).RGBA()
	return r>>8 == 0xFF && g>>8 == 0xFF && b>>8 == 0xFF
}
