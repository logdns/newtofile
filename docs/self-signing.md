# NewToFile Certificate Self-Signing

## v0.2.3 Identity Rotation

v0.2.3 uses a new, version-specific local signing identity:
`NewToFile macOS Self-Signed Release v0.2.3`. Its public certificate is
[NewToFile-SelfSigned-0.2.3.pem](signing/NewToFile-SelfSigned-0.2.3.pem), with
SHA-256 fingerprint:

`11:ED:CF:E7:34:C8:8F:C7:F4:CD:64:D5:AD:00:DD:4E:CE:A8:4B:F5:9F:7F:72:2E:5B:5F:83:04:AF:7C:C6:92`

The v0.2.2 certificate and its already-published assets remain unchanged. The
private key and release keychain are maintainer-only local material and are
never uploaded. Like the previous identity, this certificate does not establish
Apple publisher trust or notarization.

Starting with v0.2.2, preview assets use the dedicated identity
`NewToFile macOS Self-Signed Release` (subject organization `logdns`). The
v0.2.1 preview used ad-hoc signatures without a signing certificate.

- RSA 3072, SHA-256, code-signing EKU; certificate is not a CA.
- Validity: 2026-10-14 09:05:31 UTC to 2031-10-13 09:05:31 UTC.
- Certificate SHA-256 fingerprint:
  `C6:C4:E2:25:83:8F:6A:D5:67:44:19:F7:F3:36:62:4D:10:2B:B0:B8:9F:70:26:29:59:36:80:68:CB:82:04:C4`.
- Public certificate: [NewToFile-SelfSigned.pem](signing/NewToFile-SelfSigned.pem).
- App, helper, Finder extension and DMG use the same certificate. No trusted
  timestamp is claimed; self-signed previews explicitly disable timestamping.

## What It Means

This is a real certificate-backed signature, not `codesign -s -` ad-hoc signing.
It provides cryptographic byte integrity relative to this certificate and a
stable project signing identity. **It is not Apple Developer ID signing,
not Apple notarization, and does not establish system publisher trust.**
Gatekeeper rejection is expected, and Finder extension activation still needs
clean-machine acceptance testing. Signature verification is not app compatibility.

A certificate/checksum obtained alongside the artifact is not an independent
trust root. Confirm the project's source/channel before downloading. Never
install this certificate as a trusted root or disable Gatekeeper globally merely
to open a preview. Prefer a source build if the distribution model is unsuitable.

## Inspect a Download

```sh
openssl x509 -in NewToFile-SelfSigned.pem -noout -fingerprint -sha256
codesign --verify --deep --strict NewToFile.app
codesign -dv --verbose=4 NewToFile.app
```

Expect `Authority=NewToFile macOS Self-Signed Release`, not `Signature=adhoc`.
Extract the signing certificate in a disposable directory and compare the leaf:

```sh
CHECK_DIR="$(mktemp -d)"
codesign --display --extract-certificates="$CHECK_DIR/cert" NewToFile.app
openssl x509 -inform DER -in "$CHECK_DIR/cert0" -noout -fingerprint -sha256
```
This does not guarantee Gatekeeper will accept the app.

## Maintainer Key Handling

The identity lives in a dedicated local keychain outside the repository, not
in the login/system keychain search list. No system/user trust override is added.
Local keychain recovery material is restricted to the current account; protect
that account, encrypt backups, and do not upload recovery files. Public releases
contain only the public certificate. A future identity rotation must be explicitly
documented; never silently replace an already published version's signed bytes.
