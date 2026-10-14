import AppKit

let directory = CommandLine.arguments[1]
try FileManager.default.createDirectory(atPath: directory, withIntermediateDirectories: true)
for size in [16, 32, 128, 256, 512] {
    for scale in [1, 2] {
        let pixels = size * scale
        let bitmap = NSBitmapImageRep(bitmapDataPlanes: nil, pixelsWide: pixels, pixelsHigh: pixels,
                                      bitsPerSample: 8, samplesPerPixel: 4, hasAlpha: true,
                                      isPlanar: false, colorSpaceName: .deviceRGB, bytesPerRow: 0, bitsPerPixel: 0)!
        let context = NSGraphicsContext(bitmapImageRep: bitmap)!
        NSGraphicsContext.saveGraphicsState()
        NSGraphicsContext.current = context
        context.cgContext.scaleBy(x: CGFloat(pixels) / 1024, y: CGFloat(pixels) / 1024)
        let background = NSBezierPath(roundedRect: NSRect(x: 64, y: 64, width: 896, height: 896), xRadius: 200, yRadius: 200)
        NSGradient(starting: NSColor(red: 0.33, green: 0.55, blue: 0.74, alpha: 1),
                   ending: NSColor(red: 0.21, green: 0.35, blue: 0.54, alpha: 1))!.draw(in: background, angle: 110)
        NSColor(red: 0.97, green: 0.98, blue: 1.0, alpha: 1).setFill()
        NSBezierPath(roundedRect: NSRect(x: 285, y: 240, width: 420, height: 555), xRadius: 45, yRadius: 45).fill()
        NSColor(red: 0.69, green: 0.77, blue: 0.88, alpha: 1).setFill()
        for y in [620, 535, 450] {
            NSBezierPath(roundedRect: NSRect(x: 355, y: y, width: 250, height: 22), xRadius: 11, yRadius: 11).fill()
        }
        NSColor(red: 0.85, green: 0.92, blue: 1.0, alpha: 1).setFill()
        NSBezierPath(ovalIn: NSRect(x: 565, y: 180, width: 250, height: 250)).fill()
        NSColor(red: 0.21, green: 0.35, blue: 0.54, alpha: 1).setFill()
        NSBezierPath(roundedRect: NSRect(x: 675, y: 230, width: 30, height: 150), xRadius: 10, yRadius: 10).fill()
        NSBezierPath(roundedRect: NSRect(x: 615, y: 290, width: 150, height: 30), xRadius: 10, yRadius: 10).fill()
        NSGraphicsContext.restoreGraphicsState()
        let suffix = scale == 2 ? "@2x" : ""
        let output = URL(fileURLWithPath: directory).appendingPathComponent("icon_\(size)x\(size)\(suffix).png")
        try bitmap.representation(using: .png, properties: [:])!.write(to: output)
    }
}
