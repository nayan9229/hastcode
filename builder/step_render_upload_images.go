package builder

import (
	"context"
	"slices"
	"strings"

	"hastcode.com/site/internal/imagestore"
)

func init() {
	registerStep(phaseRender, 17, "upload images", uploadImages)
}

// uploadWorkers bounds concurrent HEAD/PUT requests (research §1.5: 8–16).
const uploadWorkers = 8

// uploadImages sends the images queued by storeImage (hero assets, Author
// photos, the OG badge) to the Image Store in r2 mode: HEAD each key, PUT
// only the missing ones. A failure stops the build naming the key.
func uploadImages(s *site) error {
	if s.store == nil || len(s.uploads) == 0 {
		return nil
	}
	// Two pieces with the same hero file name and bytes share keys; send
	// each key once (R2 answers 429 to concurrent writes of one key).
	objs := slices.CompactFunc(slices.SortedStableFunc(slices.Values(s.uploads), func(a, b imagestore.Object) int {
		return strings.Compare(a.Key, b.Key)
	}), func(a, b imagestore.Object) bool { return a.Key == b.Key })
	return s.store.Upload(context.Background(), objs, uploadWorkers)
}
