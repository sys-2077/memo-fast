#!/bin/sh
# memo-fast installer / uninstaller
# Usage:
#   sh install.sh              # install
#   sh install.sh --uninstall  # uninstall
set -e

BINARY_NAME="memo-fast"
INSTALL_DIR="/usr/local/bin"
REPO="codexor/memo-fast-product"

# --- helpers ---

info() { printf "  %s\n" "$1"; }
error() { printf "  ERROR: %s\n" "$1" >&2; exit 1; }

detect_os() {
    case "$(uname -s)" in
        Darwin) echo "darwin" ;;
        Linux)  echo "linux" ;;
        *)      error "Unsupported OS: $(uname -s)" ;;
    esac
}

detect_arch() {
    case "$(uname -m)" in
        x86_64|amd64)  echo "amd64" ;;
        arm64|aarch64) echo "arm64" ;;
        *)             error "Unsupported architecture: $(uname -m)" ;;
    esac
}

latest_version() {
    curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" \
        | grep '"tag_name"' \
        | sed -E 's/.*"v?([^"]+)".*/\1/'
}

# --- uninstall ---

do_uninstall() {
    info "Uninstalling ${BINARY_NAME}..."

    if [ -f "${INSTALL_DIR}/${BINARY_NAME}" ]; then
        rm -f "${INSTALL_DIR}/${BINARY_NAME}" 2>/dev/null || \
            sudo rm -f "${INSTALL_DIR}/${BINARY_NAME}"
        info "Removed ${INSTALL_DIR}/${BINARY_NAME}"
    else
        info "Binary not found at ${INSTALL_DIR}/${BINARY_NAME}"
    fi

    info ""
    info "To clean up project configs, remove .mcp/memo-fast/ in each project."
    info "To remove git hooks, run: rm .git/hooks/post-commit (if memo-fast only)"
    info ""
    info "Uninstall complete."
}

# --- install ---

do_install() {
    OS="$(detect_os)"
    ARCH="$(detect_arch)"
    VERSION="${1:-$(latest_version)}"

    if [ -z "$VERSION" ]; then
        error "Could not determine latest version. Check https://github.com/${REPO}/releases"
    fi

    # Strip leading v if present
    VERSION="${VERSION#v}"

    ARCHIVE="memo-fast-${OS}-${ARCH}.tar.gz"
    URL="https://github.com/${REPO}/releases/download/v${VERSION}/${ARCHIVE}"

    info "Installing memo-fast v${VERSION} (${OS}/${ARCH})..."
    info ""

    TMPDIR="$(mktemp -d)"
    trap 'rm -rf "${TMPDIR}"' EXIT

    info "Downloading ${URL}..."
    curl -fsSL "${URL}" -o "${TMPDIR}/${ARCHIVE}" || \
        error "Download failed. Check if v${VERSION} exists at https://github.com/${REPO}/releases"

    info "Extracting..."
    tar -xzf "${TMPDIR}/${ARCHIVE}" -C "${TMPDIR}"

    info "Installing to ${INSTALL_DIR}/${BINARY_NAME}..."
    if [ -w "${INSTALL_DIR}" ]; then
        mv "${TMPDIR}/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
    else
        sudo mv "${TMPDIR}/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
    fi
    chmod +x "${INSTALL_DIR}/${BINARY_NAME}"

    info ""
    info "memo-fast v${VERSION} installed successfully."
    info ""
    info "Get started:"
    info "  cd your-project"
    info "  memo-fast init"
    info "  memo-fast index"
    info ""
}

# --- main ---

case "${1}" in
    --uninstall)
        do_uninstall
        ;;
    --version)
        do_install "$2"
        ;;
    *)
        do_install
        ;;
esac
