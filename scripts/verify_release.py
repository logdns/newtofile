#!/usr/bin/env python3
"""Check bundle signatures, certificate identity, architecture and DMG integrity."""
import argparse
import hashlib
from pathlib import Path
import plistlib
import subprocess
import tempfile


def run(*args):
    return subprocess.check_output(args, stderr=subprocess.STDOUT, timeout=120)


def certificate_digest(path):
    der = run("/usr/bin/openssl", "x509", "-in", str(path), "-outform", "DER")
    return hashlib.sha256(der).hexdigest()


def verify_identity(path, expected):
    run("/usr/bin/codesign", "--verify", "--strict", str(path))
    with tempfile.TemporaryDirectory(prefix="NewToFile-cert-") as directory:
        prefix = Path(directory) / "cert"
        run("/usr/bin/codesign", "--display", f"--extract-certificates={prefix}", str(path))
        leaf = Path(str(prefix) + "0")
        if not leaf.is_file() or hashlib.sha256(leaf.read_bytes()).hexdigest() != expected:
            raise RuntimeError(f"Signing certificate mismatch: {path}")


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--app", required=True, type=Path)
    parser.add_argument("--arch", required=True, choices=["arm64", "x86_64"])
    parser.add_argument("--certificate", required=True, type=Path)
    parser.add_argument("--version", required=True)
    parser.add_argument("--dmg", type=Path)
    args = parser.parse_args()
    app = args.app.resolve()
    extension = app / "Contents/PlugIns/NewToFileFinder.appex"
    helper = app / "Contents/MacOS/NewToFileBridge"
    binaries = [app / "Contents/MacOS/NewToFile", helper,
                extension / "Contents/MacOS/NewToFileFinder"]
    for bundle in [app, extension]:
        with (bundle / "Contents/Info.plist").open("rb") as file:
            info = plistlib.load(file)
        if info["CFBundleShortVersionString"] != args.version:
            raise RuntimeError(f"Version mismatch: {bundle}")
    digest = certificate_digest(args.certificate)
    run("/usr/bin/codesign", "--verify", "--deep", "--strict", str(app))
    for component in [app, helper, extension]:
        verify_identity(component, digest)
    for binary in binaries:
        if run("/usr/bin/lipo", "-archs", str(binary)).decode().strip() != args.arch:
            raise RuntimeError(f"Architecture mismatch: {binary}")
    if args.dmg:
        verify_identity(args.dmg.resolve(), digest)
        run("/usr/bin/hdiutil", "verify", str(args.dmg.resolve()))
    print(f"PASS: {args.version} {args.arch}, signatures, certificate SHA256={digest}")
    print("Certificate integrity only; Apple trust/notarization is not established.")


if __name__ == "__main__":
    main()
