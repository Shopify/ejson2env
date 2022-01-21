NAME=ejson2env
RUBY_MODULE=EJSON2ENV
PACKAGE=github.com/Shopify/ejson2env
VERSION=$(shell cat VERSION)
GEM=pkg/$(NAME)-$(VERSION).gem
DEB=pkg/$(NAME)_$(VERSION)_amd64.deb

GOFILES=$(shell find . -type f -name '*.go')
MANFILES=$(shell find man -name '*.ronn' -exec echo build/{} \; | sed 's/\.ronn/\.gz/')

BUNDLE_EXEC=bundle exec

.PHONY: default all binaries gem man clean dev_bootstrap

default: all
all: gem deb
binaries: build/bin/linux-amd64 build/bin/linux-arm64 build/bin/darwin-universal build/bin/freebsd-amd64
gem: $(GEM)
deb: $(DEB)
man: $(MANFILES)

build/man/%.gz: man/%.ronn
	mkdir -p "$(@D)"
	set -euo pipefail ; $(BUNDLE_EXEC) ronn -r --pipe "$<" | gzip > "$@" || (rm -f "$<" && false)

build/bin/linux-amd64: $(GOFILES)
	mkdir -p "$(@D)"
	GOOS=linux GOARCH=amd64 go build \
	-ldflags '-s -w -X main.version="$(VERSION)"' \
	-o "$@" \
	"$(PACKAGE)/cmd/$(NAME)"

build/bin/linux-arm64: $(GOFILES)
	mkdir -p "$(@D)"
	GOOS=linux GOARCH=arm64 go build \
	-ldflags '-s -w -X main.version="$(VERSION)"' \
	-o "$@" \
	"$(PACKAGE)/cmd/$(NAME)"

build/bin/darwin-amd64: $(GOFILES)
	GOOS=darwin GOARCH=amd64 go build \
	-ldflags '-s -w -X main.version="$(VERSION)"' \
	-o "$@" \
	"$(PACKAGE)/cmd/$(NAME)"

build/bin/darwin-arm64: $(GOFILES)
	GOOS=darwin GOARCH=arm64 go build \
	-ldflags '-s -w -X main.version="$(VERSION)"' \
	-o "$@" \
	"$(PACKAGE)/cmd/$(NAME)"

build/bin/darwin-universal: build/bin/darwin-amd64 build/bin/darwin-arm64
	$(V)lipo -create -output "$@" $^

build/bin/freebsd-amd64: $(GOFILES)
	GOOS=freebsd GOARCH=amd64 go build \
	-ldflags '-s -w -X main.version="$(VERSION)"' \
	-o "$@" \
	"$(PACKAGE)/cmd/$(NAME)"

$(GEM): rubygem/$(NAME)-$(VERSION).gem
	mkdir -p $(@D)
	mv "$<" "$@"
	
rubygem/$(NAME)-$(VERSION).gem: \
	rubygem/lib/$(NAME)/version.rb \
	rubygem/build/linux-amd64/ejson2env \
	rubygem/LICENSE.txt \
	rubygem/build/darwin-universal/ejson2env \
	rubygem/build/freebsd-amd64/ejson2env \
	rubygem/man
	cd rubygem && gem build ejson2env.gemspec

rubygem/LICENSE.txt: LICENSE.txt
	cp "$<" "$@"

rubygem/man: man
	cp -a build/man $@

rubygem/build/darwin-universal/ejson2env: build/bin/darwin-universal
	mkdir -p $(@D)
	cp -a "$<" "$@"

rubygem/build/linux-amd64/ejson2env: build/bin/linux-amd64
	mkdir -p $(@D)
	cp -a "$<" "$@"

rubygem/build/freebsd-amd64/ejson2env: build/bin/freebsd-amd64
	mkdir -p $(@D)
	cp -a "$<" "$@"

rubygem/lib/$(NAME)/version.rb: VERSION
	mkdir -p $(@D)
	echo 'module $(RUBY_MODULE)\n  VERSION = "$(VERSION)"\nend' > $@

$(DEB): build/bin/linux-amd64 man
	mkdir -p $(@D)
	rm -f "$@"
	$(BUNDLE_EXEC) fpm \
		-t deb \
		-s dir \
		--name="$(NAME)" \
		--version="$(VERSION)" \
		--package="$@" \
		--license=MIT \
		--category=admin \
		--no-depends \
		--no-auto-depends \
		--architecture=amd64 \
		--maintainer="Catherine Jones <catherine.jones@shopify.com>" \
		--description="utility for decrypting ejson secrets and helping to export them as environment variables" \
		--url="https://github.com/Shopify/ejson2env" \
		./build/man/=/usr/share/man/ \
		./build/bin/linux-amd64=/usr/bin/ejson2env

clean:
	rm -rf build pkg rubygem/{LICENSE.txt,lib/ejson2env/version.rb,build,*.gem}
