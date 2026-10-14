import AppKit

public struct OfficeAvailability {
    public let word: Bool
    public let excel: Bool
    public let powerpoint: Bool
    public let wps: Bool

    public init() {
        let workspace = NSWorkspace.shared
        func installed(_ identifiers: [String]) -> Bool {
            identifiers.contains { workspace.urlForApplication(withBundleIdentifier: $0) != nil }
        }
        wps = installed(["com.kingsoft.wpsoffice.mac", "com.kingsoft.wpsoffice.mac.global", "com.kingsoft.wpsoffice.mac.cn"])
            || ["/Applications/wpsoffice.app", "/Applications/WPS Office.app",
                NSHomeDirectory() + "/Applications/wpsoffice.app",
                NSHomeDirectory() + "/Applications/WPS Office.app"].contains {
                    FileManager.default.fileExists(atPath: $0)
                }
        word = installed(["com.microsoft.Word"])
        excel = installed(["com.microsoft.Excel"])
        powerpoint = installed(["com.microsoft.Powerpoint"])
    }

    public var availableKinds: [FileKind] {
        FileKind.allCases.filter { kind in
            switch kind {
            case .docx: return word || wps
            case .xlsx: return excel || wps
            case .pptx: return powerpoint || wps
            default: return true
            }
        }
    }

    public var description: String {
        var names: [String] = []
        if word { names.append("Word") }
        if excel { names.append("Excel") }
        if powerpoint { names.append("PowerPoint") }
        if wps { names.append("WPS") }
        return names.isEmpty ? "未检测到 Office / WPS" : "已检测到 " + names.joined(separator: "、")
    }
}
