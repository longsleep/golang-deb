#!/bin/bash
set -eux
GOARCH=$(./bin/go env GOARCH)

./bin/go build -o ./bin/readabihash -linkshared -ldflags="-r ''" debian/helpers/readabihash.go

TRIPLET=$(dpkg-architecture -qDEB_HOST_MULTIARCH)

mkdir -p debian/libgolang-${GOVER}-std1/usr/lib/${TRIPLET}
mv pkg/linux_${GOARCH}_dynlink/libstd.so debian/libgolang-${GOVER}-std1/usr/lib/${TRIPLET}/libgolang-${GOVER}-std.so.1

ln -s ../../../${TRIPLET}/libgolang-${GOVER}-std.so.1 pkg/linux_${GOARCH}_dynlink/libstd.so

mkdir -p debian/golang-${GOVER}-go-shared-dev/usr/lib/go-${GOVER}/pkg/
mv pkg/linux_${GOARCH}_dynlink/ debian/golang-${GOVER}-go-shared-dev/usr/lib/go-${GOVER}/pkg/

cp bin/readabihash debian/golang-${GOVER}-go-shared-dev/usr/lib/go-${GOVER}/pkg/
