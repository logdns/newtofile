import Cocoa
import FinderSync

@objc(FinderSync)
final class FinderSync: FIFinderSync {
    override init() {
        super.init()
        updateDirectories()
        for name in [NSWorkspace.didMountNotification, NSWorkspace.didUnmountNotification] {
            NSWorkspace.shared.notificationCenter.addObserver(self, selector: #selector(updateDirectories), name: name, object: nil)
        }
    }

    @objc private func updateDirectories() {
        var directories: Set<URL> = [URL(fileURLWithPath: "/", isDirectory: true),
                                     URL(fileURLWithPath: NSHomeDirectory(), isDirectory: true)]
        directories.formUnion(FileManager.default.urls(for: .desktopDirectory, in: .userDomainMask))
        // Explicit mount roots also cover volumes Finder treats as separate trees.
        if let volumes = FileManager.default.mountedVolumeURLs(includingResourceValuesForKeys: nil) {
            directories.formUnion(volumes)
        }
        FIFinderSyncController.default().directoryURLs = directories
    }

    override var toolbarItemName: String { "新建文件" }
    override var toolbarItemToolTip: String { "使用 NewToFile 新建文件" }
    override var toolbarItemImage: NSImage {
        NSImage(systemSymbolName: "doc.badge.plus", accessibilityDescription: "新建文件")!
    }

    override func menu(for menuKind: FIMenuKind) -> NSMenu? {
        if menuKind == .toolbarItemMenu {
            let menu = NSMenu()
            let item = NSMenuItem(title: "打开 NewToFile…", action: #selector(openApp), keyEquivalent: "")
            item.target = self
            menu.addItem(item)
            return menu
        }
        guard let directory = targetDirectory(for: menuKind) else { return nil }
        let menu = NSMenu()
        let submenu = NSMenu(title: "新建文件")
        for kind in OfficeAvailability().availableKinds {
            let item = NSMenuItem(title: "\(kind.title) (.\(kind.rawValue))…", action: #selector(createFile(_:)), keyEquivalent: "")
            item.target = self
            // Capture the destination when the menu opens, not after Finder's selection changes.
            item.representedObject = CreationRequest(kind: kind, directory: directory).url
            item.image = NSImage(systemSymbolName: kind.symbol, accessibilityDescription: nil)
            submenu.addItem(item)
        }
        let root = NSMenuItem(title: "新建文件", action: nil, keyEquivalent: "")
        root.image = NSImage(systemSymbolName: "doc.badge.plus", accessibilityDescription: nil)
        root.submenu = submenu
        menu.addItem(root)
        return menu
    }

    private func targetDirectory(for menuKind: FIMenuKind) -> URL? {
        let controller = FIFinderSyncController.default()
        let selected = controller.selectedItemURLs() ?? []
        if menuKind == .contextualMenuForItems || menuKind == .contextualMenuForSidebar {
            if selected.count == 1, let first = selected.first {
                var isDirectory: ObjCBool = false
                if FileManager.default.fileExists(atPath: first.path, isDirectory: &isDirectory), isDirectory.boolValue {
                    return first
                }
                return first.deletingLastPathComponent()
            }
            if !selected.isEmpty {
                let parents = Set(selected.map { $0.deletingLastPathComponent() })
                guard parents.count == 1 else { return nil }
                return parents.first
            }
        }
        return controller.targetedURL()
    }

    @objc private func createFile(_ sender: NSMenuItem) {
        guard let url = sender.representedObject as? URL else { return }
        open(url)
    }

    @objc private func openApp() { open(URL(string: "newtofile://settings")!) }

    private func open(_ url: URL) {
        let configuration = NSWorkspace.OpenConfiguration()
        configuration.activates = true
        NSWorkspace.shared.open(url, configuration: configuration) { _, error in
            if let error { NSLog("NewToFile could not open host: %@", error.localizedDescription) }
        }
    }
}
