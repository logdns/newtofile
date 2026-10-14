import Foundation

public struct CreationRequest {
    public let kind: FileKind
    public let directory: URL

    public init(kind: FileKind, directory: URL) {
        self.kind = kind
        self.directory = directory
    }

    public var url: URL {
        var components = URLComponents()
        components.scheme = "newtofile"
        components.host = "create"
        components.queryItems = [
            URLQueryItem(name: "type", value: kind.rawValue),
            URLQueryItem(name: "directory", value: directory.path)
        ]
        return components.url!
    }

    public init?(url: URL) {
        guard let components = URLComponents(url: url, resolvingAgainstBaseURL: false),
              components.scheme == "newtofile", components.host == "create",
              let items = components.queryItems,
              let type = items.first(where: { $0.name == "type" })?.value,
              let kind = FileKind(rawValue: type),
              let path = items.first(where: { $0.name == "directory" })?.value,
              path.hasPrefix("/"), !path.contains("\0") else { return nil }
        self.kind = kind
        self.directory = URL(fileURLWithPath: path, isDirectory: true)
    }
}
