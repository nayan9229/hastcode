package buildtest

import (
	"bytes"
	"image"
	_ "image/jpeg" // decode JPEG output
	_ "image/png"  // decode PNG output
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	_ "github.com/gen2brain/avif" // decode AVIF output
	_ "github.com/gen2brain/webp" // decode WebP output
)

// Image is an image file found in build output, measured by fully decoding it.
type Image struct {
	// URL is the site path the file is served at ("/assets/...").
	URL string
	// Format is the decoder's name: "avif", "webp", "jpeg" or "png".
	Format string
	Width  int
	Height int
}

// DecodeImage decodes the file written for urlPath and reports its format
// and pixel size. The test fails if the file is missing or does not decode.
func DecodeImage(tb testing.TB, outDir, urlPath string) Image {
	tb.Helper()
	data, err := os.ReadFile(StaticPath(outDir, urlPath))
	if err != nil {
		tb.Fatalf("image %s: %v", urlPath, err)
	}
	img, format, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		tb.Fatalf("image %s does not decode: %v", urlPath, err)
	}
	b := img.Bounds()
	return Image{URL: urlPath, Format: format, Width: b.Dx(), Height: b.Dy()}
}

// Images decodes every .avif, .webp, .jpg and .png file under the output's
// static/ dir.
func Images(tb testing.TB, outDir string) []Image {
	tb.Helper()
	root := filepath.Join(outDir, "static")
	var imgs []Image
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		switch filepath.Ext(path) {
		case ".avif", ".webp", ".jpg", ".jpeg", ".png":
		default:
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		imgs = append(imgs, DecodeImage(tb, outDir, "/"+filepath.ToSlash(rel)))
		return nil
	})
	if err != nil {
		tb.Fatal(err)
	}
	return imgs
}

// SrcsetEntry is one candidate of a srcset attribute with a width descriptor.
type SrcsetEntry struct {
	URL   string
	Width int
}

// Srcset parses a srcset attribute of "url Nw" candidates. The test fails on
// a candidate without a width descriptor.
func Srcset(tb testing.TB, srcset string) []SrcsetEntry {
	tb.Helper()
	var out []SrcsetEntry
	for _, cand := range strings.Split(srcset, ",") {
		fields := strings.Fields(cand)
		if len(fields) != 2 || !strings.HasSuffix(fields[1], "w") {
			tb.Fatalf("srcset candidate %q: want \"url <width>w\"", cand)
		}
		w, err := strconv.Atoi(strings.TrimSuffix(fields[1], "w"))
		if err != nil {
			tb.Fatalf("srcset candidate %q: %v", cand, err)
		}
		out = append(out, SrcsetEntry{URL: fields[0], Width: w})
	}
	return out
}
