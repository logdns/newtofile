# Security Policy

## Supported Versions

Security fixes target the newest published version. Upgrade before reporting a
problem. The project does not currently promise an SLA or backports.

## Reporting a Vulnerability

Use GitHub's private vulnerability reporting under this repository's Security
tab when enabled. If it is unavailable, contact the maintainer through their
GitHub profile to arrange a private channel; do not post exploit details or
personal data in a public issue. Ordinary UI bugs belong in Issues.

Include the version, macOS version, source locations, expected boundary,
minimum reproduction with dummy files, and observed impact. Do not test other
users' files or probe shared systems. We will acknowledge and assess reports
on a best-effort basis and coordinate disclosure after remediation.

## Security Model

- The app is an unprivileged local utility, not a service or privileged daemon.
- External `newtofile://` links only prefill the visible form. Creation requires
  a user action; the app does not authorize writes on behalf of another user.
- The bundled Swift helper accepts JSON via stdin. Go invokes a fixed bundle
  executable without a shell. The helper has the same user's filesystem rights.
- Files use exclusive creation with no-follow protection on the final component.
  Existing files and symlinks are not overwritten. Parent directories follow
  normal filesystem semantics; this is not a hostile-directory isolation tool.
- The Finder extension is sandboxed; the main app/helper are not App Sandbox
  apps. macOS filesystem permissions and privacy controls still apply.
- No app telemetry, update downloader, or network-backed feature is implemented.

## Distribution Limits

v0.2.2 previews use a dedicated NewToFile self-signed code-signing certificate;
v0.2.1 uses ad-hoc signing. Both are **not Apple notarized**. The self-signed
certificate is not an Apple Developer ID or a system-trusted publisher identity.
Gatekeeper may block them and Finder extension registration may differ across
systems. Checksums detect transfer corruption; they are not an independent
publisher identity. Prefer a source build if this distribution model does not
meet your security requirements. Never disable Gatekeeper globally.

See [audit scope and limits](docs/security-audit.md). A source review and passing
tests do not guarantee absence of vulnerabilities or clean-machine compatibility.

The public signing certificate and fingerprint are documented in
[Self-Signing](docs/self-signing.md). Private keys/passwords are never shipped.
Do not import the certificate as a trusted root merely to install this preview.
