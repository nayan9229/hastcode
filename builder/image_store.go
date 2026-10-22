package builder

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"hastcode.com/site/internal/imagestore"
)

// The Image Store (ADR-0004, spec "Image Store"): where encoded images
// (hero variants and crops, Author photos, the OG badge) are served from.
// Steps name an image by its site path under assetsDir
// ("/assets/heroes/x-1200x675.<hash>.jpg"), get the URL pages use from
// s.imageURL, and hand the bytes to s.storeImage. In local mode the image
// is written to that path in the output; in r2 mode its key is the path
// without /assets/ ("heroes/x-…jpg"), the URL is absolute on
// images.base_url, and uploadImages (step_render_upload_images.go) HEADs
// and PUTs it.

const ruleImageStore = "image-store"

// uploadWorkers bounds concurrent HEAD/PUT requests (research §1.5: 8–16).
const uploadWorkers = 8

// openImageStore checks the image store config and, in r2 mode, reads the
// credentials and returns the client. It runs before any step, so a
// misconfigured r2 build fails before encoding anything.
func openImageStore(cfg Config) (*imagestore.Client, error) {
	ic := cfg.Images
	switch ic.Store {
	case StoreLocal:
		return nil, nil
	case StoreR2:
	default:
		return nil, fmt.Errorf("images.store %q must be %q or %q", ic.Store, StoreLocal, StoreR2)
	}
	if u, err := url.Parse(ic.BaseURL); err != nil || (u.Scheme != "https" && u.Scheme != "http") || u.Host == "" {
		return nil, fmt.Errorf("images.store is r2 but images.base_url %q is not an absolute http(s) URL", ic.BaseURL)
	}
	if ic.Bucket == "" {
		return nil, fmt.Errorf("images.store is r2 but images.bucket is blank")
	}
	creds, err := imagestore.FromEnv(cfg.Getenv)
	if err != nil {
		return nil, fmt.Errorf("images.store is r2 but %v (CI secrets; make dev uses images.store local)", err)
	}
	return imagestore.New(creds, ic.Bucket), nil
}

// imageKey is the Image Store key of an image's site path.
func imageKey(sitePath string) string { return strings.TrimPrefix(sitePath, assetsDir) }

// imageURL is the URL pages use for the image at sitePath: the site path
// itself in local mode, absolute on images.base_url in r2 mode.
func (s *site) imageURL(sitePath string) string {
	if s.store == nil {
		return sitePath
	}
	return strings.TrimSuffix(s.cfg.Images.BaseURL, "/") + "/" + imageKey(sitePath)
}

// storeImage publishes the image at sitePath: data, or the file at path
// when data is nil. Local mode writes it to the output now; r2 mode queues
// it for uploadImages.
func (s *site) storeImage(sitePath, path string, data []byte, contentType string) error {
	s.result.ImageKeys = append(s.result.ImageKeys, imageKey(sitePath))
	if s.store == nil {
		if data == nil {
			var err error
			if data, err = os.ReadFile(path); err != nil {
				return err
			}
		}
		return s.out.WriteStatic(sitePath, data)
	}
	s.uploads = append(s.uploads, imagestore.Object{Key: imageKey(sitePath), ContentType: contentType, Data: data, File: path})
	return nil
}
