// swift-tools-version: 5.9
import PackageDescription

let package = Package(
    name: "NewToFile",
    platforms: [.macOS(.v13)],
    products: [.library(name: "NewToFileCore", targets: ["NewToFileCore"])],
    targets: [
        .target(name: "NewToFileCore"),
        .testTarget(name: "NewToFileCoreTests", dependencies: ["NewToFileCore"])
    ]
)
