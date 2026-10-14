#!/bin/bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"
VERSION="$(/usr/libexec/PlistBuddy -c 'Print CFBundleShortVersionString' Resources/App-Info.plist)"
mkdir -p dist
for ARCH in arm64 x86_64; do
    OUTPUT_DIR="$ROOT/build/$ARCH" ARCH="$ARCH" ./scripts/build.sh
    APP="$ROOT/build/$ARCH/NewToFile.app"
    ditto -c -k --sequesterRsrc --keepParent "$APP" "dist/NewToFile-$VERSION-$ARCH.zip"
    STAGE="$(mktemp -d "$ROOT/build/dmg-$ARCH.XXXXXX")"
    ditto "$APP" "$STAGE/NewToFile.app"
    ln -s /Applications "$STAGE/Applications"
    hdiutil create -ov -format UDZO -volname NewToFile -srcfolder "$STAGE" \
        "dist/NewToFile-$VERSION-$ARCH.dmg"
    rm "$STAGE/Applications"
    rm -rf "$STAGE/NewToFile.app"
    rmdir "$STAGE"
done
(cd dist && shasum -a 256 "NewToFile-$VERSION-arm64.zip" "NewToFile-$VERSION-arm64.dmg" \
    "NewToFile-$VERSION-x86_64.zip" "NewToFile-$VERSION-x86_64.dmg" > SHA256SUMS.txt)
echo "Packages: $ROOT/dist"
