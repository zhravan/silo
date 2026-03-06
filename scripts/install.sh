#!/bin/sh
# Install silo from GitHub releases.
# Usage: curl -sSL https://raw.githubusercontent.com/zhravan/silo/main/scripts/install.sh | sh
# Or: SILO_VERSION=0.0.1 curl -sSL ... | sh

set -e

REPO="zhravan/silo"
GITHUB_API="https://api.github.com/repos/${REPO}"
RELEASES="${GITHUB_API}/releases"

# Resolve version (default: latest release tag)
get_latest_version() {
  curl -sSL "${RELEASES}/latest" | sed -n 's/.*"tag_name": *"v\?\([^"]*\)".*/\1/p' | head -1
}

VERSION="${SILO_VERSION:-$(get_latest_version)}"
if [ -z "$VERSION" ]; then
  echo "Could not determine version. Set SILO_VERSION or ensure ${REPO} has a release." >&2
  exit 1
fi

# Normalize: ensure v prefix for URL
TAG="v${VERSION#v}"
BASE_URL="https://github.com/${REPO}/releases/download/${TAG}"

# Detect OS and arch
OS=$(uname -s)
ARCH=$(uname -m)

case "$OS" in
  Linux)   OS=linux ;;
  Darwin)  OS=darwin ;;
  *)
    echo "Unsupported OS: $OS" >&2
    exit 1
    ;;
esac

case "$ARCH" in
  x86_64)   ARCH=amd64 ;;
  amd64)    ARCH=amd64 ;;
  aarch64)  ARCH=arm64 ;;
  arm64)    ARCH=arm64 ;;
  armv6l|armv7l) ARCH=arm ;;
  *)
    echo "Unsupported arch: $ARCH" >&2
    exit 1
    ;;
esac

BINARY="silo-${OS}-${ARCH}"
URL="${BASE_URL}/${BINARY}"
echo "Installing silo ${TAG} (${OS}/${ARCH}) from ${URL}"

TMP=$(mktemp -d)
trap 'rm -rf "$TMP"' EXIT
curl -sSL -o "${TMP}/silo" "$URL"
chmod +x "${TMP}/silo"

PREFIX="${SILO_PREFIX:-/usr/local}"
DEST="${PREFIX}/bin/silo"
if [ ! -w "${PREFIX}/bin" ] 2>/dev/null; then
  echo "Need write access to ${PREFIX}/bin (use: sudo)" >&2
  if command -v sudo >/dev/null 2>&1; then
    sudo mkdir -p "${PREFIX}/bin"
    sudo mv "${TMP}/silo" "$DEST"
  else
    mkdir -p "${HOME}/.local/bin"
    DEST="${HOME}/.local/bin/silo"
    mv "${TMP}/silo" "$DEST"
    echo "Installed to ${DEST}. Add ${HOME}/.local/bin to PATH if needed."
  fi
else
  mkdir -p "${PREFIX}/bin"
  mv "${TMP}/silo" "$DEST"
fi

echo "Installed to ${DEST}"
"$DEST" --help >/dev/null 2>&1 || true
