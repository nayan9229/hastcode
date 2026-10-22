# templ output (*_templ.go) is gitignored: every target that compiles Go
# runs `generate` first.
#
# The AVIF encoder (gen2brain/avif) runs as WASM in-process: CGO_ENABLED=0
# plus -tags nodynamic stops it loading a system libavif, so output is
# byte-identical on laptops, CI and Vercel.
export CGO_ENABLED := 0
export GOFLAGS := -tags=nodynamic

CONTENT ?= content
CONFIG  ?= site.yaml
OUT     ?= .vercel/output
MODE    ?= production
# Extra builder flags, e.g. FLAGS=-skip-search-index with MODE=preview.
FLAGS   ?=

.PHONY: generate build test vet check dev adapt

# `make dev`: build (preview mode by default, so drafts show) and serve the
# output plus the subscribe function behind the built config.json routes
# (410s, /home 301, trailing-slash 308s, 404) on http://localhost:$(PORT).
# BUTTONDOWN_API_KEY / BUTTONDOWN_BASE_URL come from the shell, else from
# .env (scripts/setup-deploy.sh writes the test newsletter's key there),
# else stay unset (every subscribe POST then answers 502).
# MODE is also read from the environment: the CI preview job sets
# MODE=preview for `vercel build`, which runs `make build` (vercel.json).
DEV_MODE ?= preview
PORT     ?= 3000

generate:
	go generate ./...

build: generate
	go run ./cmd/builder -content $(CONTENT) -config $(CONFIG) -out $(OUT) -mode $(MODE) $(FLAGS)

test: generate
	go test ./...

vet: generate
	go vet ./...

check: vet test

# `make adapt TARGET=cloudflare|vercel`: the deploy-adapter rule
# (internal/deploy). Removes the other Target's manifest files from $(OUT)
# in place, once, after the build and before the upload. Cloudflare then
# deploys with `wrangler deploy --config $(OUT)/wrangler.json`; Vercel with
# `vercel deploy --prebuilt`.
TARGET ?=
adapt:
	go run ./cmd/adapt -target "$(TARGET)" -out $(OUT)

dev:
	$(MAKE) build MODE=$(DEV_MODE)
	key="$${BUTTONDOWN_API_KEY-}"; base="$${BUTTONDOWN_BASE_URL-}"; \
	  set -a; [ ! -f .env ] || . ./.env; set +a; \
	  BUTTONDOWN_API_KEY="$${key:-$${BUTTONDOWN_API_KEY-}}" \
	  BUTTONDOWN_BASE_URL="$${base:-$${BUTTONDOWN_BASE_URL-}}" \
	  go run ./cmd/subscribe -addr localhost:$(PORT) -out $(OUT)
