#!/usr/bin/env python3
"""Exercise the packaged Go-to-Swift file service without touching user folders."""
import json
from pathlib import Path
import subprocess
import tempfile
import unittest

ROOT = Path(__file__).resolve().parent.parent
HELPER = ROOT / "build/NewToFile.app/Contents/MacOS/NewToFileBridge"


def call(**request):
    result = subprocess.run([str(HELPER)], input=json.dumps(request), text=True,
                            capture_output=True, check=True, timeout=20)
    return json.loads(result.stdout)


class PackagedAppTests(unittest.TestCase):
    def test_detect_and_create_available_formats(self):
        response = call(operation="detect")
        self.assertNotIn("error", response)
        self.assertTrue(set(["txt", "md", "rtf", "json", "csv", "html"]).issubset(response["kinds"]))
        with tempfile.TemporaryDirectory(prefix="NewToFile-中文 & ") as directory:
            for kind in response["kinds"]:
                with self.subTest(kind=kind):
                    created = call(operation="create", kind=kind, name="测试文档", directory=directory)
                    self.assertNotIn("error", created)
                    file = Path(created["path"])
                    self.assertEqual(file.parent, Path(directory))
                    self.assertEqual(file.suffix, "." + kind)
                    if kind in ["docx", "xlsx", "pptx"]:
                        self.assertEqual(file.read_bytes(), (ROOT / f"Resources/Templates/blank.{kind}").read_bytes())

    def test_duplicate_and_invalid_name(self):
        with tempfile.TemporaryDirectory() as directory:
            first = call(operation="create", kind="txt", name="Note", directory=directory)
            Path(first["path"]).write_text("Do not overwrite")
            second = call(operation="create", kind="txt", name="Note", directory=directory)
            self.assertEqual(Path(second["path"]).name, "Note 2.txt")
            self.assertEqual(Path(first["path"]).read_text(), "Do not overwrite")
            rejected = call(operation="create", kind="txt", name="../escape", directory=directory)
            self.assertIn("error", rejected)


if __name__ == "__main__":
    unittest.main()
