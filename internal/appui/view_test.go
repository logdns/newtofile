package appui

import (
	"fmt"
	"slices"
	"testing"

	"github.com/egoist/mygo/ui"
	"newtofile/internal/filebridge"
)

func testModel() *Model {
	m := New(Preferences{Directory: "/Users/test/Desktop"})
	m.Kinds = slices.Clone(filebridge.Kinds)
	m.Office = "已检测到 WPS"
	m.StandardDirectories = map[string]string{"Desktop": "/Users/test/Desktop", "Documents": "/Users/test/Documents", "Downloads": "/Users/test/Downloads"}
	return m
}

func TestSelectFormatAndPreview(t *testing.T) {
	m := testModel()
	tt := ui.NewTester(m.View, 820, 740)
	if err := tt.Click("Markdown"); err != nil {
		t.Fatal(err)
	}
	if m.Kind != "md" || !tt.HasText("未命名.md") {
		t.Fatalf("format=%s, texts=%q", m.Kind, tt.Texts())
	}
}

func TestCreateUsesLatestTextInput(t *testing.T) {
	m := testModel()
	var requests []filebridge.Request
	m.Actions.Create = func(r filebridge.Request) { requests = append(requests, r) }
	tt := ui.NewTester(m.View, 820, 740)
	if err := tt.Click("输入文件名"); err != nil {
		t.Fatal(err)
	}
	tt.Key(ui.Cmd, ui.KeyA)
	tt.Type("旅行笔记")
	if err := tt.Click("创建文件"); err != nil {
		t.Fatal(err)
	}
	if len(requests) != 1 || requests[0].Name != "旅行笔记" || requests[0].Directory != m.Directory {
		t.Fatalf("requests=%+v", requests)
	}
}

func TestCreateKeyboardShortcut(t *testing.T) {
	m := testModel()
	count := 0
	m.Actions.Create = func(filebridge.Request) { count++ }
	tt := ui.NewTester(m.View, 820, 740)
	tt.Key(ui.Cmd, ui.KeyEnter)
	if count != 1 {
		t.Fatalf("create actions=%d", count)
	}
}

func TestDisabledCreationForBlankNameAndBusy(t *testing.T) {
	for _, busy := range []bool{false, true} {
		m := testModel()
		m.Busy = busy
		if !busy {
			m.Name = " "
		}
		count := 0
		m.Actions.Create = func(filebridge.Request) { count++ }
		tt := ui.NewTester(m.View, 820, 740)
		tt.Key(ui.Cmd, ui.KeyEnter)
		if count != 0 {
			t.Fatal("disabled form created a file")
		}
	}
}

func TestSidebarFiltersFormats(t *testing.T) {
	m := testModel()
	tt := ui.NewTester(m.View, 820, 740)
	if err := tt.Click("办公文档"); err != nil {
		t.Fatal(err)
	}
	if len(m.VisibleKinds()) != 3 || tt.HasText("Markdown") {
		t.Fatalf("office filter=%v", m.VisibleKinds())
	}
	if err := tt.Click("文本与代码"); err != nil {
		t.Fatal(err)
	}
	if len(m.VisibleKinds()) != 6 || tt.HasText("Word 文档") {
		t.Fatalf("text filter=%v", m.VisibleKinds())
	}
}

func TestMissingOfficeShowsUsefulEmptyState(t *testing.T) {
	m := New(Preferences{Directory: "/tmp"})
	m.SelectCategory("office")
	tt := ui.NewTester(m.View, 820, 740)
	if !tt.HasText("暂未发现办公应用") {
		t.Fatalf("texts=%q", tt.Texts())
	}
}

func TestSettingsAppearanceAndPreferenceToggle(t *testing.T) {
	m := testModel()
	changes := 0
	m.Actions.PreferencesChanged = func() { changes++ }
	tt := ui.NewTester(m.View, 820, 740)
	if err := tt.Click("设置与扩展"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("Finder 右键扩展") {
		t.Fatal("settings did not open")
	}
	if err := tt.Click("切换创建后打开文件"); err != nil {
		t.Fatal(err)
	}
	if !m.OpenAfter || changes != 1 {
		t.Fatalf("open=%v changes=%d", m.OpenAfter, changes)
	}
	if err := tt.Click("深色"); err != nil {
		t.Fatal(err)
	}
	if m.Appearance != 2 {
		t.Fatal("appearance did not update")
	}
}

func TestSidebarLocationChangesDestination(t *testing.T) {
	m := testModel()
	tt := ui.NewTester(m.View, 820, 740)
	if err := tt.Click("下载"); err != nil {
		t.Fatal(err)
	}
	if m.Directory != "/Users/test/Downloads" {
		t.Fatalf("directory=%q", m.Directory)
	}
}

func TestFinderRequestPrefillsWithoutWriting(t *testing.T) {
	m := testModel()
	m.Actions.Create = func(filebridge.Request) { t.Fatal("URL must not write a file") }
	m.Accept(filebridge.Request{Operation: "create", Kind: "docx", Directory: "/tmp/新文档 & notes"})
	if m.Kind != "docx" || m.Directory != "/tmp/新文档 & notes" || m.Page != "create" {
		t.Fatalf("model=%+v", m)
	}
}

func TestMinimumSizeAndBothAppearancesRender(t *testing.T) {
	for _, dark := range []bool{false, true} {
		for _, page := range []string{"create", "settings"} {
			m := testModel()
			m.Page = page
			tt := ui.NewTester(m.View, 760, 680)
			tt.SetDark(dark)
			if tt.Image().Bounds().Dx() != 760 {
				t.Fatal("incorrect image size")
			}
		}
	}
}

func TestFilenamePreviewDoesNotDuplicateExtension(t *testing.T) {
	m := testModel()
	m.Name = "笔记.TXT"
	if m.Filename() != "笔记.txt" {
		t.Fatal(m.Filename())
	}
}

func TestChineseCompositionDoesNotCreatePrematurely(t *testing.T) {
	m := testModel()
	count := 0
	m.Actions.Create = func(filebridge.Request) { count++ }
	tt := ui.NewTester(m.View, 820, 740)
	if err := tt.Click("输入文件名"); err != nil {
		t.Fatal(err)
	}
	tt.Key(ui.Cmd, ui.KeyA)
	tt.Compose("工作笔记", 4)
	tt.Key(0, ui.KeyEnter)
	if count != 0 {
		t.Fatal("IME confirmation must not create a file")
	}
	tt.Compose("", 0)
	if m.Name != "工作笔记" {
		t.Fatalf("filename=%q", m.Name)
	}
	if err := tt.Click("创建文件"); err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("create actions=%d", count)
	}
}

func TestFilenameSupportsClipboardCommands(t *testing.T) {
	m := testModel()
	tt := ui.NewTester(m.View, 820, 740)
	if err := tt.Click("输入文件名"); err != nil {
		t.Fatal(err)
	}
	tt.Command("selectAll")
	tt.Command("copy")
	tt.Type("另一份")
	tt.Command("selectAll")
	tt.Command("paste")
	if m.Name != "未命名" {
		t.Fatalf("filename=%q", m.Name)
	}
}

func TestExplicitAppearancePaintsOpaqueSidebarOverNativeMaterial(t *testing.T) {
	for _, systemDark := range []bool{false, true} {
		for _, appearance := range []int{1, 2} {
			t.Run(fmt.Sprintf("systemDark=%v/appearance=%d", systemDark, appearance), func(t *testing.T) {
				m := testModel()
				m.Appearance = appearance
				tt := ui.NewTester(m.View, 820, 740)
				tt.SetVibrancy(true)
				tt.SetDark(systemDark)
				r, g, b, a := tt.Image().At(8, 120).RGBA()
				if a != 0xffff {
					t.Fatalf("sidebar leaked system material: alpha=%d", a)
				}
				light := r+g+b > 3*0x8000
				if light != (appearance == 1) {
					t.Fatalf("wrong sidebar appearance: rgb=%d,%d,%d", r, g, b)
				}
			})
		}
	}
}

func TestSystemAppearanceKeepsNativeMaterial(t *testing.T) {
	if !useNativeVibrancy(true, 0) || useNativeVibrancy(true, 1) || useNativeVibrancy(true, 2) || useNativeVibrancy(false, 0) {
		t.Fatal("incorrect material policy")
	}
}
