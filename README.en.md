# NewToFile

[中文](README.md) · [Preview Downloads](https://github.com/logdns/newtofile/releases)

A small native macOS utility for creating files from Finder or the desktop.
The Go interface uses [MyGo 0.4.8](https://mygo.egoist.dev/docs/ui), Metal and
Core Text, not a WebView. A Swift Finder Sync extension adds the context menu,
and a bundled Swift helper performs file creation.

## Features

- Create TXT, Markdown, RTF, JSON, CSV and HTML in a selected folder.
- Detect installed Word, Excel, PowerPoint or WPS and enable the corresponding
  DOCX, XLSX and PPTX formats. Office documents use actual OOXML templates.
- Native folder picker, desktop/documents/downloads shortcuts, format filters,
  filename preview, menu bar entry and preferences.
- System, light and dark appearances. Native vibrancy follows the system;
  explicit appearances use a matching opaque sidebar to avoid mixed colors.
- Duplicate names receive numeric suffixes; existing files and symlinks are
  not overwritten. Optional opening with the default application.
- No app telemetry or remote file operations. No Accessibility or Full Disk
  Access permission is required; normal macOS privacy/filesystem rules apply.

## Install

Requires macOS 13+. Choose the `arm64` package for Apple Silicon or `x86_64`
for Intel from [Releases](https://github.com/logdns/newtofile/releases).
Verify downloaded assets with the release's `SHA256SUMS.txt`.

**Current packages are ad-hoc signed, not Apple notarized, and are previews.**
Gatekeeper may block launch; extension activation varies by macOS/signing state.
Prefer building from source if needed. Only use the system's per-app Open Anyway
option for a trusted download; do not globally disable Gatekeeper. Intel builds
are cross-compiled and still require native Intel acceptance testing.

Move the app to Applications and open it once. Enable NewToFile under System
Settings > General > Login Items & Extensions > Finder (macOS 15+), or Privacy
& Security > Extensions > Finder Extensions (macOS 13/14). Right-click an empty
area of the desktop or Finder and choose New File. The interface is currently
Chinese. `Cmd+Return` creates a file; `Cmd+,` opens preferences.

If the extension is missing, toggle it, reopen Finder windows, and ensure you
do not have multiple installed/running copies. After installing Office/WPS,
refresh detection in settings and reopen the Finder menu. App bundle IDs must
be registered with macOS for detection to succeed.

An external `newtofile://` request only prefills the visible form; it does not
write a file until the user confirms creation. The main app/helper are not App
Sandbox apps, while the Finder extension is sandboxed. See [SECURITY](SECURITY.md).

## Build and Test

Requires Go 1.27.1+, Xcode/macOS SDK and Python 3 on macOS.

```sh
go mod download
go test ./...
go test -race ./...
go vet ./...
swift test
python3 scripts/test_templates.py
./scripts/build.sh
python3 scripts/test_app.py
./scripts/preview.sh
```

The app appears at `build/NewToFile.app`. Go modules are pinned; the first
dependency download needs a network connection, the installed app does not.
Run `./scripts/package.sh` for separate arm64/Intel ZIP and DMG packages.
See [release instructions](docs/releasing.md) for signing and publication.

Tests cover UI state/input/appearance, URL validation, preferences, exclusive
and concurrent creation, symlinks, helper integration and OOXML structure.
Automated integrity tests do not replace opening documents in Office/WPS or
clean-machine Finder/Gatekeeper tests. Live AppKit material should be checked
manually in both system appearances.

## Limits and Removal

No legacy DOC/XLS/PPT formats or Mac App Store packaging. Read-only/protected
folders are not writable. Not all cloud placeholders, third-party virtual
filesystems or search-result views are supported. Finder decides icon placement;
the extension targets folders, not arbitrary screen coordinates.

Preferences live in `~/Library/Application Support/NewToFile/preferences.json`.
To uninstall, disable the Finder extension, quit the menu bar app and remove
the app. Optionally remove that preferences directory; keep created documents.

## Open Source

[MIT](LICENSE), copyright logdns. Bundled dependency notices are preserved in
`Resources/ThirdParty`. See [CONTRIBUTING](CONTRIBUTING.md), [CHANGELOG](CHANGELOG.md),
[security review](docs/security-audit.md) and [SECURITY](SECURITY.md).
