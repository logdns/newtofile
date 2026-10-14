# NewToFile Security Review — 2026-10-14

## Scope and Result

**Quick, partial, source-only review** of clean source commit
`de0fe4505515c074e1106a70902d129b9366db49`. This is not an exhaustive audit,
penetration test, dependency advisory scan, or assurance of production readiness.
No prior compatible ledgers were available. No invocation budget was set:
four independent reconnaissance calls, two hunters in one wave, one fresh
final coverage critic; verifier capacity was reserved but no candidate survived.
Candidate validation/final record verification had an empty candidate set.

Reviewed first-party Go/Swift app, helper, Finder extension, preferences,
file creation, bundled template source/resources, dependencies' pinning inputs,
CI, build and package scripts. **0 confirmed findings; 0 needs-validation
vulnerability records.** This does not mean there are no vulnerabilities.

## Controls Observed in Source

- External custom-scheme links only prefill a visible form; creation requires a
  separate user action. Pending operations snapshot their inputs and reject
  retargeting while busy.
- Fixed sibling helper path, JSON stdin and no shell command interpolation.
  Helper validates kinds, directories, names and Office availability again.
- Final-leaf exclusive/no-follow creation protects existing files and symlinks;
  partial writes and EINTR are handled, collisions are bounded.
- Preferences use private temporary files and rename rather than truncating
  an existing final-leaf symlink target.
- CI has read-only repository permission, SHA-pinned actions, no persisted
  checkout credentials, no declared signing/publication secrets or publisher.
- Packaging separates architectures and explicitly documents ad-hoc signing
  and absent notarization instead of claiming Developer ID identity.

## Confirmed Findings

None. No complete boundary violation met the evidence gate.

## Needs Validation

No source-grounded vulnerability hypothesis remained for independent validation.
Deployment and acceptance gaps below are limitations, not severity-bearing findings.
The audit executed no target-controlled code or probes, contacted no live service,
and retained no runtime artifacts. macOS hard memory caps could not be established
via finite RLIMIT_AS/DATA/RSS; the complete audit execution sandbox contract was
therefore unavailable. All ledger evidence is source-only with artifact null.
Separately authorized engineering tests/packages are not audit reproduction proof.

## Hardening (Not Vulnerabilities)

- Bound same-user helper stdin/preferences reads for robustness.
- Consider descriptor-anchored parent-directory writes and identity-checked error
  cleanup if hostile shared-directory isolation becomes a product requirement.
  Current parent directories follow normal filesystem semantics, not a sandbox.
- Use fresh build output for releases to avoid accidental stale bundle files.
- Ignore local credential/signing files and record source/toolchain/digest manifests.
- Document preview CLI overwrite behavior and check its close errors.

## Coverage

The validated ledger contains 12 units: 11 covered by source checks, 0 candidate,
0 blocked, 0 deferred, 1 out of scope. The final critic found the missing explicit
MyGo single-instance IPC/lock-lifecycle exclusion; it is now recorded out of scope.
No follow-up wave ran, in accordance with quick profile. No prior confirmations
were carried and no changed-source prior findings were revalidated.

Third-party native/FFI memory safety, dependency advisories, dependency IPC endpoint
identity/ownership, hosted GitHub policies and live release infrastructure are not
covered. Explicit source scopes do not establish Gatekeeper/TCC grants, effective
Finder registration, URL-handler selection among installed copies, live AppKit
compositing, Office/WPS opening, or Intel runtime behavior. Validate these on
controlled disposable/clean systems before calling a release production-ready.

Full machine-readable records and assignment maps are held in the external run
folder; this public document omits workstation-specific paths. Both findings and
coverage schema validators pass. Format validation is not security proof.

## Post-Review Engineering Changes

After the audited snapshot, the release adds this report and sanitized screenshots,
installation/build notes and ignore rules, plus a test-only temporary-directory
override (`NEWTOFILE_TEST_TMPDIR`) and guarded test cleanup. Runtime app/helper/Finder
sources, resource inputs, dependency pins and build/package/CI code remain identical
to the audited commit. These additions do not expand the audit's execution claim.
