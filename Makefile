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

.PHONY: generate build test vet check dev adapt prune-images lint-workflows fetch refresh intake lead lead-dismiss

# `make dev`: build (preview mode by default, so drafts show) and serve the
# output plus the subscribe function behind the built config.json routes
# (410s, /home 301, trailing-slash 308s, 404) on http://localhost:$(PORT).
# BUTTONDOWN_API_KEY / BUTTONDOWN_BASE_URL come from the shell, else from
# .env (scripts/setup-deploy.sh writes the test newsletter's key there),
# else stay unset (every subscribe POST then answers 502).
# MODE and FLAGS are also read from the environment: the QA Publish sets
# MODE=production FLAGS=-image-store=r2 for `vercel build`, which runs `make build` (vercel.json).
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

# `make lint-workflows`: actionlint over .github/workflows/ (tests.yml runs
# it on every push to main). Pinned; `go run` fetches it into the module
# cache, so there is nothing to install. It also runs shellcheck over the
# `run:` scripts when shellcheck is on PATH. The one ignored message is
# GitHub's `concurrency.queue` (publish.yml), which 1.7.12 predates.
ACTIONLINT := github.com/rhysd/actionlint/cmd/actionlint@v1.7.12
lint-workflows:
	go run $(ACTIONLINT) -ignore 'unexpected key "queue" for "concurrency" section'

# `make adapt TARGET=cloudflare|vercel`: the deploy-adapter rule
# (internal/deploy). Removes the other Target's manifest files from $(OUT)
# in place, once, after the build and before the upload. Cloudflare then
# deploys with `wrangler deploy --config $(OUT)/wrangler.json`; Vercel with
# `vercel deploy --prebuilt`.
TARGET ?=
adapt:
	go run ./cmd/adapt -target "$(TARGET)" -out $(OUT)

# `make prune-images`: list the R2 objects the current $(CONTENT) would not
# produce (sizes and a total). Dry run: deletes nothing. Needs the r2
# build's environment (CLOUDFLARE_ACCOUNT_ID, R2_ACCESS_KEY_ID,
# R2_SECRET_ACCESS_KEY). To delete the listed orphans:
# `make prune-images PRUNE_FLAGS=--delete`. See README "Image Store".
PRUNE_FLAGS ?=
prune-images: generate
	go run ./cmd/prune -content $(CONTENT) -config $(CONFIG) $(PRUNE_FLAGS)

# `make fetch REPO=owner/name MD=content/explainers/<vertical>/<slug>.md`:
# write the Repo Profile's Fact Pack beside its Markdown and the README
# text under the build cache (node_modules/.cache/hastcode/readmes/).
# `make refresh`: re-fetch every Fact Pack under $(CONTENT), rewrite the
# changed ones and print one line per changed repository; the command
# exits 3 when anything changed, 0 when nothing did (make then prints
# "Error 3" and itself exits 2: the "changed" signal, not a failure). Both
# read GITHUB_TOKEN and refuse without it. The binary is built first
# because `go run` would turn exit 3 into 1. A Refresh is manual: review
# the diff, commit, push (README "Refresh"). See README "Repo
# Profiles: fetching a Fact Pack".
FETCH_BIN := node_modules/.cache/hastcode/bin/fetch
REPO ?=
MD   ?=
fetch:
	mkdir -p $(dir $(FETCH_BIN)) && go build -o $(FETCH_BIN) ./cmd/fetch
	$(FETCH_BIN) -repo "$(REPO)" -markdown "$(MD)"
refresh:
	mkdir -p $(dir $(FETCH_BIN)) && go build -o $(FETCH_BIN) ./cmd/fetch
	$(FETCH_BIN) -refresh -content $(CONTENT)

# `make intake`: check every Source on the Source Watchlist (sources.yaml)
# and list the new Leads in $(INTAKE_STATE)/queue.md (ADR-0007). Local
# only, like a Refresh; never in GitHub Actions. The queue and the
# seen-state live in $(INTAKE_STATE), gitignored; nothing is written under
# $(CONTENT). `make lead-dismiss ID=…` hides a Lead for good; `make lead
# ID=… [KIND=…] [VERTICAL=…]` starts a draft skeleton from one. Repository
# Sources need GITHUB_TOKEN (the REST API, as for `make fetch`). See README
# "Intake".
INTAKE_STATE ?= .intake
WATCHLIST    ?= sources.yaml
INTAKE_FLAGS  = -watchlist $(WATCHLIST) -config $(CONFIG) -content $(CONTENT) -state $(INTAKE_STATE)
ID       ?=
KIND     ?=
VERTICAL ?=
intake: generate
	go run ./cmd/intake $(INTAKE_FLAGS)
lead-dismiss: generate
	go run ./cmd/intake $(INTAKE_FLAGS) -dismiss "$(ID)"
lead: generate
	go run ./cmd/intake $(INTAKE_FLAGS) -pick "$(ID)" -kind "$(KIND)" -vertical "$(VERTICAL)"

dev:
	$(MAKE) build MODE=$(DEV_MODE)
	key="$${BUTTONDOWN_API_KEY-}"; base="$${BUTTONDOWN_BASE_URL-}"; \
	  set -a; [ ! -f .env ] || . ./.env; set +a; \
	  BUTTONDOWN_API_KEY="$${key:-$${BUTTONDOWN_API_KEY-}}" \
	  BUTTONDOWN_BASE_URL="$${base:-$${BUTTONDOWN_BASE_URL-}}" \
	  go run ./cmd/subscribe -addr localhost:$(PORT) -out $(OUT)
