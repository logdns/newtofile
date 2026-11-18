# v0.2.4 Scoped Security Review

Date: 2026-11-18. Source-only review of Office/WPS detection and its use by the
Finder extension and packaged file helper. This is not a full application,
dependency, or release infrastructure audit.

## Result

No security vulnerability was confirmed in the reviewed changes.

- `Sources/NewToFileCore/OfficeAvailability.swift` uses fixed bundle identifiers
  and standard application names. Its fallback reads bundle metadata from
  macOS application directories; it does not load or execute the detected app.
- A matching bundle ID is a format-availability hint, not a signature or
  publisher check. The discovered application path is not passed to the
  creation or open operation. The existing optional open-after-create behavior
  still delegates to macOS's default file handler.
- `Sources/FileHelper/main.swift` still checks format availability before
  creating an Office document. Plain-text and other non-Office formats skip
  application detection; file name and destination validation, template
  selection, and the Finder sandbox configuration are unchanged.

## Limits

Separate engineering tests cover registered identifiers, realistic bundle
metadata fixtures, format mapping, and the packaged helper with installed WPS.
They are not audit sandbox evidence. Real Microsoft Office detection/opening,
WPS opening/saving, and Finder extension launch and sandbox visibility on a
clean Mac were not verified. Earlier scope is recorded in
[the original review](security-audit.md).
