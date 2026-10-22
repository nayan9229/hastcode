package builder_test

import (
	"bytes"
	"image"
	"net/http"
	"os"
	"path"
	"slices"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"hastcode.com/site/internal/buildtest"
)

// r2Tree is ValidTree plus an Author photo: every kind of encoded image
// (hero variants and crops, Author photo, OG badge).
func r2Tree() buildtest.Tree {
	return buildtest.ValidTree().WithAuthor(authorWithPhoto())
}

// referencedImageURLs are the image URLs the built pages point at: the
// article's <picture> sources and <img>, og:image, the Article JSON-LD
// images, the Author photo and Person image, and the home page's OG badge.
func referencedImageURLs(t *testing.T, out string) []string {
	t.Helper()
	var urls []string
	add := func(u string) {
		if u != "" && !slices.Contains(urls, u) {
			urls = append(urls, u)
		}
	}
	article := buildtest.Page(t, out, raftPage)
	article.Find("article picture source, article picture img").Each(func(_ int, s *goquery.Selection) {
		for _, e := range buildtest.Srcset(t, s.AttrOr("srcset", "")) {
			add(e.URL)
		}
	})
	add(article.Find("article picture img").AttrOr("src", ""))
	add(article.Find(`meta[property="og:image"]`).AttrOr("content", ""))
	for _, im := range buildtest.JSONLDOfType(t, article, "Article")["image"].([]any) {
		add(im.(string))
	}
	author := buildtest.Page(t, out, adaURL)
	add(author.Find("main .h-card img.u-photo").AttrOr("src", ""))
	add(buildtest.JSONLDOfType(t, author, "Person")["image"].(string))
	add(buildtest.Page(t, out, "/").Find(`meta[property="og:image"]`).AttrOr("content", ""))
	slices.Sort(urls)
	return urls
}

func TestImageStore_R2UploadsEveryImageOnceWithImmutableHeadersAndWritesNoneToOutput(t *testing.T) {
	store := buildtest.NewImageStore(t)
	out := buildtest.MustBuild(t, store.Configure(buildtest.Config()), r2Tree())

	var wantKeys []string
	for _, u := range referencedImageURLs(t, out) {
		key, ok := strings.CutPrefix(u, store.BaseURL)
		if !ok {
			t.Fatalf("image URL %s is not on the Image Store %s", u, store.BaseURL)
		}
		wantKeys = append(wantKeys, key)
	}
	// 1600x900 hero: 5 widths x (AVIF, JPEG) + 4 crops, less the 16:9 crop
	// that is the 1200 px JPEG; 1 Author photo; 1 OG badge.
	if len(wantKeys) != 15 {
		t.Fatalf("pages reference %d images, want 15: %v", len(wantKeys), wantKeys)
	}
	if got := store.RequestKeys(http.MethodPut); !slices.Equal(got, wantKeys) {
		t.Errorf("PUT keys:\n got %v\nwant %v", got, wantKeys)
	}

	wantType := map[string]string{".avif": "image/avif", ".jpg": "image/jpeg", ".png": "image/png"}
	for _, r := range store.Requests(http.MethodPut) {
		if !strings.HasPrefix(r.Key, "heroes/raft-leader-election-") && !strings.HasPrefix(r.Key, "authors/ada-") && !strings.HasPrefix(r.Key, "og-badge.") {
			t.Errorf("key %s does not mirror the asset path without /assets/", r.Key)
		}
		if got, want := r.Header.Get("Content-Type"), wantType[path.Ext(r.Key)]; got != want {
			t.Errorf("%s Content-Type = %q, want %q", r.Key, got, want)
		}
		if got := r.Header.Get("Cache-Control"); got != "public, max-age=31536000, immutable" {
			t.Errorf("%s Cache-Control = %q", r.Key, got)
		}
		for h := range r.Header {
			if strings.HasPrefix(strings.ToLower(h), "x-amz-checksum-") || strings.EqualFold(h, "x-amz-sdk-checksum-algorithm") {
				t.Errorf("%s carries checksum header %s (R2 wants checksums only when required)", r.Key, h)
			}
		}
	}

	fallback := buildtest.Page(t, out, raftPage).Find("article picture img").AttrOr("src", "")
	obj, ok := store.Object(strings.TrimPrefix(fallback, store.BaseURL))
	if !ok {
		t.Fatalf("fallback <img> %s not in the store", fallback)
	}
	img, format, err := image.Decode(bytes.NewReader(obj.Data))
	if err != nil || format != "jpeg" || img.Bounds().Dx() != 1200 || img.Bounds().Dy() != 675 {
		t.Errorf("stored fallback image: format %q, %v, err %v; want a 1200x675 JPEG", format, img.Bounds(), err)
	}

	if imgs := buildtest.Images(t, out); len(imgs) != 0 {
		t.Errorf("r2 build wrote %d image files to the output, want none: %v", len(imgs), imgs)
	}
}

func TestImageStore_R2PagesReferenceEveryImageByAbsoluteStoreURL(t *testing.T) {
	store := buildtest.NewImageStore(t)
	out := buildtest.MustBuild(t, store.Configure(buildtest.Config()), r2Tree().WithBrief(buildtest.ValidBrief()))

	onStore := func(page, what, u string) {
		t.Helper()
		if !strings.HasPrefix(u, store.BaseURL) {
			t.Errorf("%s %s = %q, want an absolute URL on %s", page, what, u, store.BaseURL)
		}
	}
	pages := 0
	for _, page := range buildtest.URLPaths(t, out) {
		if !strings.HasSuffix(page, "/") && !strings.HasSuffix(page, ".html") {
			continue
		}
		pages++
		doc := buildtest.Page(t, out, page)
		doc.Find("img[src]").Each(func(_ int, s *goquery.Selection) { onStore(page, "<img src>", s.AttrOr("src", "")) })
		// srcset candidates are "url [descriptor]" (listing thumbnails have one bare URL).
		candidates := func(what, srcset string) {
			for _, c := range strings.Split(srcset, ",") {
				onStore(page, what, strings.Fields(c)[0])
			}
		}
		doc.Find("img[srcset], source[srcset]").Each(func(_ int, s *goquery.Selection) { candidates("srcset", s.AttrOr("srcset", "")) })
		doc.Find(`link[rel="preload"][as="image"]`).Each(func(_ int, s *goquery.Selection) {
			candidates("preload imagesrcset", s.AttrOr("imagesrcset", ""))
		})
		doc.Find(`meta[property="og:image"], meta[name="twitter:image"]`).Each(func(_ int, s *goquery.Selection) {
			onStore(page, "og/twitter image", s.AttrOr("content", ""))
		})
		for _, block := range buildtest.JSONLD(t, doc) {
			switch im := block["image"].(type) {
			case string:
				onStore(page, "JSON-LD image", im)
			case []any:
				for _, u := range im {
					onStore(page, "JSON-LD image", u.(string))
				}
			}
		}
	}
	if pages < 10 {
		t.Fatalf("checked only %d pages", pages)
	}
	brief := buildtest.Page(t, out, "/briefs/etcd-raft-rewrite/")
	if news := buildtest.JSONLDOfType(t, brief, "NewsArticle"); news == nil || len(news["image"].([]any)) != 3 {
		t.Errorf("Brief NewsArticle images = %v, want 3 crops", news["image"])
	}

	// Nothing anywhere (pages, feeds, sitemaps, manifests) points at a
	// local image path or prefixes site_url onto a store URL.
	for _, s := range []string{"/assets/heroes/", "/assets/authors/", "/assets/og-badge", "hastcode.testhttps", "hastcode.test/https"} {
		if hits := buildtest.FilesContaining(t, out, s); len(hits) != 0 {
			t.Errorf("%q appears in %v", s, hits)
		}
	}
}

func TestImageStore_SecondR2BuildHEADsEveryKeyAndPUTsNothing(t *testing.T) {
	store := buildtest.NewImageStore(t)
	cfg := store.Configure(buildtest.Config())
	buildtest.MustBuild(t, cfg, r2Tree())
	uploaded := store.RequestKeys(http.MethodPut)
	store.ResetRequests()

	buildtest.MustBuild(t, cfg, r2Tree())

	if puts := store.RequestKeys(http.MethodPut); len(puts) != 0 {
		t.Errorf("second build PUT %v, want nothing", puts)
	}
	if heads := store.RequestKeys(http.MethodHead); !slices.Equal(heads, uploaded) {
		t.Errorf("second build HEADs:\n got %v\nwant every key once %v", heads, uploaded)
	}
}

func TestImageStore_R2BuildPUTsOnlyKeysTheStoreLacks(t *testing.T) {
	first := buildtest.NewImageStore(t)
	buildtest.MustBuild(t, first.Configure(buildtest.Config()), r2Tree())
	keys := first.RequestKeys(http.MethodPut)

	store := buildtest.NewImageStore(t)
	seeded := keys[0]
	store.Seed(seeded, buildtest.StoredObject{Data: []byte("already there")})
	buildtest.MustBuild(t, store.Configure(buildtest.Config()), r2Tree())

	if got := store.RequestKeys(http.MethodPut); !slices.Equal(got, keys[1:]) {
		t.Errorf("PUT keys:\n got %v\nwant all but the seeded %s", got, seeded)
	}
	if o, _ := store.Object(seeded); string(o.Data) != "already there" {
		t.Errorf("seeded object %s was overwritten", seeded)
	}
}

func TestImageStore_R2WithoutACredentialFailsBeforeEncodingNamingIt(t *testing.T) {
	for _, missing := range []string{"CLOUDFLARE_ACCOUNT_ID", "R2_ACCESS_KEY_ID", "R2_SECRET_ACCESS_KEY"} {
		t.Run(missing, func(t *testing.T) {
			store := buildtest.NewImageStore(t)
			env := store.Env()
			delete(env, missing)
			cfg := store.ConfigureEnv(buildtest.Config(), env)
			cfg.Images.CacheDir = t.TempDir()

			out, res := buildtest.Build(t, cfg, r2Tree())

			buildtest.RequireErrorRule(t, res, "", "image-store", missing)
			if entries, _ := os.ReadDir(cfg.Images.CacheDir); len(entries) != 0 {
				t.Errorf("image cache has %d entries: the build encoded before failing", len(entries))
			}
			if reqs := store.Requests(""); len(reqs) != 0 {
				t.Errorf("store received %d requests", len(reqs))
			}
			if buildtest.Exists(out, "/") {
				t.Error("the build wrote pages")
			}
		})
	}
}

func TestImageStore_R2NeedsBaseURLAndBucket(t *testing.T) {
	store := buildtest.NewImageStore(t)
	noBase := store.Configure(buildtest.Config())
	noBase.Images.BaseURL = ""
	_, res := buildtest.Build(t, noBase, r2Tree())
	buildtest.RequireErrorRule(t, res, "", "image-store", "images.base_url")

	noBucket := store.Configure(buildtest.Config())
	noBucket.Images.Bucket = ""
	_, res = buildtest.Build(t, noBucket, r2Tree())
	buildtest.RequireErrorRule(t, res, "", "image-store", "images.bucket")

	unknown := buildtest.Config()
	unknown.Images.Store = "s3"
	_, res = buildtest.Build(t, unknown, r2Tree())
	buildtest.RequireErrorRule(t, res, "", "image-store", "images.store")
}

func TestImageStore_LocalModeNeedsNoCredentials(t *testing.T) {
	cfg := buildtest.Config()
	cfg.Getenv = func(string) string { return "" }
	out := buildtest.MustBuild(t, cfg, r2Tree())
	src := buildtest.Page(t, out, raftPage).Find("article picture img").AttrOr("src", "")
	if !strings.HasPrefix(src, "/assets/heroes/") || !buildtest.Exists(out, src) {
		t.Errorf("local hero <img src> = %q, want a written file under /assets/heroes/", src)
	}
}

func TestImageStore_FailedUploadFailsTheBuildNamingTheKey(t *testing.T) {
	first := buildtest.NewImageStore(t)
	buildtest.MustBuild(t, first.Configure(buildtest.Config()), r2Tree())
	key := first.RequestKeys(http.MethodPut)[0]

	for _, method := range []string{http.MethodPut, http.MethodHead} {
		t.Run(method, func(t *testing.T) {
			store := buildtest.NewImageStore(t)
			store.Fail(method, key, http.StatusForbidden)
			_, res := buildtest.Build(t, store.Configure(buildtest.Config()), r2Tree())
			buildtest.RequireError(t, res, "", key)
			if _, ok := store.Object(key); ok {
				t.Errorf("%s stored despite the failure", key)
			}
		})
	}
}

// The prune command (ticket 04) keeps the keys a build reports; a local
// build reports the same keys an r2 build uploads, without uploading.
func TestImageStore_ResultListsTheKeysInEitherMode(t *testing.T) {
	store := buildtest.NewImageStore(t)
	_, r2 := buildtest.Build(t, store.Configure(buildtest.Config()), r2Tree())
	_, local := buildtest.Build(t, buildtest.Config(), r2Tree())

	uploaded := store.RequestKeys(http.MethodPut)
	if len(uploaded) == 0 || !slices.Equal(r2.ImageKeys, uploaded) {
		t.Errorf("r2 ImageKeys:\n got %v\nwant the uploaded %v", r2.ImageKeys, uploaded)
	}
	if !slices.Equal(local.ImageKeys, uploaded) {
		t.Errorf("local ImageKeys:\n got %v\nwant %v", local.ImageKeys, uploaded)
	}
}
