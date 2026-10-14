#!/usr/bin/env python3
"""Build minimal, self-contained OOXML documents with no external dependencies."""
from pathlib import Path
import zipfile

ROOT = Path(__file__).resolve().parent.parent
DEST = ROOT / "Resources" / "Templates"
REL = "http://schemas.openxmlformats.org/package/2006/relationships"
OFFICE_REL = "http://schemas.openxmlformats.org/officeDocument/2006/relationships"
W = "http://schemas.openxmlformats.org/wordprocessingml/2006/main"
S = "http://schemas.openxmlformats.org/spreadsheetml/2006/main"
P = "http://schemas.openxmlformats.org/presentationml/2006/main"
A = "http://schemas.openxmlformats.org/drawingml/2006/main"


def relationships(entries):
    return f'<Relationships xmlns="{REL}">' + "".join(
        f'<Relationship Id="{rid}" Type="{OFFICE_REL}/{kind}" Target="{target}"/>'
        for rid, kind, target in entries
    ) + "</Relationships>"


def package(extension, files, content_types, main):
    files = dict(files)
    files["[Content_Types].xml"] = (
        '<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">'
        '<Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>'
        '<Default Extension="xml" ContentType="application/xml"/>'
        + "".join(f'<Override PartName="/{path}" ContentType="{mime}"/>' for path, mime in content_types.items())
        + "</Types>"
    )
    files["_rels/.rels"] = relationships([("rId1", "officeDocument", main)])
    with zipfile.ZipFile(DEST / f"blank.{extension}", "w", zipfile.ZIP_DEFLATED) as archive:
        for path, content in sorted(files.items()):
            # Stable timestamps make rebuilds byte-for-byte reproducible.
            info = zipfile.ZipInfo(path, date_time=(2026, 1, 1, 0, 0, 0))
            info.compress_type = zipfile.ZIP_DEFLATED
            archive.writestr(info, '<?xml version="1.0" encoding="UTF-8" standalone="yes"?>\n' + content)


def word():
    package("docx", {
        "word/document.xml": f'<w:document xmlns:w="{W}"><w:body><w:p/><w:sectPr>'
        '<w:pgSz w:w="11906" w:h="16838"/>'
        '<w:pgMar w:top="1440" w:right="1440" w:bottom="1440" w:left="1440" w:header="708" w:footer="708" w:gutter="0"/>'
        '</w:sectPr></w:body></w:document>'
    }, {
        "word/document.xml": "application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"
    }, "word/document.xml")


def excel():
    package("xlsx", {
        "xl/workbook.xml": f'<workbook xmlns="{S}" xmlns:r="{OFFICE_REL}">'
        '<bookViews><workbookView activeTab="0"/></bookViews>'
        '<sheets><sheet name="Sheet1" sheetId="1" r:id="rId1"/></sheets></workbook>',
        "xl/_rels/workbook.xml.rels": relationships([("rId1", "worksheet", "worksheets/sheet1.xml")]),
        "xl/worksheets/sheet1.xml": f'<worksheet xmlns="{S}"><dimension ref="A1"/>'
        '<sheetViews><sheetView workbookViewId="0"/></sheetViews><sheetFormatPr defaultRowHeight="15"/>'
        '<sheetData/></worksheet>'
    }, {
        "xl/workbook.xml": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet.main+xml",
        "xl/worksheets/sheet1.xml": "application/vnd.openxmlformats-officedocument.spreadsheetml.worksheet+xml"
    }, "xl/workbook.xml")


def powerpoint():
    tree = '<p:spTree><p:nvGrpSpPr><p:cNvPr id="1" name=""/><p:cNvGrpSpPr/><p:nvPr/></p:nvGrpSpPr>' \
           '<p:grpSpPr><a:xfrm><a:off x="0" y="0"/><a:ext cx="0" cy="0"/>' \
           '<a:chOff x="0" y="0"/><a:chExt cx="0" cy="0"/></a:xfrm></p:grpSpPr></p:spTree>'
    namespaces = f'xmlns:p="{P}" xmlns:a="{A}" xmlns:r="{OFFICE_REL}"'
    colors = {"dk1": "000000", "lt1": "FFFFFF", "dk2": "1F2937", "lt2": "F3F4F6",
              "accent1": "176A60", "accent2": "E39C45", "accent3": "4E8098", "accent4": "79946A",
              "accent5": "BA6252", "accent6": "8B8070", "hlink": "0563C1", "folHlink": "954F72"}
    color_scheme = "".join(f'<a:{key}><a:srgbClr val="{value}"/></a:{key}>' for key, value in colors.items())
    fill = '<a:solidFill><a:schemeClr val="phClr"/></a:solidFill>'
    font = '<a:latin typeface="Arial"/><a:ea typeface=""/><a:cs typeface=""/>'
    color_map = '<p:clrMap accent1="accent1" accent2="accent2" accent3="accent3" accent4="accent4" ' \
                'accent5="accent5" accent6="accent6" bg1="lt1" bg2="lt2" folHlink="folHlink" ' \
                'hlink="hlink" tx1="dk1" tx2="dk2"/>'
    files = {
        "ppt/presentation.xml": f'<p:presentation {namespaces}>'
        '<p:sldMasterIdLst><p:sldMasterId id="2147483648" r:id="rId1"/></p:sldMasterIdLst>'
        '<p:sldIdLst><p:sldId id="256" r:id="rId2"/></p:sldIdLst>'
        '<p:sldSz cx="12192000" cy="6858000" type="screen16x9"/>'
        '<p:notesSz cx="6858000" cy="9144000"/><p:defaultTextStyle/></p:presentation>',
        "ppt/_rels/presentation.xml.rels": relationships([
            ("rId1", "slideMaster", "slideMasters/slideMaster1.xml"), ("rId2", "slide", "slides/slide1.xml")]),
        "ppt/slides/slide1.xml": f'<p:sld {namespaces}><p:cSld>{tree}</p:cSld>'
        '<p:clrMapOvr><a:masterClrMapping/></p:clrMapOvr></p:sld>',
        "ppt/slides/_rels/slide1.xml.rels": relationships([("rId1", "slideLayout", "../slideLayouts/slideLayout1.xml")]),
        "ppt/slideLayouts/slideLayout1.xml": f'<p:sldLayout {namespaces} type="blank" preserve="1">'
        f'<p:cSld name="Blank">{tree}</p:cSld><p:clrMapOvr><a:masterClrMapping/></p:clrMapOvr></p:sldLayout>',
        "ppt/slideLayouts/_rels/slideLayout1.xml.rels": relationships([("rId1", "slideMaster", "../slideMasters/slideMaster1.xml")]),
        "ppt/slideMasters/slideMaster1.xml": f'<p:sldMaster {namespaces}><p:cSld>{tree}</p:cSld>{color_map}'
        '<p:sldLayoutIdLst><p:sldLayoutId id="2147483649" r:id="rId1"/></p:sldLayoutIdLst>'
        '<p:txStyles><p:titleStyle/><p:bodyStyle/><p:otherStyle/></p:txStyles></p:sldMaster>',
        "ppt/slideMasters/_rels/slideMaster1.xml.rels": relationships([
            ("rId1", "slideLayout", "../slideLayouts/slideLayout1.xml"), ("rId2", "theme", "../theme/theme1.xml")]),
        "ppt/theme/theme1.xml": f'<a:theme xmlns:a="{A}" name="NewToFile"><a:themeElements>'
        f'<a:clrScheme name="NewToFile">{color_scheme}</a:clrScheme>'
        f'<a:fontScheme name="NewToFile"><a:majorFont>{font}</a:majorFont><a:minorFont>{font}</a:minorFont></a:fontScheme>'
        f'<a:fmtScheme name="NewToFile"><a:fillStyleLst>{fill * 3}</a:fillStyleLst><a:lnStyleLst>'
        + "".join(f'<a:ln w="{width}" cap="flat" cmpd="sng" algn="ctr">{fill}<a:prstDash val="solid"/>'
                  '<a:miter lim="800000"/></a:ln>' for width in (6350, 12700, 19050))
        + '</a:lnStyleLst><a:effectStyleLst>' + '<a:effectStyle><a:effectLst/></a:effectStyle>' * 3
        + f'</a:effectStyleLst><a:bgFillStyleLst>{fill * 3}</a:bgFillStyleLst></a:fmtScheme>'
        '</a:themeElements><a:objectDefaults/><a:extraClrSchemeLst/></a:theme>'
    }
    types = {
        "ppt/presentation.xml": "presentation.main",
        "ppt/slides/slide1.xml": "slide",
        "ppt/slideLayouts/slideLayout1.xml": "slideLayout",
        "ppt/slideMasters/slideMaster1.xml": "slideMaster"
    }
    content_types = {path: f"application/vnd.openxmlformats-officedocument.presentationml.{kind}+xml" for path, kind in types.items()}
    content_types["ppt/theme/theme1.xml"] = "application/vnd.openxmlformats-officedocument.theme+xml"
    package("pptx", files, content_types, "ppt/presentation.xml")


if __name__ == "__main__":
    DEST.mkdir(parents=True, exist_ok=True)
    word()
    excel()
    powerpoint()
    print("Generated DOCX, XLSX and PPTX templates.")
