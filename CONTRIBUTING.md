# Contributing

Thanks for helping improve NewToFile. Discussions and issues may be written in
Chinese or English. Please search existing issues before opening a new one.

## Development

Use macOS 13+, Go 1.27.1+, Xcode with its macOS SDK, and Python 3.
See [README](README.md) and [release instructions](docs/releasing.md).
Run the Go, Swift, template, packaged-helper, and appearance checks documented
in the README. Include a regression test with behavioral changes.

Keep the MyGo interface native; do not introduce a WebView or a frontend build.
Keep Finder URL requests as visible form prefills, never automatic writes.
Do not weaken exclusive file creation or allow existing files to be overwritten.
Use disposable directories for tests, never personal documents.

## Pull Requests

Describe the problem, proposed behavior, tests run, and remaining limitations.
For UI changes, include sanitized screenshots in both appearances, including
an explicit app appearance opposite to the system appearance.
For template changes, test ZIP/XML integrity and describe Office/WPS versions
used for manual compatibility testing; integrity alone is not compatibility.

Do not commit build output, private paths, credentials, signing keys, or
notarization profiles. Contributions are licensed under the project's MIT
license. Preserve third-party copyright and license notices.

For vulnerabilities, use [SECURITY.md](SECURITY.md), not a public issue.
