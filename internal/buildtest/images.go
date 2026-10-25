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
)

// Image is an image file found in build output, measured by fully decoding it.
type Image struct {
	// URL is the site path the file is served at ("/assets/...").
	URL string
	// Format is the decoder's name: "avif", "jpeg" or "png".
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
// static/ dir. No WebP decoder is registered, so a .webp file fails the
// test.
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

// AVIFChroma reads the chroma subsampling an AVIF file written for urlPath
// declares in its AV1 codec configuration box (av1C, AV1-ISOBMFF 2.3.3):
// "4:4:4", "4:2:2", "4:2:0" or "monochrome".
func AVIFChroma(tb testing.TB, outDir, urlPath string) string {
	tb.Helper()
	data, err := os.ReadFile(StaticPath(outDir, urlPath))
	if err != nil {
		tb.Fatalf("image %s: %v", urlPath, err)
	}
	i := bytes.Index(data, []byte("av1C"))
	if i < 0 || i+7 > len(data) {
		tb.Fatalf("image %s has no av1C box", urlPath)
	}
	// After the box type: marker+version, profile+level, then
	// tier|high_bitdepth|twelve_bit|monochrome|subsampling_x|subsampling_y|position(2).
	flags := data[i+6]
	mono, sx, sy := flags&0x10 != 0, flags&0x08 != 0, flags&0x04 != 0
	switch {
	case mono:
		return "monochrome"
	case sx && sy:
		return "4:2:0"
	case sx:
		return "4:2:2"
	default:
		return "4:4:4"
	}
}

// FileSize is the length in bytes of the file written for urlPath. The test
// fails if it is missing.
// ContentSize is the size of the page at urlPath without its inline
// stylesheet (the <style> element in <head>), which the 50 KB gate does not
// count (ADR-0006).
func ContentSize(tb testing.TB, outDir, urlPath string) int64 {
	tb.Helper()
	data, err := os.ReadFile(StaticPath(outDir, urlPath))
	if err != nil {
		tb.Fatalf("file %s: %v", urlPath, err)
	}
	s := string(data)
	if i := strings.Index(s, "<style>"); i >= 0 {
		if j := strings.Index(s[i:], "</style>"); j >= 0 {
			return int64(len(s) - (j + len("</style>")))
		}
	}
	return int64(len(s))
}

func FileSize(tb testing.TB, outDir, urlPath string) int64 {
	tb.Helper()
	info, err := os.Stat(StaticPath(outDir, urlPath))
	if err != nil {
		tb.Fatalf("file %s: %v", urlPath, err)
	}
	return info.Size()
}
