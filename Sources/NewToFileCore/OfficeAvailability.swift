import AppKit

public struct OfficeAvailability {
    public let word: Bool
    public let excel: Bool
    public let powerpoint: Bool
    public let wps: Bool

    public init() {
        self.init(applicationDirectories: FileManager.default.urls(
            for: .applicationDirectory, in: [.localDomainMask, .userDomainMask]
        )) { NSWorkspace.shared.urlForApplication(withBundleIdentifier: $0) != nil }
    }

    init(applicationDirectories: [URL], isRegistered: (String) -> Bool) {
        func installed(_ identifiers: [String], bundleNames: [String]) -> Bool {
            identifiers.contains(where: isRegistered) || applicationDirectories.contains { directory in
                bundleNames.contains { name in
                    let url = directory.appendingPathComponent(name, isDirectory: true)
                    guard let identifier = Bundle(url: url)?.bundleIdentifier else { return false }
                    return identifiers.contains(identifier)
                }
            }
        }
        word = installed(["com.microsoft.Word"], bundleNames: ["Microsoft Word.app"])
        excel = installed(["com.microsoft.Excel"], bundleNames: ["Microsoft Excel.app"])
        powerpoint = installed(["com.microsoft.Powerpoint"], bundleNames: ["Microsoft PowerPoint.app"])
        wps = installed(["com.kingsoft.wpsoffice.mac", "com.kingsoft.wpsoffice.mac.global", "com.kingsoft.wpsoffice.mac.cn"],
                        bundleNames: ["wpsoffice.app", "WPS Office.app"])
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
