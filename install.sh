#!/bin/sh
# ccinspect installer — works on macOS, Linux, and WSL
# Usage: curl -fsSL https://raw.githubusercontent.com/thieung/ccinspect/main/install.sh | sh

set -e

REPO="thieung/ccinspect"
INSTALL_DIR="/usr/local/bin"
BINARY="ccinspect"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

info()  { printf "${CYAN}▸${NC} %s\n" "$1"; }
ok()    { printf "${GREEN}✓${NC} %s\n" "$1"; }
warn()  { printf "${YELLOW}⚠${NC} %s\n" "$1"; }
error() { printf "${RED}✗${NC} %s\n" "$1" >&2; exit 1; }

# Detect OS
detect_os() {
  case "$(uname -s)" in
    Linux*)  echo "linux" ;;
    Darwin*) echo "darwin" ;;
    *)       error "Unsupported OS: $(uname -s). Use Windows? Try: install.ps1" ;;
  esac
}

# Detect architecture
detect_arch() {
  case "$(uname -m)" in
    x86_64|amd64) echo "amd64" ;;
    arm64|aarch64) echo "arm64" ;;
    *)             error "Unsupported architecture: $(uname -m)" ;;
  esac
}

# Get latest version from GitHub
get_latest_version() {
  if command -v curl > /dev/null 2>&1; then
    curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"v?([^"]+)".*/\1/'
  elif command -v wget > /dev/null 2>&1; then
    wget -qO- "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"v?([^"]+)".*/\1/'
  else
    error "Neither curl nor wget found. Please install one of them."
  fi
}

# Download file
download() {
  local url="$1" dest="$2"
  if command -v curl > /dev/null 2>&1; then
    curl -fsSL "$url" -o "$dest"
  elif command -v wget > /dev/null 2>&1; then
    wget -q "$url" -O "$dest"
  fi
}

main() {
  printf "\n${CYAN}╭─────────────────────────────────╮${NC}\n"
  printf "${CYAN}│${NC}   ${GREEN}ccinspect${NC} installer            ${CYAN}│${NC}\n"
  printf "${CYAN}╰─────────────────────────────────╯${NC}\n\n"

  OS=$(detect_os)
  ARCH=$(detect_arch)
  info "Detected: ${OS}/${ARCH}"

  # Get version
  VERSION="${CCINSPECT_VERSION:-}"
  if [ -z "$VERSION" ]; then
    info "Fetching latest version..."
    VERSION=$(get_latest_version)
  fi

  if [ -z "$VERSION" ]; then
    error "Could not determine latest version. Set CCINSPECT_VERSION manually."
  fi
  ok "Version: v${VERSION}"

  # Build download URL
  ARCHIVE="${BINARY}_${VERSION}_${OS}_${ARCH}.tar.gz"
  URL="https://github.com/${REPO}/releases/download/v${VERSION}/${ARCHIVE}"
  info "Downloading ${ARCHIVE}..."

  # Download to temp dir
  TMP_DIR=$(mktemp -d)
  trap 'rm -rf "$TMP_DIR"' EXIT

  download "$URL" "${TMP_DIR}/${ARCHIVE}" || error "Download failed. Check https://github.com/${REPO}/releases for available versions."
  ok "Downloaded"

  # Extract
  info "Extracting..."
  tar -xzf "${TMP_DIR}/${ARCHIVE}" -C "$TMP_DIR"
  ok "Extracted"

  # Install
  info "Installing to ${INSTALL_DIR}..."
  if [ -w "$INSTALL_DIR" ]; then
    mv "${TMP_DIR}/${BINARY}" "${INSTALL_DIR}/${BINARY}"
  else
    warn "Need sudo to write to ${INSTALL_DIR}"
    sudo mv "${TMP_DIR}/${BINARY}" "${INSTALL_DIR}/${BINARY}"
  fi
  chmod +x "${INSTALL_DIR}/${BINARY}"
  ok "Installed to ${INSTALL_DIR}/${BINARY}"

  # Verify
  printf "\n"
  if command -v "$BINARY" > /dev/null 2>&1; then
    ok "$(${BINARY} --version)"
    printf "\n${GREEN}🎉 Installation complete!${NC}\n"
    printf "   Run ${CYAN}ccinspect --help${NC} to get started.\n\n"
  else
    warn "Binary installed but not found in PATH."
    printf "   Add this to your shell profile:\n"
    printf "   ${CYAN}export PATH=\"${INSTALL_DIR}:\$PATH\"${NC}\n\n"
  fi
}

main
