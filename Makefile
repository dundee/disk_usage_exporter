NAME := disk_usage_exporter
PACKAGE := github.com/dundee/$(NAME)
VERSION := $(shell git describe --tags 2>/dev/null)
GIT_SHA := $(shell git rev-parse HEAD 2>/dev/null)
GOFLAGS ?= -buildmode=pie -trimpath -mod=readonly -modcacherw
LDFLAGS := -s -w -extldflags '-static' \
	-X '$(PACKAGE)/build.BuildVersion=$(VERSION)' \
	-X '$(PACKAGE)/build.BuildCommitSha=$(GIT_SHA)' \
	-X '$(PACKAGE)/build.BuildDate=$(shell LC_ALL=en_US.UTF-8 date)'

all: build-all clean-uncompressed-dist shasums

run:
	go run .

build:
	@echo "Version: " $(VERSION)
	mkdir -p dist
	GOFLAGS="$(GOFLAGS)" CGO_ENABLED=0 go build -ldflags="$(LDFLAGS)" -o dist/$(NAME) .

build-all:
	@echo "Version: " $(VERSION)
	-mkdir dist
	-CGO_ENABLED=0 gox \
		-os="darwin" \
		-arch="amd64" \
		-output="dist/disk_usage_exporter_{{.OS}}_{{.Arch}}" \
		-ldflags="$(LDFLAGS)"

	-CGO_ENABLED=0 gox \
		-os="linux freebsd netbsd openbsd" \
		-output="dist/disk_usage_exporter_{{.OS}}_{{.Arch}}" \
		-ldflags="$(LDFLAGS)"

	cd dist; GOFLAGS="$(GOFLAGS)" CGO_ENABLED=0 go build -ldflags="$(LDFLAGS)" -o disk_usage_exporter_linux_amd64 ..

	cd dist; CGO_ENABLED=0 GOOS=linux GOARM=5 GOARCH=arm go build -ldflags="$(LDFLAGS)" -o disk_usage_exporter_linux_armv5l ..
	cd dist; CGO_ENABLED=0 GOOS=linux GOARM=6 GOARCH=arm go build -ldflags="$(LDFLAGS)" -o disk_usage_exporter_linux_armv6l ..
	cd dist; CGO_ENABLED=0 GOOS=linux GOARM=7 GOARCH=arm go build -ldflags="$(LDFLAGS)" -o disk_usage_exporter_linux_armv7l ..
	cd dist; CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o disk_usage_exporter_linux_arm64 ..

	cd dist; for file in disk_usage_exporter_linux_* disk_usage_exporter_darwin_* disk_usage_exporter_netbsd_* disk_usage_exporter_openbsd_* disk_usage_exporter_freebsd_*; do tar czf $$file.tgz $$file; done

clean:
	-rm -r dist

clean-uncompressed-dist:
	find dist -type f -not -name '*.tgz' -not -name '*.zip' -delete

shasums:
	cd dist; sha256sum * > sha256sums.txt
	cd dist; gpg --sign --armor --detach-sign sha256sums.txt

.PHONY: run build clean
