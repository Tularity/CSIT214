#!/bin/sh
# Builds the single-file release binaries into bin/: the React client is built
# and embedded, and each target is cross-compiled without cgo.
set -eu

root=$(cd "$(dirname "$0")/.." && pwd)
dist="$root/backend/web/dist"

cleanup() {
  find "$dist" -mindepth 1 ! -name .gitkeep -exec rm -rf {} +
}
trap cleanup EXIT

echo "building the client"
(cd "$root/frontend" && npm ci --no-audit --no-fund && npm run build)
cleanup
cp -R "$root/frontend/dist/." "$dist/"

mkdir -p "$root/bin"
for target in \
  "darwin arm64 csit214-macos-apple-silicon" \
  "darwin amd64 csit214-macos-intel" \
  "windows amd64 csit214-windows.exe" \
  "linux amd64 csit214-linux"
do
  set -- $target
  echo "building bin/$3"
  (cd "$root/backend" && \
    CGO_ENABLED=0 GOOS=$1 GOARCH=$2 go build -trimpath \
      -ldflags "-s -w -X main.releaseBuild=true" \
      -o "$root/bin/$3" .)
done

ls -l "$root/bin"
