#!/bin/bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"
OUTPUT_DIR="${OUTPUT_DIR:-$ROOT/build}"
mkdir -p "$OUTPUT_DIR"
OUTPUT_DIR="$(cd "$OUTPUT_DIR" && pwd)"
APP="$OUTPUT_DIR/NewToFile.app"
EXT="$APP/Contents/PlugIns/NewToFileFinder.appex"
SDK="$(xcrun --sdk macosx --show-sdk-path)"
ARCH="${ARCH:-$(uname -m)}"
SIGN_IDENTITY="${SIGN_IDENTITY:--}"
SIGN_ARGS=(--force --sign "$SIGN_IDENTITY" --options runtime)
if [[ -n "${SIGN_KEYCHAIN:-}" ]]; then SIGN_ARGS+=(--keychain "$SIGN_KEYCHAIN"); fi
if [[ -n "${SIGN_TIMESTAMP:-}" ]]; then SIGN_ARGS+=(--timestamp="$SIGN_TIMESTAMP"); fi
case "$ARCH" in
    arm64) GOARCH=arm64 ;;
    x86_64) GOARCH=amd64 ;;
    *) echo "Unsupported architecture: $ARCH" >&2; exit 1 ;;
esac

mkdir -p "$APP/Contents/MacOS" "$APP/Contents/Resources/Templates" "$EXT/Contents/MacOS"
python3 scripts/generate_templates.py
cp Resources/App-Info.plist "$APP/Contents/Info.plist"
cp Resources/Finder-Info.plist "$EXT/Contents/Info.plist"
cp Resources/Templates/blank.* "$APP/Contents/Resources/Templates/"
if [[ -d Resources/ThirdParty ]]; then
    mkdir -p "$APP/Contents/Resources/ThirdParty"
    chmod -R u+w "$APP/Contents/Resources/ThirdParty"
    cp -R Resources/ThirdParty/. "$APP/Contents/Resources/ThirdParty/"
    chmod -R u+w "$APP/Contents/Resources/ThirdParty"
fi
xcrun swift scripts/make_icon.swift "$OUTPUT_DIR/AppIcon.iconset"
iconutil -c icns "$OUTPUT_DIR/AppIcon.iconset" -o "$APP/Contents/Resources/AppIcon.icns"

COMMON=(-sdk "$SDK" -target "$ARCH-apple-macosx13.0" -swift-version 5 -O)
xcrun swiftc "${COMMON[@]}" -module-name NewToFileBridge \
    Sources/NewToFileCore/*.swift Sources/FileHelper/main.swift \
    -framework AppKit -o "$APP/Contents/MacOS/NewToFileBridge"
GOOS=darwin GOARCH="$GOARCH" CGO_ENABLED=0 go build -trimpath \
    -ldflags "-s -w -X github.com/egoist/mygo.packageURLSchemes=newtofile" \
    -o "$APP/Contents/MacOS/NewToFile" ./cmd/newtofile
xcrun swiftc "${COMMON[@]}" -module-name NewToFileFinder \
    -parse-as-library -application-extension \
    Sources/NewToFileCore/FileKind.swift Sources/NewToFileCore/CreationRequest.swift \
    Sources/NewToFileCore/OfficeAvailability.swift Sources/FinderExtension/FinderSync.swift \
    -framework Cocoa -framework FinderSync -Xlinker -e -Xlinker _NSExtensionMain \
    -o "$EXT/Contents/MacOS/NewToFileFinder"

codesign "${SIGN_ARGS[@]}" \
    --entitlements Resources/Finder.entitlements "$EXT"
codesign "${SIGN_ARGS[@]}" "$APP/Contents/MacOS/NewToFileBridge"
codesign "${SIGN_ARGS[@]}" "$APP"
codesign --verify --deep --strict "$APP"
echo "Built: $APP"
