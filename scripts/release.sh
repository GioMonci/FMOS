#!/usr/bin/env bash
#
# Builds the files for an FMOS release into build/release/v<version>/:
#
#   fmos-x64.exe     Windows on Intel or AMD
#   fmos-arm64.exe   Windows on ARM
#   fmos-macos       macOS on Apple Silicon
#   SHA256SUMS.txt   the fingerprint (SHA-256) of each file above
#   notes.md         release notes for you to fill in (not uploaded)
#
# Every program file is built from committed code, with -trimpath (so no
# folder paths from this computer end up inside) and with the release version.
# The script checks its own work, then prints how to publish. It never
# uploads anything itself.
#
# Usage (from macOS or Linux): scripts/release.sh 0.1.0
# On Windows, use scripts/release.ps1, which does the same steps and checks.

set -euo pipefail

version="${1:-}"
if [[ ! "$version" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  echo "Usage: scripts/release.sh <version>   for example: scripts/release.sh 0.1.0" >&2
  exit 1
fi

cd "$(dirname "$0")/.."

if [[ -n "$(git status --porcelain)" ]]; then
  echo "There are uncommitted changes. Commit them first, so the release matches the code on GitHub." >&2
  exit 1
fi
if git rev-parse -q --verify "refs/tags/v$version" >/dev/null; then
  echo "The tag v$version already exists. Pick a new version number." >&2
  exit 1
fi

out="build/release/v$version"
if [[ -e "$out" ]]; then
  echo "$out already exists. Delete that folder first if you want to build this version again." >&2
  exit 1
fi
mkdir -p "$out"
files=(fmos-x64.exe fmos-arm64.exe fmos-macos)

# build <GOOS> <GOARCH> <file name>
build() {
  echo "Building $3..."
  CGO_ENABLED=0 GOOS="$1" GOARCH="$2" go build -trimpath \
    -ldflags "-X main.version=$version" -o "$out/$3" ./cmd/fmos
}
build windows amd64 fmos-x64.exe
build windows arm64 fmos-arm64.exe
build darwin arm64 fmos-macos

(cd "$out" && shasum -a 256 "${files[@]}" > SHA256SUMS.txt)

# Check 1: no folder paths from this computer inside the files.
for file in "${files[@]}"; do
  if grep -qF "$HOME" "$out/$file"; then
    echo "Stopped: $file still contains paths from this computer ($HOME). It was not built with -trimpath." >&2
    exit 1
  fi
done

# Check 2: the fingerprints in SHA256SUMS.txt match the files.
(cd "$out" && shasum -a 256 --check --quiet SHA256SUMS.txt)

# Check 3: the program reports the release version (only possible on an Apple Silicon Mac).
if [[ "$(uname -s)/$(uname -m)" == "Darwin/arm64" ]]; then
  reported="$("$out/fmos-macos" --version)"
  if [[ "$reported" != "FMOS $version" ]]; then
    echo "Stopped: fmos-macos reports \"$reported\" instead of \"FMOS $version\"." >&2
    exit 1
  fi
fi

cat > "$out/notes.md" <<EOF
## What's new in FMOS $version

- (Write a few short, plain-English lines. FMOS shows them on its Updates screen.)
EOF

commit="$(git rev-parse HEAD)"
echo
echo "Release files for FMOS $version are ready in $out:"
while read -r hash file; do
  printf '  %-15s  fingerprint starts with %s\n' "$file" "${hash:0:8}"
done < "$out/SHA256SUMS.txt"
echo
echo "Checked: built with -trimpath, fingerprints match, version is $version."
echo
echo "Next:"
echo "  1. Write the release notes in $out/notes.md"
echo "  2. Make sure this commit is on GitHub: git push"
echo "  3. Publish, either with the GitHub CLI:"
echo "       gh release create v$version --target $commit --title \"FMOS $version\" \\"
echo "         --notes-file $out/notes.md \\"
echo "         $out/fmos-x64.exe $out/fmos-arm64.exe $out/fmos-macos $out/SHA256SUMS.txt"
echo "     or on github.com: Releases > Draft a new release > tag v$version,"
echo "     paste notes.md, and attach the 3 programs and SHA256SUMS.txt."
