# Screenshot Sources

Documentation update: 2026-10-14. These images illustrate installation and
appearance; they do not establish clean-machine Gatekeeper acceptance, Finder
extension activation, Intel compatibility, or Apple notarization.

| File | Source and scope |
| --- | --- |
| `light.png`, `dark.png`, `settings-light.png`, `settings-dark.png` | Existing MyGo native-renderer previews with sanitized sample paths. They omit live AppKit window chrome and are not installation evidence. |
| `finder-extension-setup.png` | Screenshot supplied by the user on 2026-10-14, copied unchanged. Shows the NewToFile settings page and a system extension detail sheet with a File Provider-labelled switch under NewToFile. The switch is off. System labels do not establish implementation type: the bundled extension registers `com.apple.FinderSync`. This screenshot is not evidence of successful context-menu activation. |
| `macos-first-open-warning.png` | Apple Support's Chinese macOS Sequoia example of an unnotarized app warning. Uses Example App, not NewToFile. |
| `macos-open-anyway.png` | Apple Support's Chinese macOS Sequoia example of Privacy & Security > Open Anyway. Uses an example app, not a local acceptance test. |
| `macos-quarantine-terminal.png` | Local Terminal window capture on 2026-10-14 displaying a controlled command example. The `xattr` and `open` lines were printed, not executed. No app quarantine attributes or system trust settings were changed for this screenshot. |

## Apple Attribution

Apple screenshots remain Apple copyrighted material; the repository's MIT
license does not grant rights to these third-party screenshots. Original article:
[Safely open apps on your Mac](https://support.apple.com/zh-cn/102445).

Original image URLs:

- [Unnotarized app warning](https://cdsassets.apple.com/live/7WUAS350/images/macos/sequoia/locale/zh-cn/macos-sequoia-app-not-opened-could-not-verify-free-from-malware.png)
- [Privacy & Security > Open Anyway](https://cdsassets.apple.com/live/7WUAS350/images/macos/sequoia/locale/zh-cn/macos-sequoia-system-settings-privacy-and-security-open-app-anyway.png)

Apple's example UI may differ from other macOS releases. Consult the linked
article for current instructions. Local documentation work did not reproduce
or bypass Gatekeeper rejection and did not change Finder extension enablement.

## Documentation Checks

Checked on 2026-10-14 against base source commit `8d8bdca`:

- All 38 local links and heading anchors in the modified documentation resolve.
- The four new PNG files have valid chunk checksums and were visually inspected.
- The supplied Finder screenshot is byte-identical to the user attachment.
- The documented `xattr -dr com.apple.quarantine` command was tested only on a
  disposable temporary fixture: it removed that attribute inside the selected
  bundle, retained other attributes, and left a sibling file's quarantine intact.
- `git diff --check` passes. No application source, signatures, release tags or
  binaries changed; no new application build, full test suite or audit was run.

These checks do not demonstrate successful first launch or Finder activation
on a clean Mac. The terminal screenshot is not evidence of either outcome.
