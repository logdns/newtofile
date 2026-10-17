# v0.2.3 Scoped Security Review

Date: 2026-10-17. Reviewed source commit `b5a0abdcb58afa6fce8b1c74449c43f56b18edb9`.
Profile: quick, scoped, source-only review of the About links, their UI regression
test, version/plist metadata, changelog and public project links. This is **not** a
full or comprehensive audit, penetration test, dependency advisory scan, or review
of hosted GitHub policies.

## Result

- Confirmed vulnerabilities: 0.
- Unresolved source-grounded vulnerability hypotheses: 0.
- Scope: fixed About links and associated public metadata only. Unchanged app/file
  creation, IPC/Finder paths, dependencies, CI/release service, and OS behavior
  were not re-audited. The broader, older quick review is documented in
  [security-audit.md](security-audit.md), but is not a substitute for this scope.

## Source Trace

`internal/appui/settings.go` declares the GitHub and xinai.de destinations as
compile-time HTTPS constants and passes them directly to MyGo `ui.Link`. There is
no user-controlled URL, custom scheme, shell execution, file access, credential,
or authorization decision in this feature. `internal/appui/view_test.go` clicks
both links and verifies the exact URLs observed by MyGo's headless browser-opening
capture.

The release uses a distinct version-specific self-signed certificate. The
certificate verifies artifact integrity relative to its public key, but does not
establish Apple Developer ID trust or notarization. The private key remains in a
local maintainer keychain; only the public certificate is packaged.

## Validation Limits

The audit did not execute target-controlled code under the audit sandbox, probe
live websites or GitHub, or validate browser policy. Separate engineering tests
and build checks are listed in [v0.2.3 validation](validation-0.2.3.md); they are
not audit sandbox evidence. Website ownership/availability, default-browser
behavior, clean-machine Gatekeeper/Finder activation, Intel runtime and actual
Office/WPS acceptance remain outside this source review.

Machine-readable scoped records and validators are retained outside the
repository at `/Users/ideadev/security-audit-skill/newtofile/run-2/` on the
maintainer's workstation; they are not release assets.
