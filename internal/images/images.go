// Package images resizes, crops and encodes raster images (AVIF, WebP,
// JPEG) with an on-disk cache keyed on the source bytes and encoder
// settings. It knows nothing about heroes, pages or URLs: the builder
// decides which variants a hero needs and where they are served.
//
// AVIF and WebP use gen2brain/avif and gen2brain/webp (libavif and libwebp
// as WASM, no cgo). Build with CGO_ENABLED=0 -tags nodynamic so they never
// load a system library and output is identical on every machine.
package images

import (
	"bytes"
	"cmp"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png" // decode PNG sources
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/gen2brain/avif"
	"github.com/gen2brain/webp"
	"golang.org/x/image/draw"
	"golang.org/x/sync/errgroup"
)

// cacheVersion is part of every cache key. Bump it when resizing or
// encoding changes in a way the settings below do not capture.
const cacheVersion = "1"

// Format is an output encoding.
type Format string

const (
	AVIF Format = "avif"
	WebP Format = "webp"
	JPEG Format = "jpeg"
)

// Ext is the file extension for the format, with the dot.
func (f Format) Ext() string {
	if f == JPEG {
		return ".jpg"
	}
	return "." + string(f)
}

// MIME is the format's media type.
func (f Format) MIME() string { return "image/" + string(f) }

// Variant is one output image: the source scaled to cover exactly
// Width x Height (center-cropped when the aspect ratio differs) and encoded
// as Format.
type Variant struct {
	Format Format
	Width  int
	Height int
}

func (v Variant) fileName() string {
	return fmt.Sprintf("%s-%dx%d%s", v.Format, v.Width, v.Height, v.Format.Ext())
}

// Job asks for a set of variants of one source file.
type Job struct {
	Source   string // path of the source raster
	Variants []Variant
}

// File is an encoded variant, stored in the cache directory.
type File struct {
	Variant
	// Path is the encoded file in the cache directory.
	Path string
	// Hash is the hex SHA-256 of the encoded bytes.
	Hash string
}

// Result is the outcome of one Job. Err is set when the source cannot be
// read or decoded or a variant cannot be encoded; Files is then nil.
type Result struct {
	Files []File // in the order of Job.Variants
	Err   error
}

// Probe reads just enough of a source raster to report its size and format.
func Probe(path string) (width, height int, format string, err error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, 0, "", err
	}
	defer f.Close()
	cfg, format, err := image.DecodeConfig(f)
	if err != nil {
		return 0, 0, "", fmt.Errorf("not a supported raster image (PNG, JPEG, WebP, AVIF): %w", err)
	}
	return cfg.Width, cfg.Height, format, nil
}

// Encoder runs jobs through a worker pool, reusing cached variants.
type Encoder struct {
	// CacheDir holds encoded variants between builds; it is created if
	// missing. Required.
	CacheDir string
	// Fast trades file size for encode speed (tests, local previews).
	Fast bool
	// Workers bounds parallel encodes; 0 means GOMAXPROCS.
	Workers int
}

type settings struct {
	avif avif.Options
	webp webp.Options
	jpeg jpeg.Options
}

func (e *Encoder) settings() settings {
	if e.Fast {
		return settings{
			avif: avif.Options{Quality: 50, Speed: 10, ChromaSubsampling: image.YCbCrSubsampleRatio420},
			webp: webp.Options{Quality: 75, Method: 0},
			jpeg: jpeg.Options{Quality: 75},
		}
	}
	return settings{
		avif: avif.Options{Quality: 60, Speed: 8, ChromaSubsampling: image.YCbCrSubsampleRatio420},
		webp: webp.Options{Quality: 80, Method: 4},
		jpeg: jpeg.Options{Quality: 82},
	}
}

// Run encodes every job's variants, in parallel across jobs and variants.
// A variant already in the cache (same source bytes, same settings) is not
// re-encoded. Content problems are reported per job in Result.Err; the
// returned error is for cache I/O failures only.
func (e *Encoder) Run(jobs []Job) ([]Result, error) {
	if e.CacheDir == "" {
		return nil, errors.New("images: CacheDir is required")
	}
	set := e.settings()
	results := make([]Result, len(jobs))
	sources := make([]*source, len(jobs))

	// Hash sources and find what is missing from the cache.
	type task struct {
		job int
		v   Variant
	}
	var tasks []task
	for i, j := range jobs {
		src, err := openSource(j.Source, set, e.CacheDir)
		if err != nil {
			results[i].Err = err
			continue
		}
		sources[i] = src
		for _, v := range j.Variants {
			if _, err := os.Stat(src.cached(v)); err == nil {
				continue
			}
			tasks = append(tasks, task{i, v})
		}
	}

	// Encode misses in parallel. Each source is decoded once, on first use.
	var (
		g      errgroup.Group
		mu     sync.Mutex
		jobErr = make([]error, len(jobs))
	)
	g.SetLimit(cmp.Or(e.Workers, runtime.GOMAXPROCS(0)))
	for _, t := range tasks {
		src := sources[t.job]
		g.Go(func() error {
			data, err := src.encode(t.v, set)
			if err != nil {
				mu.Lock()
				if jobErr[t.job] == nil {
					jobErr[t.job] = err
				}
				mu.Unlock()
				return nil
			}
			return writeAtomic(src.cached(t.v), data)
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}

	// Collect cached files and hash them for content-addressed names.
	for i, j := range jobs {
		if results[i].Err != nil {
			continue
		}
		if jobErr[i] != nil {
			results[i].Err = jobErr[i]
			continue
		}
		files := make([]File, 0, len(j.Variants))
		for _, v := range j.Variants {
			p := sources[i].cached(v)
			sum, err := hashFile(p)
			if err != nil {
				return nil, err
			}
			files = append(files, File{Variant: v, Path: p, Hash: sum})
		}
		results[i].Files = files
	}
	return results, nil
}

// source is one job's input: its cache directory and, once needed, its
// decoded pixels.
type source struct {
	path string
	dir  string // cache directory for this source + settings

	decode func() (image.Image, error) // decodes once, on first call
}

func openSource(path string, set settings, cacheRoot string) (*source, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	h := sha256.New()
	fmt.Fprintf(h, "v%s avif%+v webp%+v jpeg%+v\n", cacheVersion, set.avif, set.webp, set.jpeg)
	if _, err := io.Copy(h, f); err != nil {
		return nil, err
	}
	key := hex.EncodeToString(h.Sum(nil))
	s := &source{path: path, dir: filepath.Join(cacheRoot, key[:2], key)}
	s.decode = sync.OnceValues(func() (image.Image, error) {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		img, _, err := image.Decode(bytes.NewReader(data))
		if err != nil {
			return nil, fmt.Errorf("decode %s: %w", filepath.Base(path), err)
		}
		return img, nil
	})
	return s, nil
}

func (s *source) cached(v Variant) string { return filepath.Join(s.dir, v.fileName()) }

func (s *source) encode(v Variant, set settings) ([]byte, error) {
	img, err := s.decode()
	if err != nil {
		return nil, err
	}
	scaled := cover(img, v.Width, v.Height)
	var buf bytes.Buffer
	switch v.Format {
	case AVIF:
		err = avif.Encode(&buf, scaled, set.avif)
	case WebP:
		err = webp.Encode(&buf, scaled, set.webp)
	case JPEG:
		err = jpeg.Encode(&buf, scaled, &set.jpeg)
	default:
		err = fmt.Errorf("unknown format %q", v.Format)
	}
	if err != nil {
		return nil, fmt.Errorf("encode %s %dx%d: %w", strings.ToUpper(string(v.Format)), v.Width, v.Height, err)
	}
	return buf.Bytes(), nil
}

// cover scales src to fill exactly w x h, cropping the centre of whichever
// dimension overflows.
func cover(src image.Image, w, h int) image.Image {
	b := src.Bounds()
	sw, sh := b.Dx(), b.Dy()
	crop := b
	if sw*h > sh*w { // source is wider than the target: trim the sides
		cw := (sh*w + h/2) / h
		x := b.Min.X + (sw-cw)/2
		crop = image.Rect(x, b.Min.Y, x+cw, b.Max.Y)
	} else if sw*h < sh*w { // source is taller: trim top and bottom
		ch := (sw*h + w/2) / w
		y := b.Min.Y + (sh-ch)/2
		crop = image.Rect(b.Min.X, y, b.Max.X, y+ch)
	}
	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.CatmullRom.Scale(dst, dst.Bounds(), src, crop, draw.Src, nil)
	return dst
}

func hashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// writeAtomic writes data so a concurrent or interrupted build never sees a
// partial cache file.
func writeAtomic(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), ".tmp-*")
	if err != nil {
		return err
	}
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		os.Remove(tmp.Name())
		return err
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmp.Name())
		return err
	}
	return os.Rename(tmp.Name(), path)
}
