#!/usr/bin/env bash
set -e

# Build gen for multiple platforms. Outputs to dist/.
# Usage: ./scripts/build.sh [version]
# Version defaults to "dev" if not set.

VERSION="${1:-dev}"
DIST="dist"
BINARY_NAME="gen"
REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$REPO_ROOT"

mkdir -p "$DIST"

build() {
	local goos="$1"
	local goarch="$2"
	local archive_ext="${3:-tar.gz}"
	local exe="$BINARY_NAME"
	[[ "$goos" == "windows" ]] && exe="${BINARY_NAME}.exe"

	export CGO_ENABLED=0
	export GOOS="$goos"
	export GOARCH="$goarch"
	echo "Building $goos/$goarch -> $DIST/gen_${goos}_${goarch}.${archive_ext}"
	go build -ldflags "-X main.Version=$VERSION" -o "$DIST/$exe" ./cmd/gen

	if [[ "$archive_ext" == "zip" ]]; then
		zip -j "$DIST/gen_${goos}_${goarch}.zip" "$DIST/$exe"
	else
		tar -czf "$DIST/gen_${goos}_${goarch}.tar.gz" -C "$DIST" "$exe"
	fi
	rm -f "$DIST/$exe"
}

build darwin arm64
build darwin amd64
build linux amd64
build linux arm64
build windows amd64 zip

echo "Done. Artifacts in $DIST/"
