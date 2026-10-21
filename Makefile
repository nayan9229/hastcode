# templ output (*_templ.go) is gitignored: every target that compiles Go
# runs `generate` first.
#
# Image encoders (gen2brain/avif, gen2brain/webp) run as WASM in-process:
# CGO_ENABLED=0 plus -tags nodynamic stops them loading a system libavif or
# libwebp, so output is byte-identical on laptops, CI and Vercel.
export CGO_ENABLED := 0
export GOFLAGS := -tags=nodynamic

CONTENT ?= content
CONFIG  ?= site.yaml
OUT     ?= .vercel/output
MODE    ?= production

.PHONY: generate build test vet check

generate:
	go generate ./...

build: generate
	go run ./cmd/builder -content $(CONTENT) -config $(CONFIG) -out $(OUT) -mode $(MODE)

test: generate
	go test ./...

vet: generate
	go vet ./...

check: vet test
