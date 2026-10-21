# templ output (*_templ.go) is gitignored: every target that compiles Go
# runs `generate` first.

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
