import XCTest
@testable import NewToFileCore

final class FileCreatorTests: XCTestCase {
    private var directory: URL!
    private var templates: URL!
    private var creator: FileCreator!

    override func setUpWithError() throws {
        directory = FileManager.default.temporaryDirectory.appendingPathComponent(UUID().uuidString, isDirectory: true)
        try FileManager.default.createDirectory(at: directory, withIntermediateDirectories: true)
        templates = URL(fileURLWithPath: #filePath).deletingLastPathComponent().deletingLastPathComponent()
            .deletingLastPathComponent().appendingPathComponent("Resources/Templates", isDirectory: true)
        creator = FileCreator(templateDirectory: templates)
    }

    override func tearDownWithError() throws {
        try FileManager.default.removeItem(at: directory)
    }

    func testPlainTextIsEmptyAndHasCorrectExtension() throws {
        let url = try creator.create(kind: .txt, name: "笔记", directory: directory)
        XCTAssertEqual(url.lastPathComponent, "笔记.txt")
        XCTAssertEqual(try Data(contentsOf: url), Data())
    }

    func testExtensionIsNotDuplicated() throws {
        let url = try creator.create(kind: .txt, name: " Note.TXT ", directory: directory)
        XCTAssertEqual(url.lastPathComponent, "Note.txt")
    }

    func testExistingFileIsNeverOverwritten() throws {
        let original = directory.appendingPathComponent("Note.txt")
        try Data("keep me".utf8).write(to: original)
        let second = try creator.create(kind: .txt, name: "Note", directory: directory)
        let third = try creator.create(kind: .txt, name: "Note", directory: directory)
        XCTAssertEqual(second.lastPathComponent, "Note 2.txt")
        XCTAssertEqual(third.lastPathComponent, "Note 3.txt")
        XCTAssertEqual(try String(contentsOf: original, encoding: .utf8), "keep me")
    }

    func testExistingSymlinkIsNeverFollowed() throws {
        let destination = directory.appendingPathComponent("original.txt")
        try Data("safe".utf8).write(to: destination)
        try FileManager.default.createSymbolicLink(at: directory.appendingPathComponent("Note.txt"), withDestinationURL: destination)
        let created = try creator.create(kind: .txt, name: "Note", directory: directory)
        XCTAssertEqual(created.lastPathComponent, "Note 2.txt")
        XCTAssertEqual(try String(contentsOf: destination, encoding: .utf8), "safe")
    }

    func testInvalidNamesAreRejected() {
        for name in ["", " ", ".", "..", ".txt", "../escape", "a/b", "a:b", "a\0b", "a\nb"] {
            XCTAssertThrowsError(try creator.create(kind: .txt, name: name, directory: directory), name)
        }
    }

    func testMissingDirectoryIsRejected() {
        XCTAssertThrowsError(try creator.create(kind: .txt, name: "Note", directory: directory.appendingPathComponent("missing")))
    }

    func testRegularFileIsNotAcceptedAsDirectory() throws {
        let file = try creator.create(kind: .txt, name: "Note", directory: directory)
        XCTAssertThrowsError(try creator.create(kind: .txt, name: "Other", directory: file))
    }

    func testMissingTemplateDoesNotLeaveAnEmptyFile() {
        let missing = FileCreator(templateDirectory: directory)
        XCTAssertThrowsError(try missing.create(kind: .docx, name: "Note", directory: directory))
        XCTAssertFalse(FileManager.default.fileExists(atPath: directory.appendingPathComponent("Note.docx").path))
    }

    func testEachOfficeFormatCopiesRealZIPTemplate() throws {
        for kind in [FileKind.docx, .xlsx, .pptx] {
            let file = try creator.create(kind: kind, name: "Document", directory: directory)
            let data = try Data(contentsOf: file)
            XCTAssertEqual(Array(data.prefix(4)), [0x50, 0x4b, 0x03, 0x04])
            XCTAssertEqual(data, try Data(contentsOf: templates.appendingPathComponent("blank.\(kind.rawValue)")))
        }
    }

    func testJSONIsParseable() throws {
        let file = try creator.create(kind: .json, name: "Data", directory: directory)
        let object = try JSONSerialization.jsonObject(with: Data(contentsOf: file)) as? [String: Any]
        XCTAssertEqual(object?.count, 0)
    }

    func testEveryNonOfficeKindCanBeCreated() throws {
        for kind in FileKind.allCases where !kind.isOffice {
            let file = try creator.create(kind: kind, name: "Test", directory: directory)
            XCTAssertEqual(try Data(contentsOf: file), kind.initialContent)
        }
    }

    func testConcurrentCreationUsesDistinctNames() {
        let lock = NSLock()
        var created: [URL] = []
        var failures: [Error] = []
        DispatchQueue.concurrentPerform(iterations: 32) { _ in
            do {
                let url = try creator.create(kind: .txt, name: "Shared", directory: directory)
                lock.lock()
                created.append(url)
                lock.unlock()
            } catch {
                lock.lock()
                failures.append(error)
                lock.unlock()
            }
        }
        XCTAssertTrue(failures.isEmpty)
        XCTAssertEqual(Set(created).count, 32)
    }

    func testHiddenFilenameIsAllowed() throws {
        let file = try creator.create(kind: .txt, name: ".notes", directory: directory)
        XCTAssertEqual(file.lastPathComponent, ".notes.txt")
    }
}
