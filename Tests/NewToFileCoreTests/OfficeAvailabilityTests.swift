import XCTest
@testable import NewToFileCore

final class OfficeAvailabilityTests: XCTestCase {
    func testMicrosoftAppsEnableOnlyTheirOwnFormatsWhenRegistered() {
        let apps: [(String, FileKind)] = [
            ("com.microsoft.Word", .docx),
            ("com.microsoft.Excel", .xlsx),
            ("com.microsoft.Powerpoint", .pptx)
        ]

        for (identifier, kind) in apps {
            let availability = OfficeAvailability(applicationDirectories: []) { $0 == identifier }
            XCTAssertEqual(availability.availableKinds.filter(\.isOffice), [kind], identifier)
        }
    }

    func testUnregisteredMicrosoftAppsInMultipleApplicationsDirectories() throws {
        try withApplicationsRoot { root in
            let local = root.appendingPathComponent("local/Applications")
            let user = root.appendingPathComponent("user/Applications")
            try writeBundle("Microsoft Word.app", identifier: "com.microsoft.Word", in: local)
            try writeBundle("Microsoft Excel.app", identifier: "com.microsoft.Excel", in: user)
            try writeBundle("Microsoft PowerPoint.app", identifier: "com.microsoft.Powerpoint", in: local)

            let availability = OfficeAvailability(applicationDirectories: [local, user]) { _ in false }
            XCTAssertEqual(availability.availableKinds.filter(\.isOffice), [.docx, .xlsx, .pptx])
            XCTAssertEqual(availability.description, "已检测到 Word、Excel、PowerPoint")
        }
    }

    func testOfficeHubAndIncorrectBundleIDDoNotEnableFormats() throws {
        try withApplicationsRoot { root in
            try writeBundle("Microsoft Word.app", identifier: "com.example.Word", in: root)
            let availability = OfficeAvailability(applicationDirectories: [root]) { $0 == "com.microsoft.Office" }
            XCTAssertEqual(availability.availableKinds, FileKind.allCases.filter { !$0.isOffice })
        }
    }

    func testWPSStillEnablesAllOfficeFormats() throws {
        let registered = OfficeAvailability(applicationDirectories: []) { $0 == "com.kingsoft.wpsoffice.mac" }
        XCTAssertEqual(registered.availableKinds.filter(\.isOffice), [.docx, .xlsx, .pptx])

        try withApplicationsRoot { root in
            try writeBundle("WPS Office.app", identifier: "com.kingsoft.wpsoffice.mac.global", in: root)
            let unregistered = OfficeAvailability(applicationDirectories: [root]) { _ in false }
            XCTAssertEqual(unregistered.availableKinds.filter(\.isOffice), [.docx, .xlsx, .pptx])
        }
    }

    private func withApplicationsRoot(_ test: (URL) throws -> Void) throws {
        let root = FileManager.default.temporaryDirectory.appendingPathComponent(UUID().uuidString)
        try FileManager.default.createDirectory(at: root, withIntermediateDirectories: true)
        defer { try? FileManager.default.removeItem(at: root) }
        try test(root)
    }

    private func writeBundle(_ name: String, identifier: String, in directory: URL) throws {
        let contents = directory.appendingPathComponent(name).appendingPathComponent("Contents")
        try FileManager.default.createDirectory(at: contents, withIntermediateDirectories: true)
        let info = try PropertyListSerialization.data(
            fromPropertyList: ["CFBundleIdentifier": identifier, "CFBundlePackageType": "APPL"],
            format: .xml, options: 0
        )
        try info.write(to: contents.appendingPathComponent("Info.plist"))
    }
}
