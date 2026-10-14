# NewToFile · 随处新建文件

一个使用 [MyGo Native UI](https://mygo.egoist.dev/docs/ui) 编写的原生 macOS 小工具。主界面为 Go / MyGo 0.4.8，使用 Metal 和 Core Text 绘制，不使用 HTML、JavaScript 或 WebView；Finder 扩展与安全文件创建核心保留 Swift / AppKit 实现。

在 Finder 文件夹或桌面的右键菜单里选择「新建文件」，输入名称后创建文件；也可以从菜单栏入口选择任意文件夹。

[English](README.en.md) · [下载预览版](https://github.com/logdns/newtofile/releases) · [更新记录](CHANGELOG.md)

> v0.2.2 公开包使用 NewToFile 独立证书自签名，仍为未公证预览版。请阅读安装限制，勿关闭系统安全保护。

## 原生界面

![浅色界面原生渲染预览](docs/screenshots/light.png)

[深色预览](docs/screenshots/dark.png) · [修复后的浅色设置页](docs/screenshots/settings-light.png)

以上由 MyGo 原生渲染器生成，路径已脱敏，不包含 AppKit 窗口边框。

- 隐藏式标题栏保留系统红黄绿窗口按钮；可拖动标题栏和侧栏上方移动窗口。
- 跟随系统时使用 Vibrancy 半透明侧栏；手动指定浅色/深色时使用匹配外观的实色侧栏，避免材质与文字混色。按「全部格式 / 文本与代码 / 办公文档」分类。
- 桌面、文稿、下载快捷入口，使用系统原生文件夹选择面板。
- 文件格式图标、实时文件名预览、创建状态与错误提示。
- 跟随系统浅色 / 深色外观，也可以在设置中单独选择外观。
- 独立设置页包含 Finder 扩展引导、办公应用检测、创建后打开选项。
- 支持中文输入法、键盘焦点、复制粘贴和 VoiceOver 语义。`Command + N` 新建窗口入口，`Command + Return` 创建文件，`Command + ,` 打开设置。

## 功能

- 支持 TXT、Markdown、RTF、JSON、CSV、HTML。
- 检测到 Microsoft Word、Excel、PowerPoint 时，分别显示 DOCX、XLSX、PPTX；检测到 WPS 时显示全部三种办公格式。
- 办公文件使用真实 OOXML 模板：一页空白 Word 文档、一个空白 Excel 工作表、一张 16:9 空白幻灯片。不是只有后缀的空文件。
- 同名文件自动追加 ` 2`、` 3` 等序号，不覆盖原文件或符号链接；可选创建后直接打开，否则在 Finder 中选中新文件。
- 右键单个文件夹时在该文件夹内创建；右键文件时在其所在文件夹创建；空白处右键使用当前文件夹。
- 外接磁盘挂载后自动更新 Finder 扩展监听范围。只读磁盘、系统保护位置和没有写权限的目录不能创建。
- 所有操作在本机完成，不联网，不上传文件名、内容或路径。不需要授予辅助功能或全磁盘访问权限。

## 环境

- macOS 13 Ventura 或更高版本。
- 构建需要 Go 1.27.1 或更高版本、Xcode / 对应 SDK 和 Python 3。可先运行 `go version` 和 `xcode-select -p` 检查工具链。
- `go.mod` / `go.sum` 固定 MyGo 0.4.8 及依赖；首次构建需要下载 Go 模块，安装后的 App 运行不需要联网。
- `Package.swift` 用于 Swift 核心库和测试；App、Swift 文件服务与 Finder 扩展由构建脚本直接编译、打包，不需要 Node.js、前端构建或 XcodeGen。

## 构建与安装

```bash
./scripts/build.sh
open build
```

生成 `build/NewToFile.app`，默认编译本机 CPU 架构并使用本地 ad-hoc 签名。将它拖入「应用程序」，首次打开一次，然后：

1. macOS 15 / 26 及更新系统：在「系统设置 → 通用 → 登录项与扩展 → Finder」中启用 **NewToFile**。不同版本可能显示为「Finder 扩展」。
2. macOS 13 / 14：在「系统设置 → 隐私与安全性 → 扩展 → Finder 扩展」中启用 **NewToFile**。
3. 在桌面或 Finder 文件夹的空白处右键，选择「新建文件 → 纯文本 (.txt)…」。
4. 输入名称并点击「创建文件」，快捷键是 `Command + Return`。

App 可以关闭窗口后留在菜单栏运行；即使退出 App，已启用的 Finder 扩展也会在使用时重新打开它。不自动启用扩展、不修改 Finder 偏好设置、不自动设置开机启动。

如果扩展没有出现，确认 App 已在「应用程序」并打开过，重新开关扩展，关闭并重新打开 Finder 窗口；必要时手动重新启动 Finder。不要同时保留多个已启动的 NewToFile 副本，以免 URL 协议指向旧版本。

安装 Office / WPS 后点击格式区域的刷新按钮或「设置与扩展 → 重新检测」，并重新打开右键菜单。系统必须能够识别安装应用的 Bundle ID；压缩包内或没有启动注册过的应用可能检测不到。

偏好设置保存在 `~/Library/Application Support/NewToFile/preferences.json`。第一次升级到 MyGo 界面时，会读取旧版的保存位置和「创建后打开」偏好，不更改原偏好数据。

## 验证

```bash
go test ./...
go test -race ./...
go vet ./...
swift test
python3 scripts/test_templates.py
python3 scripts/test_app.py
./scripts/preview.sh
```

- Go 的 MyGo `ui.NewTester` 测试覆盖格式选择、分类、名称编辑、创建快捷键、位置切换、设置开关、双外观及最小窗口尺寸；另有请求解析和偏好持久化测试。
- Swift 测试覆盖文件内容、同名保护、符号链接、32 路并发创建、非法名称和 Finder → App 请求编码。
- Python 检查办公模板 ZIP 完整性、XML、内容类型及内部关系引用。此检查不替代在 Office / WPS 中实际打开文件的兼容性验收。
- `test_app.py` 在临时目录测试已打包的 Swift 文件服务，验证真实可用格式的创建、中文路径、同名保护和错误回执。需要先构建 App。
- 预览脚本通过 MyGo 原生渲染器输出 `build/preview.png`、`build/preview-dark.png`、`build/preview-settings.png`；不需要屏幕录制权限，不创建用户文件、不读写用户偏好。另输出 `preview-forced-light.png`、`preview-forced-dark.png`，模拟系统与 App 外观相反的情况。静态预览不等于真实 AppKit 材质验收。

### 手动验收

- 桌面空白处、普通文件夹、中文/空格/`&` 路径、外接可写磁盘分别新建 TXT。
- 对单个文件夹、单个文件及多选同目录文件右键，确认目标位置正确。
- 连续创建同名文件，确认原文件内容不变。
- 未安装 Office / WPS 时不显示办公格式；安装对应应用后刷新并确认菜单变化。
- 在 Word / Excel / PowerPoint 或 WPS 中实际打开新文件，确认无修复提示且能保存。
- 在只读目录创建，确认显示错误、不留下半成品。
- 测试 App 已退出时从 Finder 发起创建，确认能打开 App 并保留目标路径。

## 签名与分发

GitHub Releases 提供 Apple Silicon / Intel 分开的 ZIP 和 DMG 预览包，以及 `SHA256SUMS.txt`。v0.2.2 产物为独立证书自签名（v0.2.1 为 ad-hoc），**尚未通过 Apple 公证**，不是生产就绪的正式安装包。Gatekeeper 可能阻止启动，Finder 扩展在部分系统上可能要求有效的开发签名。Intel 包为交叉构建，尚需 Intel 实机验收。

下载与电脑架构匹配的包，将 App 拖入「应用程序」。可先运行 `shasum -a 256 -c SHA256SUMS.txt` 核对同目录下载的全部资产。若系统提示无法验证开发者，优先从源码构建；仅对确认来源的 App 使用系统设置提供的「仍要打开」操作。不同 macOS 版本对未公证 App 的限制不同，此操作不保证可用。**不要全局禁用 Gatekeeper，也不要对整个下载目录递归移除隔离标记。**

有 Apple Developer 证书时可以使用：

```bash
SIGN_IDENTITY="Developer ID Application: Your Name (TEAMID)" ./scripts/build.sh
```

公开分发前需要为 App 和扩展使用一致的有效签名，提交 Apple 公证、装订票据，并在干净的 macOS 用户环境中验证安装、扩展启用和办公兼容性。运行 `./scripts/package.sh` 可在独立架构目录构建两种版本，并输出 ZIP、DMG 和校验和。单架构可使用 `ARCH=x86_64 OUTPUT_DIR=build/x86_64 ./scripts/build.sh`。详见 [发布流程](docs/releasing.md)。

## 结构与权限

```text
cmd/newtofile/              MyGo App 生命周期、菜单栏、窗口和系统对话框
internal/appui/             MyGo 原生界面、状态及 UI 测试
internal/filebridge/        Go → Swift 文件服务与 Finder URL 解析
internal/preferences/      本地偏好持久化
Sources/FileHelper/        Swift JSON 文件服务，复用原文件创建核心
Sources/FinderExtension/    Finder 右键菜单及目标文件夹解析
Sources/NewToFileCore/      文件类型、应用检测、文件创建及 URL 请求
Resources/Templates/       可重建的空白办公模板
Resources/                 Bundle 配置和扩展沙盒权限
Tests/                     Swift 单元测试
scripts/                   构建、模板生成、校验和原生 UI 预览
```

Finder 扩展在沙盒中运行，只负责菜单和路径；通过 `newtofile://create` 唤起 MyGo 主 App。主 App 使用单实例机制，所有 URL 请求只会预填可见表单，必须点击创建，不会仅凭外部 URL 自动写文件。

主 App 和随包签名的 `NewToFileBridge` 文件服务不启用 App Sandbox，以便写入用户确认的目标位置。Go 使用标准输入 / 输出 JSON 调用本机 Swift 文件服务，不执行拼接 shell 命令；文件服务继续使用已测试的独占创建逻辑和真实模板。macOS 的桌面、文稿、下载、移动卷访问授权和其他文件系统权限仍然有效，权限被拒绝时会明确报错。未做 Mac App Store 分发适配。

当前不支持旧二进制格式 `.doc`、`.xls`、`.ppt`，也不支持 iCloud 尚未下载目录、所有第三方虚拟文件系统及搜索结果视图。Finder 扩展监听的是文件夹，不是屏幕坐标；文件的桌面图标摆放位置由 Finder 决定。

## 开源与反馈

项目源码采用 [MIT 许可证](LICENSE)。第三方组件保留各自许可证，见 `Resources/ThirdParty/`，这些声明随 App 分发。

- 问题与功能建议：[GitHub Issues](https://github.com/logdns/newtofile/issues)。
- 参与开发：[CONTRIBUTING.md](CONTRIBUTING.md)。
- 漏洞反馈与权限边界：[SECURITY.md](SECURITY.md)。
- 本次审计范围、结果与限制：[安全审计记录](docs/security-audit.md)。
- 卸载：先禁用 Finder 扩展、退出菜单栏 App，再移除应用；如不需要保留偏好，可删除 `~/Library/Application Support/NewToFile/`。不要删除你创建的文档。

[Version 0.2.1 validation record](docs/validation.md)

[自签名说明与证书指纹](docs/self-signing.md) · [v0.2.2 验证记录](docs/validation-0.2.2.md)
