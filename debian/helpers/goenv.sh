#!/bin/sh
set -eu

__goos__deb_arch_os() {
	case "$1" in
		kfreebsd) echo freebsd ;;
		linux) echo "$1" ;;
		*) echo >&2 "error: unrecongized DEB_*_ARCH_OS: $1"; exit 1 ;;
	esac
}

__goarch__deb_arch_cpu() {
	case "$1" in
		amd64|arm|arm64|loong64|mips|ppc64|riscv64|s390x) echo "$1" ;;
		i386) echo 386 ;;
		mips64el) echo mips64le ;;
		mipsel) echo mipsle ;;
		ppc64el) echo ppc64le ;;
		*) echo >&2 "error: unrecongized DEB_*_ARCH_CPU: $1"; exit 1 ;;
	esac
}

#       build machine
#           The machine the package is built on.
#
#       host machine
#           The machine the package is built for.

DEB_HOST_ARCH="$(dpkg-architecture --query DEB_HOST_ARCH 2>/dev/null)"
# set _OS and _CPU explicitly via --force so they are never autodetected via CC (this makes it easier to test this script for various target architectures via a simple "DEB_HOST_ARCH=xxx ./debian/helpers/goenv.sh env")
DEB_HOST_ARCH_OS="$(dpkg-architecture --force --host-arch "$DEB_HOST_ARCH" --query DEB_HOST_ARCH_OS)"
DEB_HOST_ARCH_CPU="$(dpkg-architecture --force --host-arch "$DEB_HOST_ARCH" --query DEB_HOST_ARCH_CPU)"
export DEB_HOST_ARCH DEB_HOST_ARCH_OS DEB_HOST_ARCH_CPU

DEB_BUILD_ARCH_OS="$(dpkg-architecture --query DEB_BUILD_ARCH_OS)"
GOHOSTOS="$(__goos__deb_arch_os "$DEB_BUILD_ARCH_OS")"
GOOS="$(__goos__deb_arch_os "$DEB_HOST_ARCH_OS")"
export GOHOSTOS GOOS

DEB_BUILD_ARCH_CPU="$(dpkg-architecture --query DEB_BUILD_ARCH_CPU)"
GOHOSTARCH="$(__goarch__deb_arch_cpu "$DEB_BUILD_ARCH_CPU")"
GOARCH="$(__goarch__deb_arch_cpu "$DEB_HOST_ARCH_CPU")"
export GOHOSTARCH GOARCH

if [ -z "$GOHOSTOS" -o -z "$GOOS" -o -z "$GOHOSTARCH" -o -z "$GOARCH" ]; then
	exit 1
fi

# Avoid all "go" invocations downloading different toolchains during build.
export GOTOOLCHAIN=local

# Always not use sse2. This is important to ensure that the binaries we build
# (both when compiling golang on the buildds and when users cross-compile for
# 386) can actually run on older CPUs (where old means e.g. an AMD Athlon XP
# 2400+). See http://bugs.debian.org/753160 and
# https://code.google.com/p/go/issues/detail?id=8152
export GO386=softfloat

unset GOARM
if [ "$GOARCH" = 'arm' ]; then
	# start with GOARM=5 for maximum compatibility (see note about GO386 above)
	GOARM=5
	case "$DEB_HOST_ARCH" in
		armhf) GOARM=6 ;; # TODO detect Debian vs Raspbian and upgrade this to 7 by default?
	esac
	export GOARM
fi

# set CC_FOR_os_arch variables appropriately for supported architectures so that cross-compile even with cgo "just works" in more cases (also consistently across architectures for better reproducibility)
linuxArchList="$(dpkg-architecture --list-known --match-wildcard 'gnu-linux-any')"
# hotly contested/overlapping architectures: let's arbitrarily choose "armhf" as the cross-compile target for "GOARCH=arm" unless we're explicitly building the "armel" package
# this matches upstream's behavior for GOARM: https://github.com/golang/go/blob/go1.25.0/src/cmd/dist/util.go#L397-L405 (set to "7" if unspecified and cross-compiling)
armArchForCC='armhf'
if [ "$GOARCH" = 'arm' ]; then
	armArchForCC="$DEB_HOST_ARCH"
fi
for dpkgArch in $linuxArchList; do
	archCpu="$(dpkg-architecture --force --host-arch "$dpkgArch" --query DEB_HOST_ARCH_CPU 2>/dev/null)"
	if goArch="$(__goarch__deb_arch_cpu "$archCpu" 2>/dev/null)"; then
		if [ "$goArch" = 'arm' ] && [ "$dpkgArch" != "$armArchForCC" ]; then
			continue
		fi
		gnuType="$(dpkg-architecture --force --host-arch "$dpkgArch" --query DEB_HOST_GNU_TYPE 2>/dev/null)"
		export  "CC_FOR_linux_${goArch}=${gnuType}-gcc"
		export "CXX_FOR_linux_${goArch}=${gnuType}-g++"
		unset gnuType
	fi
	unset archCpu goArch
done
unset linuxArchList dpkgArch

unset CGO_ENABLED
if [ "$GOOS" = 'linux' ] && [ "$GOARCH" != "$GOHOSTARCH" ] && eval 'test -n "${CC_FOR_'"${GOHOSTOS}_${GOHOSTARCH}"':-}" && "${CC_FOR_'"${GOHOSTOS}_${GOHOSTARCH}"'}" --version > /dev/null'; then
	# if we're cross-compiling, let's check whether we can safely explicitly enable CGO (and whether we should)
	# https://github.com/golang/go/blob/go1.24.5/src/cmd/dist/build.go#L1779-L1792 ("cgoEnabled" map)
	# minus https://github.com/golang/go/blob/go1.24.5/src/cmd/dist/build.go#L1822-L1830 ("broken" map)
	# (we could use "--dist-tool" to pre-compile the "dist" tool to then use "dist list -json" to get "CgoSupported" instead of hard-coding this list, but then we need something that can parse JSON too, and all that just feels like way too much just to get this list of platforms which should enable CGO by default -Tianon)
	# $ go tool dist list -json | jq -r 'map(select(.CgoSupported and .GOOS == "linux") | .GOOS + "/" + .GOARCH) | (map(length) | max) as $max | map(. + (" " * ($max - length))) | join(" |\\\n") + " )"'
	case "$GOOS/$GOARCH" in
		linux/386      |\
		linux/amd64    |\
		linux/arm      |\
		linux/arm64    |\
		linux/loong64  |\
		linux/mips     |\
		linux/mips64   |\
		linux/mips64le |\
		linux/mipsle   |\
		linux/ppc64le  |\
		linux/riscv64  |\
		linux/s390x    )
			if eval 'test -n "${CC_FOR_'"${GOOS}_${GOARCH}"':-}" && "${CC_FOR_'"${GOOS}_${GOARCH}"'}" --version > /dev/null'; then
				export CGO_ENABLED=1
			fi
			;;
	esac
fi

exec "$@"
