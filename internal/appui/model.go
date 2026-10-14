package appui

import (
	"path/filepath"
	"slices"
	"strings"

	"newtofile/internal/filebridge"
)

type Preferences struct {
	Directory  string `json:"directory"`
	OpenAfter  bool   `json:"openAfter"`
	Appearance int    `json:"appearance"`
}

type Actions struct {
	Create             func(filebridge.Request)
	Refresh            func()
	ChooseDirectory    func()
	OpenSettings       func()
	Reveal             func(string)
	PreferencesChanged func()
}

const Version = "0.2.2"

type Model struct {
	Preferences
	Actions             Actions
	Kind                string
	Name                string
	Category            string
	Page                string
	Kinds               []string
	Office              string
	Error               string
	LastCreated         string
	Busy                bool
	Detecting           bool
	StandardDirectories map[string]string
}

func New(preferences Preferences) *Model {
	return &Model{Preferences: preferences, Kind: "txt", Name: "未命名", Category: "all", Page: "create",
		Kinds: slices.Clone(filebridge.Kinds[:6]), Office: "正在检测 Office / WPS…"}
}

func (m *Model) VisibleKinds() []string {
	var kinds []string
	for _, kind := range m.Kinds {
		if m.Category == "office" && !isOffice(kind) {
			continue
		}
		if m.Category == "text" && isOffice(kind) {
			continue
		}
		kinds = append(kinds, kind)
	}
	return kinds
}

func (m *Model) SelectCategory(category string) { m.Page, m.Category = "create", category }

func (m *Model) SelectKind(kind string) {
	if m.Busy || !slices.Contains(m.Kinds, kind) {
		return
	}
	m.Kind, m.LastCreated, m.Error = kind, "", ""
}

func (m *Model) Filename() string {
	name := strings.TrimSpace(m.Name)
	if strings.HasSuffix(strings.ToLower(name), "."+m.Kind) {
		name = name[:len(name)-len(m.Kind)-1]
	}
	if name == "" {
		name = "未命名"
	}
	return name + "." + m.Kind
}

func (m *Model) Create() {
	if m.Busy || strings.TrimSpace(m.Name) == "" || m.Actions.Create == nil {
		return
	}
	if !slices.Contains(m.Kinds, m.Kind) {
		m.Error = "当前格式不可用，请刷新应用检测。"
		return
	}
	m.Error, m.LastCreated = "", ""
	m.Actions.Create(filebridge.Request{Operation: "create", Kind: m.Kind, Name: m.Name, Directory: m.Directory})
}

func (m *Model) Accept(request filebridge.Request) {
	if request.Operation == "settings" {
		m.Page = "settings"
		return
	}
	m.Page, m.Category, m.Kind = "create", "all", request.Kind
	m.Directory, m.Name, m.LastCreated, m.Error = request.Directory, "未命名", "", ""
	if !slices.Contains(m.Kinds, m.Kind) {
		m.Kind = "txt"
		m.Error = "此办公格式暂不可用，请刷新 Office / WPS 检测后重试。"
	}
}

func (m *Model) LocationName() string {
	switch filepath.Base(m.Directory) {
	case "Desktop":
		return "桌面"
	case "Documents":
		return "文稿"
	case "Downloads":
		return "下载"
	default:
		return filepath.Base(m.Directory)
	}
}

func isOffice(kind string) bool { return kind == "docx" || kind == "xlsx" || kind == "pptx" }

type format struct{ title, detail, color string }

var formats = map[string]format{
	"txt":  {"纯文本", "简单、通用的文字文件", "#697586"},
	"md":   {"Markdown", "用标记记录灵感", "#55657D"},
	"rtf":  {"富文本文档", "保留字体与段落格式", "#B28152"},
	"json": {"JSON", "结构化数据文件", "#A38340"},
	"csv":  {"CSV 表格", "轻量的表格数据", "#538B78"},
	"html": {"HTML 页面", "一个新的网页起点", "#C17854"},
	"docx": {"Word 文档", "一页空白办公文档", "#377AC5"},
	"xlsx": {"Excel 工作簿", "一个空白工作表", "#388568"},
	"pptx": {"PowerPoint", "一张 16:9 空白幻灯片", "#CE7953"},
}
