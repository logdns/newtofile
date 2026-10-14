import Foundation
import Darwin

public enum CreationError: LocalizedError {
    case invalidName
    case invalidDirectory
    case missingTemplate(String)
    case cannotCreate(Int32)
    case tooManyDuplicates

    public var errorDescription: String? {
        switch self {
        case .invalidName: return "文件名不能为空，不能是 . 或 ..，也不能包含 /、: 或控制字符。"
        case .invalidDirectory: return "目标文件夹不存在或不可访问。请选择一个可写的文件夹。"
        case .missingTemplate(let ext): return "找不到 \(ext.uppercased()) 模板，请重新安装 App。"
        case .cannotCreate(let code): return "无法创建文件：\(String(cString: strerror(code)))。请检查文件夹权限。"
        case .tooManyDuplicates: return "同名文件过多，请使用其他文件名。"
        }
    }
}

public struct FileCreator {
    public let templateDirectory: URL

    public init(templateDirectory: URL) {
        self.templateDirectory = templateDirectory
    }

    public func create(kind: FileKind, name: String, directory: URL) throws -> URL {
        let stem = try Self.normalizedName(name, kind: kind)
        var isDirectory: ObjCBool = false
        guard directory.isFileURL,
              FileManager.default.fileExists(atPath: directory.path, isDirectory: &isDirectory),
              isDirectory.boolValue else { throw CreationError.invalidDirectory }

        let data: Data
        if let content = kind.initialContent {
            data = content
        } else {
            let template = templateDirectory.appendingPathComponent("blank.\(kind.rawValue)")
            guard FileManager.default.fileExists(atPath: template.path) else {
                throw CreationError.missingTemplate(kind.rawValue)
            }
            data = try Data(contentsOf: template)
        }

        for index in 0..<10_000 {
            let suffix = index == 0 ? "" : " \(index + 1)"
            let url = directory.appendingPathComponent("\(stem)\(suffix).\(kind.rawValue)")
            // Exclusive creation prevents races and never overwrites an existing file or symlink.
            let fd = open(url.path, O_WRONLY | O_CREAT | O_EXCL | O_NOFOLLOW, mode_t(0o666))
            if fd == -1 {
                if errno == EEXIST { continue }
                throw CreationError.cannotCreate(errno)
            }
            do {
                try data.withUnsafeBytes { buffer in
                    var offset = 0
                    while offset < buffer.count {
                        let count = write(fd, buffer.baseAddress!.advanced(by: offset), buffer.count - offset)
                        if count == -1 {
                            if errno == EINTR { continue }
                            throw CreationError.cannotCreate(errno)
                        }
                        guard count > 0 else { throw CreationError.cannotCreate(EIO) }
                        offset += count
                    }
                }
            } catch {
                close(fd)
                unlink(url.path)
                throw error
            }
            if close(fd) == -1 {
                let code = errno
                unlink(url.path)
                throw CreationError.cannotCreate(code)
            }
            return url
        }
        throw CreationError.tooManyDuplicates
    }

    public static func normalizedName(_ input: String, kind: FileKind) throws -> String {
        var name = input.trimmingCharacters(in: .whitespacesAndNewlines)
        guard !name.isEmpty, name != ".", name != "..",
              !name.contains("/"), !name.contains(":"),
              name.rangeOfCharacter(from: .controlCharacters) == nil else {
            throw CreationError.invalidName
        }
        let ext = ".\(kind.rawValue)"
        if name.lowercased().hasSuffix(ext) { name.removeLast(ext.count) }
        guard !name.isEmpty, name != ".", name != ".." else { throw CreationError.invalidName }
        return name
    }
}
