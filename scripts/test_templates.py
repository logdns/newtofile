#!/usr/bin/env python3
"""Check OOXML ZIP integrity, content types and internal relationship targets."""
from pathlib import Path
import posixpath
import unittest
import xml.etree.ElementTree as ET
import zipfile

ROOT = Path(__file__).resolve().parent.parent
REL = "{http://schemas.openxmlformats.org/package/2006/relationships}"
CT = "{http://schemas.openxmlformats.org/package/2006/content-types}"


class TemplateTests(unittest.TestCase):
    def test_packages(self):
        for extension, main in [("docx", "word/document.xml"), ("xlsx", "xl/workbook.xml"), ("pptx", "ppt/presentation.xml")]:
            with self.subTest(extension=extension), zipfile.ZipFile(ROOT / f"Resources/Templates/blank.{extension}") as archive:
                self.assertIsNone(archive.testzip())
                names = set(archive.namelist())
                self.assertIn(main, names)
                types = ET.fromstring(archive.read("[Content_Types].xml"))
                overrides = {node.attrib["PartName"].lstrip("/"): node.attrib["ContentType"] for node in types.findall(CT + "Override")}
                self.assertIn(main, overrides)
                for name in names:
                    root = ET.fromstring(archive.read(name))
                    if not name.endswith(".rels"):
                        continue
                    source_directory = posixpath.dirname(posixpath.dirname(name))
                    for rel in root.findall(REL + "Relationship"):
                        self.assertNotEqual(rel.attrib.get("TargetMode"), "External")
                        target = posixpath.normpath(posixpath.join(source_directory, rel.attrib["Target"]))
                        self.assertIn(target, names, f"Missing relationship target: {name} -> {target}")
                for name in overrides:
                    self.assertIn(name, names)

    def test_empty_document_structures(self):
        specs = [("docx", "word/document.xml", "{http://schemas.openxmlformats.org/wordprocessingml/2006/main}body"),
                 ("xlsx", "xl/worksheets/sheet1.xml", "{http://schemas.openxmlformats.org/spreadsheetml/2006/main}sheetData"),
                 ("pptx", "ppt/slides/slide1.xml", "{http://schemas.openxmlformats.org/presentationml/2006/main}cSld")]
        for extension, part, element in specs:
            with self.subTest(extension=extension), zipfile.ZipFile(ROOT / f"Resources/Templates/blank.{extension}") as archive:
                self.assertIsNotNone(ET.fromstring(archive.read(part)).find(element))


if __name__ == "__main__":
    unittest.main()
