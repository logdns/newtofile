# Building and Releasing

## Prerequisites

Use a macOS host with Go 1.27.1+, Xcode/macOS SDK and Python 3. Dependencies are
pinned in `go.mod`/`go.sum`; use `go mod download` for initial provisioning.
Review dependency licenses in `Resources/ThirdParty` when upgrading them.

## Local Validation

```sh
go test ./...
go test -race ./...
go vet ./...
swift test
python3 scripts/test_templates.py
./scripts/build.sh
python3 scripts/test_app.py
./scripts/preview.sh
```

Validate both system appearances and both explicit opposite-system appearances.
The MyGo tester validates rendered content, not the live AppKit compositor.
Manually check Finder integration, privacy-denied folders, removable drives,
Chinese input, and actual Office/WPS opening before declaring production-ready.

For an externally sandboxed test run, set `NEWTOFILE_TEST_TMPDIR` to an allowed
disposable directory; Swift's Foundation temporary directory may ignore `TMPDIR`.
Swift 6.4's default SwiftBuild engine may request large sparse cache files; when
a file-size-limited validation environment blocks it, use
`swift test --build-system native` while that fallback is supported. If an outer
OS sandbox is already enforced, `--disable-sandbox` disables only SwiftPM's nested
sandbox, not the outer policy. Do not disable isolation to validate untrusted code.

## Packages

```sh
./scripts/package.sh
cd dist
shasum -a 256 -c SHA256SUMS.txt
```

The script builds separate `build/release-VERSION.RANDOM/arm64/NewToFile.app` and
`build/release-VERSION.RANDOM/x86_64/NewToFile.app`, then creates architecture-specific ZIP and DMG
files in `dist/`. These are separate packages, not a universal binary.
Building Intel on Apple Silicon does not prove Intel runtime compatibility.

The package script allocates fresh per-release/per-architecture output directories.
The lower-level build script still reuses a caller-specified existing bundle. Record toolchain versions and the exact source
commit alongside final package digests in the release manifest.
Verify all three Mach-O executables with `file` and `lipo -archs`, and verify
each bundle with `codesign --verify --deep --strict`.

The default is ad-hoc signing. Preview releases must say **not notarized** and
explain Gatekeeper/Finder extension limitations. Do not use unrelated signing
identities, claim notarization, or automatically disable macOS protections.

With your own Developer ID identity:

```sh
SIGN_IDENTITY='Developer ID Application: Your Name (TEAMID)' ./scripts/package.sh
```

For a production release, submit the ZIP using your locally configured
`xcrun notarytool submit --wait --keychain-profile PROFILE`, staple and validate
the app, rebuild the ZIP/DMG from the stapled app, notarize/staple the DMG as
needed, and **recompute checksums**. Validate Gatekeeper acceptance and extension
activation on a clean machine. Keep credentials outside the repository.

## GitHub Publication

Update `internal/appui/model.go`, both bundle plists, and `CHANGELOG.md` together.
Commit the reviewed source, tag it, and publish a release from that exact tag.
Attach ZIPs, DMGs, and `SHA256SUMS.txt`; GitHub supplies source archives.
Use a pre-release for unnotarized preview builds and state known validation gaps.
Do not reuse a published version/tag for different binaries. Download the
uploaded assets into a disposable directory and verify their checksums.

## Certificate Self-Signing

Use your own dedicated local code-signing identity, never another project's identity:

```sh
SIGN_IDENTITY=YOUR_CERTIFICATE_SHA1 \
SIGN_KEYCHAIN="$HOME/Library/Application Support/NewToFile/Signing/NewToFile-Release.keychain-db" \
SIGN_TIMESTAMP=none DIST_DIR=dist/0.2.2 ./scripts/package.sh
```

The keychain must be unlocked by its owner. Do not store the keychain password,
private key or PKCS#12 export in the repository or release. The script signs
the extension, helper, app and final DMGs, then computes checksums after signing.
Self-signatures validate byte integrity, not Apple trust. This is still a preview.
See [certificate policy](self-signing.md).

Verify a certificate-signed release with:

```sh
python3 scripts/verify_release.py --app PATH_TO_APP --arch arm64 \
  --version 0.2.2 --certificate docs/signing/NewToFile-SelfSigned.pem \
  --dmg dist/0.2.2/NewToFile-0.2.2-arm64.dmg
NEWTOFILE_APP=PATH_TO_APP python3 scripts/test_app.py
```
