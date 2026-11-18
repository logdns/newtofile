# v0.2.4 Self-Signed Preview Validation

Date: 2026-11-18. Local validation used an Apple Silicon Mac running macOS
27.2.1, Xcode 27.0 / Swift 6.4, Go 1.27.1 and Python 3.14.7.

## Engineering Checks

- Go tests, race checks, and `go vet` passed.
- All 19 Swift tests passed using `swift test --build-system native`, including
  four regression tests for Microsoft Office and WPS detection. The tests
  cover registered identifiers, application-bundle fixtures in both standard
  directory scopes, missing or mismatched applications, and format mapping.
- Both template checks passed. The packaged helper was tested with installed
  WPS, including creation of DOCX, XLSX, and PPTX from bundled templates.
- [Scoped security review](security-review-0.2.4.md) covers this detection
  change; it does not replace a full audit.

## Release Checks

The release manifest identifies the source commit, toolchain, assets and
public certificate fingerprint. The release checks cover both architectures,
the signed App and nested executables, the ZIP and DMG payloads, and the
downloaded SHA-256 checksums.

The Intel package is cross-built, not tested on Intel hardware. No real
Microsoft Office installation, WPS opening/saving, or clean-machine Finder
activation was tested. The preview is certificate-self-signed and not Apple
notarized. macOS may block its launch or extension activation.
