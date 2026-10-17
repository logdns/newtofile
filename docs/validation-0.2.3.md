# v0.2.3 Self-Signed Preview Validation

Date: 2026-10-17. This record covers NewToFile's native About links for the
GitHub project and xinai.de. The release remains a self-signed, unnotarized
preview, not an Apple-trusted distribution.

## Changes

- Added native About-section links with fixed HTTPS destinations.
- Added a MyGo UI regression test that clicks both links and verifies the exact
  URLs dispatched to the default-browser handler.
- Incremented app/Finder bundle version to 0.2.3 (build 5).
- Rotated to a new, version-specific self-signed certificate; v0.2.2 assets and
  certificate are not modified.

## Engineering Validation

On Apple Silicon macOS 27.0.1 with Go 1.27.1, Swift 6.4, Xcode 27.0 and
Python 3.14:

- `go test ./...`, `go test -race ./...`, `go vet ./...` passed.
- `swift test --build-system native` passed all 15 XCTest cases. SwiftPM reports
  that the native build-system option is deprecated; it was used as documented
  for this toolchain's sparse-cache behavior.
- `python3 scripts/test_templates.py` passed 2 checks; `python3 scripts/test_app.py`
  passed 2 packaged-helper integration tests using the arm64 app.
- `./scripts/preview.sh` rendered MyGo previews; `./scripts/build.sh` built and
  codesign-verified the local ad-hoc development bundle.
- The arm64 and x86_64 release bundles passed the version, executable
  architecture, nested signature and v0.2.3 certificate verifier; both DMGs
  passed signature and image-integrity checks. The packaged helper integration
  test passed on arm64. Intel helper execution is not supported by this Apple
  Silicon host (`Bad CPU type in executable`); Intel runtime tests are not
  claimed.
- Final release source commit, build manifest, public certificate fingerprint
  and package SHA-256 values are recorded in `release-manifest.json`,
  `NewToFile-SelfSigned-0.2.3.pem` and `SHA256SUMS.txt`.

## Security Review

A quick scoped source-only review of the fixed URL links found no confirmed
security issues or unresolved source-grounded hypotheses. This is not a full
re-audit of the app, Finder extension, dependencies, or hosted release services;
see the retained security review record for broader review history and limits.

## Known Gaps

- Intel artifacts are cross-compiled and signature/architecture checked, but
  cannot execute on this Apple Silicon host; Intel hardware validation remains.
- No clean-machine Gatekeeper, Finder extension activation, Office/WPS opening,
  or accessibility acceptance was performed.
- The self-signed certificate is not Apple Developer ID signing and is not
  notarized. Gatekeeper may reject this preview. No trusted timestamp is claimed.
