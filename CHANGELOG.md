# Changelog

## 0.2.2 - 2026-10-14

- Support a dedicated signing keychain and explicit timestamp policy.
- Publish certificate-self-signed App, helper, Finder extension and DMG previews;
  include the public certificate and fingerprint, never the private key.
- Build release bundles into fresh per-release directories to avoid stale files.
- Allow packaged-helper tests to select an extracted release app.
- Add a release verifier for certificate identity, nested signatures, version,
  executable architecture and DMG integrity.
- Preserve the existing v0.2.1 assets and distinguish certificate self-signing from
  ad-hoc signatures and Apple notarization.

## 0.2.1 - 2026-10-14

- Fix mixed appearances and unreadable sidebar text when the app's light/dark
  preference differs from macOS. Native vibrancy remains for system appearance;
  explicit appearance uses an opaque, theme-matched sidebar.
- Add opposite-system appearance regression tests and native renderer previews.
- Unify displayed and bundle versions.
- Isolate arm64 and x86_64 build output and add ZIP/DMG packaging with SHA-256.
- Add MIT license, contributor/security policies, English documentation, release
  instructions, and a bounded security review record.

## 0.2.0 - Unreleased Local Build

- Introduce the MyGo 0.4.8 native Go interface and Swift Finder extension/helper.
- Add text formats, Office/WPS detection, real OOXML templates, collision-safe
  creation, appearance preferences, and local automated tests.
