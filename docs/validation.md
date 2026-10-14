# v0.2.1 Validation Record

Date: 2026-10-14. Host: macOS / Apple Silicon, Go 1.27.1, Swift 6.4, Xcode 27.0.

## Passed

| Check | Result |
| --- | --- |
| Go tests | 20 tests; all packages pass |
| Go race tests | All packages pass |
| `go vet ./...` | Pass |
| Swift core tests | 15 XCTest cases, no failures/skips |
| OOXML structural tests | 2 Python cases pass |
| Packaged helper tests | 2 Python cases pass; WPS detected, all nine formats created |
| Appearance previews | System light/dark and explicit opposite-system light/dark rendered |
| Appearance regression | Sidebar opaque/theme matched with native-material availability and opposing system appearance |

Go/Swift/helper engineering acceptance runs use an offline OS file/network policy,
empty allowlisted environment, scratch-local caches, disposable directories and
resource/time limits. These are separate from the source-only security review;
the host could not establish that review's required hard memory cap. They are not
vulnerability reproduction evidence or a claim of full audit sandbox compliance.

Swift 6.4's default SwiftBuild cache initialization hit the validation file-size
cap. The tests pass using the still-supported native build system, with SwiftPM's
nested sandbox disabled under the outer OS policy. Foundation's temporary-directory
choice required `NEWTOFILE_TEST_TMPDIR` to keep fixtures inside scratch. Initial
sandbox-denied temp/cache writes are environment failures, not silently counted as
successful tests. Apple `xcrun` printed denied lookup-cache warnings; final compiler,
test and signature commands returned success. Disk image packaging is performed
separately with Apple's host tooling, not as security-audit execution.

## Release Verification

Both architectures are built into fresh independent directories. Release assets
include a source/toolchain/digest manifest and SHA-256 checksums. Nested signatures,
architecture metadata, ZIP contents, DMG integrity and downloaded GitHub asset
checksums are verified during publication; the release manifest records results.

## Not Established

- Developer ID publisher authentication, Apple notarization or Gatekeeper acceptance.
- Clean-machine Finder extension activation and effective TCC permission prompts.
- Intel hardware execution (x86_64 is cross-compiled).
- Live AppKit vibrancy/window composition and VoiceOver acceptance.
- Actual document opening/saving in Microsoft Office or WPS. Structural integrity
  and creating the template copies are not an application-compatibility guarantee.
- Full dependency memory-safety/advisory analysis or native single-instance IPC audit.

This release remains an explicitly unnotarized preview, not production certified.
