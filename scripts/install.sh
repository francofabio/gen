#!/usr/bin/env bash
set -e

# Install gen from GitHub Releases. No sudo. Prefer ~/.local/bin or $HOME/bin.
# Usage: curl -fsSL https://raw.githubusercontent.com/<user>/gen/main/scripts/install.sh | sh

BINARY_NAME="gen"
REPO_URL="${GEN_REPO_URL:-https://github.com/francofabio/gen}"
BASE_URL="${REPO_URL}/releases/latest/download"

detect_arch() {
	local os="$1"
	local arch="$2"
	case "$os" in
		Darwin)
			case "$arch" in
				arm64|aarch64) echo "darwin_arm64" ;;
				x86_64|amd64) echo "darwin_amd64" ;;
				*) echo "unsupported" ;;
			esac ;;
		Linux)
			case "$arch" in
				aarch64|arm64) echo "linux_arm64" ;;
				x86_64|amd64) echo "linux_amd64" ;;
				*) echo "unsupported" ;;
			esac ;;
		*)
			echo "unsupported"
			;;
	esac
}

main() {
	local os
	os="$(uname -s)"
	local arch
	arch="$(uname -m)"
	local target
	target="$(detect_arch "$os" "$arch")"
	if [[ "$target" == "unsupported" ]]; then
		echo "gen: unsupported platform: $os/$arch" >&2
		exit 1
	fi

	local archive_name="gen_${target}.tar.gz"
	local url="${BASE_URL}/${archive_name}"
	local tmpdir
	tmpdir="$(mktemp -d)"
	trap 'rm -rf "$tmpdir"' EXIT

	echo "Downloading $url ..." >&2
	if ! curl -fsSL "$url" -o "$tmpdir/$archive_name"; then
		echo "gen: download failed. Check the URL or your connection." >&2
		exit 1
	fi
	tar -xzf "$tmpdir/$archive_name" -C "$tmpdir"

	local bindir
	if [[ -n "$HOME" ]]; then
		if [[ -d "$HOME/.local/bin" ]]; then
			bindir="$HOME/.local/bin"
		elif [[ -d "$HOME/bin" ]]; then
			bindir="$HOME/bin"
		else
			bindir="$HOME/.local/bin"
			mkdir -p "$bindir"
		fi
	else
		echo "gen: HOME environment variable is not set" >&2
		exit 1
	fi

	cp "$tmpdir/$BINARY_NAME" "$bindir/$BINARY_NAME"
	chmod +x "$bindir/$BINARY_NAME"
	echo "Installed at $bindir/$BINARY_NAME" >&2

	if ! echo ":$PATH:" | grep -q ":$bindir:"; then
		echo "" >&2
		echo "The directory $bindir is not in your PATH. Add it to your shell:" >&2
		echo "  export PATH=\"\$HOME/.local/bin:\$PATH\"" >&2
		echo "" >&2
	fi
}

main "$@"
