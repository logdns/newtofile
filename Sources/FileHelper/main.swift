import AppKit
import Foundation

struct HelperRequest: Decodable {
    let operation: String
    let kind: String?
    let name: String?
    let directory: String?
}

struct HelperResponse: Encodable {
    var kinds: [String]?
    var office: String?
    var path: String?
    var error: String?
    var directory: String?
    var openAfter: Bool?
}

func respond(_ response: HelperResponse) {
    if let data = try? JSONEncoder().encode(response) {
        FileHandle.standardOutput.write(data)
        FileHandle.standardOutput.write(Data([10]))
    }
}

do {
    let input = FileHandle.standardInput.readDataToEndOfFile()
    let request = try JSONDecoder().decode(HelperRequest.self, from: input)
    switch request.operation {
    case "detect":
        let office = OfficeAvailability()
        respond(HelperResponse(kinds: office.availableKinds.map(\.rawValue), office: office.description))
    case "legacyPreferences":
        let defaults = UserDefaults(suiteName: "com.newtofile.app")
        respond(HelperResponse(directory: defaults?.string(forKey: "lastDirectory"),
                               openAfter: defaults?.bool(forKey: "openAfterCreation")))
    case "create":
        guard let rawKind = request.kind, let kind = FileKind(rawValue: rawKind),
              let name = request.name, let directory = request.directory,
              directory.hasPrefix("/"), !directory.contains("\0") else {
            throw CreationError.invalidName
        }
        guard OfficeAvailability().availableKinds.contains(kind) else {
            respond(HelperResponse(error: "未检测到对应的 Office / WPS，请刷新应用检测。"))
            exit(0)
        }
        let executable = URL(fileURLWithPath: CommandLine.arguments[0]).standardizedFileURL
        let resources = executable.deletingLastPathComponent().deletingLastPathComponent()
            .appendingPathComponent("Resources/Templates", isDirectory: true)
        let file = try FileCreator(templateDirectory: resources).create(
            kind: kind, name: name, directory: URL(fileURLWithPath: directory, isDirectory: true))
        respond(HelperResponse(path: file.path))
    default:
        respond(HelperResponse(error: "不支持的文件操作。"))
    }
} catch {
    respond(HelperResponse(error: error.localizedDescription))
}
