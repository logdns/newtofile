#!/bin/bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"
VERSION="$(/usr/libexec/PlistBuddy -c 'Print CFBundleShortVersionString' Resources/App-Info.plist)"
DIST="${DIST_DIR:-$ROOT/dist}"
mkdir -p "$DIST" "$ROOT/build"
DIST="$(cd "$DIST" && pwd)"
BUILD_ROOT="$(mktemp -d "$ROOT/build/release-$VERSION.XXXXXX")"
SIGN_ARGS=(--force --sign "${SIGN_IDENTITY:--}")
if [[ -n "${SIGN_KEYCHAIN:-}" ]]; then SIGN_ARGS+=(--keychain "$SIGN_KEYCHAIN"); fi
if [[ -n "${SIGN_TIMESTAMP:-}" ]]; then SIGN_ARGS+=(--timestamp="$SIGN_TIMESTAMP"); fi
for ARCH in arm64 x86_64; do
    OUTPUT_DIR="$BUILD_ROOT/$ARCH" ARCH="$ARCH" ./scripts/build.sh
    APP="$BUILD_ROOT/$ARCH/NewToFile.app"
    ditto -c -k --sequesterRsrc --keepParent "$APP" "$DIST/NewToFile-$VERSION-$ARCH.zip"
    STAGE="$(mktemp -d "$ROOT/build/dmg-$ARCH.XXXXXX")"
    ditto "$APP" "$STAGE/NewToFile.app"
    ln -s /Applications "$STAGE/Applications"
    hdiutil create -ov -format UDZO -volname NewToFile -srcfolder "$STAGE" \
        "$DIST/NewToFile-$VERSION-$ARCH.dmg"
    codesign "${SIGN_ARGS[@]}" "$DIST/NewToFile-$VERSION-$ARCH.dmg"
    codesign --verify --strict "$DIST/NewToFile-$VERSION-$ARCH.dmg"
    rm "$STAGE/Applications"
    rm -rf "$STAGE/NewToFile.app"
    rmdir "$STAGE"
done
(cd "$DIST" && shasum -a 256 "NewToFile-$VERSION-arm64.zip" "NewToFile-$VERSION-arm64.dmg" \
    "NewToFile-$VERSION-x86_64.zip" "NewToFile-$VERSION-x86_64.dmg" > SHA256SUMS.txt)
echo "Builds: $BUILD_ROOT"
echo "Packages: $DIST"
