import XCTest
@testable import NewToFileCore

final class CreationRequestTests: XCTestCase {
    func testRequestRoundTripPreservesSpecialCharacters() {
        let path = "/Users/test/Desktop/新文件 & 100% # + ?"
        for kind in FileKind.allCases {
            let original = CreationRequest(kind: kind, directory: URL(fileURLWithPath: path))
            let decoded = CreationRequest(url: original.url)
            XCTAssertEqual(decoded?.directory.path, path)
            XCTAssertEqual(decoded?.kind, kind)
        }
    }

    func testMalformedRequestsAreRejected() {
        for value in [
            "https://create?type=txt&directory=/tmp",
            "newtofile://wrong?type=txt&directory=/tmp",
            "newtofile://create?type=exe&directory=/tmp",
            "newtofile://create?type=txt&directory=relative",
            "newtofile://create?type=txt",
            "newtofile://create?directory=/tmp",
            "newtofile://create?type=txt&directory=/tmp%00bad"
        ] {
            XCTAssertNil(CreationRequest(url: URL(string: value)!), value)
        }
    }
}
