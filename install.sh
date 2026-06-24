#!/bin/sh
# Installer for http-runner.
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/idesyatov/http-runner/master/install.sh | sh
#
# Environment overrides:
#   VERSION   release tag to install (e.g. v1.4.0). Default: latest release.
#   BINDIR    install directory. Default: /usr/local/bin.
set -eu

REPO="idesyatov/http-runner"
BIN="http-runner"
BINDIR="${BINDIR:-/usr/local/bin}"

err() {
	echo "error: $*" >&2
	exit 1
}

require() {
	command -v "$1" >/dev/null 2>&1 || err "required tool not found: $1"
}

fetch() {
	# fetch <url> <output-file>
	curl -fsSL --retry 3 --connect-timeout 10 -o "$2" "$1"
}

sha256_of() {
	if command -v sha256sum >/dev/null 2>&1; then
		sha256sum "$1" | awk '{print $1}'
	elif command -v shasum >/dev/null 2>&1; then
		shasum -a 256 "$1" | awk '{print $1}'
	else
		err "no sha256 tool found (need sha256sum or shasum)"
	fi
}

require curl
require tar

# Detect OS.
os=$(uname -s)
case "$os" in
	Linux) os="linux" ;;
	Darwin) os="darwin" ;;
	*) err "unsupported OS: $os" ;;
esac

# Detect architecture.
arch=$(uname -m)
case "$arch" in
	x86_64 | amd64) arch="amd64" ;;
	aarch64 | arm64) arch="arm64" ;;
	*) err "unsupported architecture: $arch" ;;
esac

# Resolve version (tag) to install.
tag="${VERSION:-}"
if [ -z "$tag" ]; then
	tag=$(curl -fsSL --retry 3 "https://api.github.com/repos/${REPO}/releases/latest" |
		grep '"tag_name":' | head -n1 | sed -E 's/.*"tag_name": *"([^"]+)".*/\1/')
	[ -n "$tag" ] || err "could not resolve the latest release tag"
fi

# Archive names use the version without the leading 'v'.
ver="${tag#v}"
archive="${BIN}_${ver}_${os}_${arch}.tar.gz"
base="https://github.com/${REPO}/releases/download/${tag}"

tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

echo "Downloading ${archive} (${tag})..."
fetch "${base}/${archive}" "${tmp}/${archive}"
fetch "${base}/checksums.txt" "${tmp}/checksums.txt"

# Verify checksum.
expected=$(grep " ${archive}\$" "${tmp}/checksums.txt" | awk '{print $1}')
[ -n "$expected" ] || err "checksum for ${archive} not found in checksums.txt"
actual=$(sha256_of "${tmp}/${archive}")
[ "$expected" = "$actual" ] || err "checksum mismatch: expected ${expected}, got ${actual}"

# Extract and install.
tar -xzf "${tmp}/${archive}" -C "$tmp"
[ -f "${tmp}/${BIN}" ] || err "binary ${BIN} not found in archive"
chmod +x "${tmp}/${BIN}"

if [ -w "$BINDIR" ]; then
	mv "${tmp}/${BIN}" "${BINDIR}/${BIN}"
else
	echo "Installing to ${BINDIR} requires elevated permissions..."
	sudo mv "${tmp}/${BIN}" "${BINDIR}/${BIN}"
fi

echo "Installed ${BIN} ${tag} to ${BINDIR}/${BIN}"
