package builder

import (
	"fmt"
	"path"
	"strings"

	"hastcode.com/site/internal/images"
	"hastcode.com/site/internal/views"
)

func init() {
	registerStep(phaseTransform, 55, "encode author photos", encodeAuthorPhotos)
}

// encodeAuthorPhotos encodes every valid photo as one square JPEG.
func encodeAuthorPhotos(s *site) error {
	var (
		jobs   []images.Job
		owners []string
	)
	for _, a := range s.sortedAuthors() {
		ph := s.authorPhotos[a.Slug]
		if ph == nil || s.failed[a.File] {
			continue
		}
		jobs = append(jobs, images.Job{Source: ph.source, Variants: []images.Variant{
			{Format: images.JPEG, Width: authorPhotoSize, Height: authorPhotoSize},
		}})
		owners = append(owners, a.Slug)
	}
	if len(jobs) == 0 {
		return nil
	}
	enc := &images.Encoder{CacheDir: s.cfg.Images.CacheDir, Fast: s.cfg.Images.Fast}
	results, err := enc.Run(jobs)
	if err != nil {
		return err
	}
	for i, res := range results {
		a := s.authors[owners[i]]
		if res.Err != nil {
			s.fail(a.File, ruleAuthorPhoto, "photo: %v", res.Err)
			continue
		}
		f := res.Files[0]
		base := strings.TrimSuffix(path.Base(a.Photo), path.Ext(a.Photo))
		ph := s.authorPhotos[a.Slug]
		ph.cached = f.Path
		ph.view = views.Image{
			URL:    fmt.Sprintf("%s%s-%dx%d.%s%s", authorPhotoDir, base, f.Width, f.Height, f.Hash[:10], f.Format.Ext()),
			Width:  f.Width,
			Height: f.Height,
		}
	}
	return nil
}
