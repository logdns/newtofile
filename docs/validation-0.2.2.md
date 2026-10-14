# v0.2.2 Self-Signed Preview Validation

Date: 2026-10-14. Host: Apple Silicon, macOS 27.0.1, Xcode 27.0,
Swift 6.4, Go 1.27.1. This is engineering acceptance, not a new security audit.
The historical security review in `security-audit.md` applies to its recorded
v0.2.1 snapshot, not the changed release tooling in this version.

## Automated Acceptance

- 20 Go test functions pass; all packages pass the race detector and `go vet`.
- 15 Swift XCTest cases pass with no failures/skips, including 32 concurrent
  creators and existing-file/symlink preservation.
- 2 Python OOXML structure tests and 2 packaged-helper tests pass. The helper
  from the self-signed Apple Silicon ZIP creates all nine formats when WPS is
  detected locally, including Chinese paths and collision/error handling.
- Explicit light-on-dark and dark-on-light native renderer previews succeed.
- The local package script builds both architectures in fresh directories and
  successfully signs the app, helper, Finder extension and final DMGs with the
  dedicated NewToFile certificate. No unrelated certificate is used.
- `verify_release.py` checks both bundles' versions, all three executable
  architectures, nested signatures, leaf-certificate SHA-256 and DMG integrity.
- Invalid requested versions and a modified signed helper are rejected in
  negative verification checks.

The final package verification, immutable source commit, CI run, signing identity,
asset sizes/digests and downloaded-asset results are recorded in the release
manifest. Source-only changes in release documentation do not substitute for
revalidating final package bytes.

## Expected Trust Limitation

A local `spctl --assess --type execute` check **rejects** the self-signed app
(exit 3). This is expected and not hidden: certificate-backed code integrity
passes, but Apple publisher trust/notarization does not. No trusted-root import,
Gatekeeper bypass or system-wide security change is performed. The certificate
is not Apple Developer ID, and no trusted signing timestamp is asserted.

## Still Not Established

Intel hardware execution, clean-machine Finder activation/TCC, full live AppKit
and accessibility acceptance, and actual Office/WPS opening/saving remain manual
acceptance items. Cross-compilation and OOXML structure are not those guarantees.
GitHub CI uses the build script's default ad-hoc signing; project-private signing
keys are not placed on a hosted runner.
