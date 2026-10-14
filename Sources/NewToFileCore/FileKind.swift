import Foundation

public enum FileKind: String, CaseIterable, Identifiable {
    case txt, md, rtf, json, csv, html, docx, xlsx, pptx

    public var id: String { rawValue }

    public var title: String {
        switch self {
        case .txt: return "纯文本"
        case .md: return "Markdown"
        case .rtf: return "富文本文档"
        case .json: return "JSON"
        case .csv: return "CSV 表格"
        case .html: return "HTML 页面"
        case .docx: return "Word 文档"
        case .xlsx: return "Excel 工作簿"
        case .pptx: return "PowerPoint 演示文稿"
        }
    }

    public var symbol: String {
        switch self {
        case .txt: return "doc.text"
        case .md: return "text.alignleft"
        case .rtf: return "doc.richtext"
        case .json: return "curlybraces"
        case .csv, .xlsx: return "tablecells"
        case .html: return "chevron.left.forwardslash.chevron.right"
        case .docx: return "doc.text.fill"
        case .pptx: return "rectangle.on.rectangle"
        }
    }

    public var isOffice: Bool { [.docx, .xlsx, .pptx].contains(self) }

    public var initialContent: Data? {
        let text: String
        switch self {
        case .txt, .md, .csv: text = ""
        case .rtf: text = "{\\rtf1\\ansi\\deff0 {\\fonttbl {\\f0 Helvetica;}}\\f0\\fs24 }"
        case .json: text = "{}\n"
        case .html: text = "<!doctype html>\n<html lang=\"zh-CN\">\n<head>\n  <meta charset=\"utf-8\">\n  <meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">\n  <title>Untitled</title>\n</head>\n<body>\n</body>\n</html>\n"
        case .docx, .xlsx, .pptx: return nil
        }
        return Data(text.utf8)
    }
}
