package appui

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"
)

func (m *Model) View(c *ui.Context) {
	t := *c.Theme()
	if m.Appearance == 1 {
		t = *ui.LightTheme()
	}
	if m.Appearance == 2 {
		t = *ui.DarkTheme()
	}
	t.Radius, t.FontSize = 7, 13
	c.SetTheme(&t)
	if useNativeVibrancy(c.Vibrancy(), m.Appearance) {
		c.Root().Background(ui.Transparent)
	}
	ui.Row(c).Fill().AlignItems(ui.Stretch).Children(func() {
		m.sidebar(c)
		ui.Column(c).Grow(1).MinWidth(490).Background(t.Background).Children(func() {
			m.titlebar(c)
			if m.Page == "settings" {
				m.settings(c)
			} else {
				m.createPage(c)
			}
		})
	})
}

func (m *Model) sidebar(c *ui.Context) {
	t := c.Theme()
	side := ui.Column(c).Width(188).Shrink(0).Padding(0, 14, 18).Gap(5).BorderWidth(0, 1, 0, 0).BorderColor(t.Border.Alpha(0.55))
	// AppKit's material follows the system, not a MyGo-only theme override.
	// Give explicit appearances an opaque, matching surface to keep text legible.
	if !useNativeVibrancy(c.Vibrancy(), m.Appearance) {
		side.Background(t.Surface)
	}
	side.TextColor(t.Text)
	side.Children(func() {
		ui.Box(c).Height(58).DragWindow()
		ui.Row(c).Gap(9).Padding(4, 8, 24).Children(func() {
			ui.Box(c).Size(32, 32).Radius(9).Gradient(ui.Hex("#548BBD"), ui.Hex("#365A8A"), 120).Center().Children(func() {
				ui.Icon(c, documentIcon).FontSize(20).TextColor(ui.Hex("#FFFFFF"))
			})
			ui.Column(c).Gap(2).Children(func() {
				ui.Text(c, "NewToFile").FontSize(14).FontWeight(650)
				ui.Text(c, "随处新建文件").FontSize(10).TextColor(t.TextMuted)
			})
		})
		ui.Text(c, "文件格式").FontSize(11).FontWeight(600).TextColor(t.TextMuted).Padding(0, 9, 6)
		m.nav(c, "全部格式", gridIcon, m.Page == "create" && m.Category == "all", func() { m.SelectCategory("all") })
		m.nav(c, "文本与代码", documentIcon, m.Page == "create" && m.Category == "text", func() { m.SelectCategory("text") })
		m.nav(c, "办公文档", officeIcon, m.Page == "create" && m.Category == "office", func() { m.SelectCategory("office") })
		ui.Text(c, "保存位置").FontSize(11).FontWeight(600).TextColor(t.TextMuted).Padding(23, 9, 6)
		for _, place := range []struct {
			title, folder string
			svg           *ui.SVG
		}{
			{"桌面", "Desktop", desktopIcon}, {"文稿", "Documents", folderIcon}, {"下载", "Downloads", downloadIcon},
		} {
			m.nav(c, place.title, place.svg, false, func() {
				if m.Busy {
					return
				}
				if directory := m.StandardDirectories[place.folder]; directory != "" {
					m.Directory = directory
				}
				m.Page, m.Error, m.LastCreated = "create", "", ""
			})
		}
		m.nav(c, "选择文件夹…", folderIcon, false, func() {
			if !m.Busy && m.Actions.ChooseDirectory != nil {
				m.Actions.ChooseDirectory()
			}
		})
		ui.Spacer(c)
		m.nav(c, "设置与扩展", settingsIcon, m.Page == "settings", func() { m.Page = "settings" })
		ui.Row(c).Padding(10, 9, 0).Gap(5).Children(func() {
			ui.Box(c).Size(5, 5).Radius(3).Background(t.Success)
			ui.Text(c, "本机运行 · 无需联网").FontSize(10).TextColor(t.TextMuted)
		})
	})
}

func (m *Model) nav(c *ui.Context, title string, svg *ui.SVG, selected bool, action func()) {
	t := c.Theme()
	row := ui.ButtonBase(c.Key("nav-"+title)).Label(title).Selected(selected).Justify(ui.Start).TextColor(t.Text).Padding(8, 9).Gap(9).Radius(7)
	if selected {
		row.Background(t.Accent.Alpha(0.12)).TextColor(t.Accent)
	} else if row.Hovered() {
		row.Background(t.Text.Alpha(0.045))
	}
	row.Children(func() {
		ui.Icon(c, svg).FontSize(16).TextColor(t.Accent)
		ui.Text(c, title).FontSize(12).FontWeight(500).SingleLine()
	})
	row.OnClick(action)
}

func (m *Model) titlebar(c *ui.Context) {
	t := c.Theme()
	ui.Row(c).Height(58).Shrink(0).Padding(0, 26).Gap(8).DragWindow().BorderWidth(0, 0, 1, 0).BorderColor(t.Border.Alpha(0.5)).Children(func() {
		title := "新建文件"
		if m.Page == "settings" {
			title = "设置与扩展"
		}
		ui.Text(c, title).FontSize(13).FontWeight(600)
		ui.Spacer(c)
		ui.Text(c, "⌘ N").FontSize(11).TextColor(t.TextMuted).Font("monospace")
	})
}

func (m *Model) createPage(c *ui.Context) {
	t := c.Theme()
	ui.Scroll(c.Key("create-scroll")).Grow(1).Padding(20, 28, 16).Gap(16).Children(func() {
		ui.Row(c).AlignItems(ui.Start).Children(func() {
			ui.Column(c).Gap(6).Grow(1).Children(func() {
				ui.Text(c, "一个新的开始").FontSize(23).FontWeight(650)
				ui.Text(c, "选好格式，把想法留在一个新文件里。").FontSize(12).TextColor(t.TextMuted)
			})
			ui.Box(c).Size(38, 38).Radius(12).Background(t.Accent.Alpha(0.07)).Center().Children(func() {
				ui.Icon(c, plusIcon).FontSize(20).TextColor(t.Accent)
			})
		})
		ui.Column(c).Gap(10).Children(func() {
			ui.Row(c).Children(func() {
				ui.Text(c, "文件格式").FontSize(12).FontWeight(600)
				ui.Spacer(c)
				ui.Textf(c, "%d 种可用", len(m.VisibleKinds())).FontSize(11).TextColor(t.TextMuted)
			})
			ui.Grid(c).Columns(3).Gap(9).Children(func() {
				for _, kind := range m.VisibleKinds() {
					m.formatCard(c, kind)
				}
			})
			if len(m.VisibleKinds()) == 0 {
				ui.Column(c).Padding(20).Gap(8).Radius(10).Background(t.Surface).Children(func() {
					ui.Text(c, "暂未发现办公应用").FontWeight(600)
					ui.Text(c, "安装 Office 或 WPS 后，刷新检测即可新建办公文档。").FontSize(12).TextColor(t.TextMuted)
				})
			}
			ui.Row(c).Gap(5).Children(func() {
				ui.Box(c).Size(5, 5).Radius(3).Background(t.Success)
				ui.Text(c, m.Office).FontSize(10).TextColor(t.TextMuted).SingleLine().Grow(1)
				if ui.Button(c.Key("refresh"), "").Label("刷新应用检测").Tooltip("刷新 Office / WPS 检测").Disabled(m.Detecting).Padding(4).Children(func() {
					ui.Icon(c, refreshIcon).FontSize(13).TextColor(t.TextMuted)
				}).Clicked() && m.Actions.Refresh != nil {
					m.Actions.Refresh()
				}
			})
		})
		ui.Column(c).Gap(10).Children(func() {
			ui.Text(c, "名称与位置").FontSize(12).FontWeight(600)
			ui.Column(c).Radius(10).Border(1, t.Border.Alpha(0.7)).Background(panelColor(t)).Children(func() {
				ui.Row(c).Padding(10, 14).Gap(12).Children(func() {
					ui.Text(c, "文件名").Width(46).FontSize(12).TextColor(t.TextMuted)
					ui.TextInput(c.Key("filename"), &m.Name).Label("输入文件名").Placeholder("未命名").AutoFocus().Disabled(m.Busy).Grow(1).Height(30)
					ui.Text(c, "."+m.Kind).Font("monospace").FontSize(12).TextColor(t.TextMuted)
				})
				ui.Divider(c).Background(t.Border.Alpha(0.6)).Margin(0, 14)
				ui.Row(c).Padding(10, 14).Gap(12).Children(func() {
					ui.Text(c, "保存到").Width(46).FontSize(12).TextColor(t.TextMuted)
					ui.Icon(c, folderIcon).FontSize(17).TextColor(t.Accent)
					ui.Column(c).Grow(1).MinWidth(0).Gap(3).Children(func() {
						ui.Text(c, m.LocationName()).FontSize(12).FontWeight(500).SingleLine()
						ui.Text(c, m.Directory).FontSize(10).TextColor(t.TextMuted).SingleLine().Tooltip(m.Directory)
					})
					if ui.Button(c, "更改…").Label("更改保存位置").Disabled(m.Busy).FontSize(11).Clicked() && m.Actions.ChooseDirectory != nil {
						m.Actions.ChooseDirectory()
					}
				})
			})
		})
		ui.Row(c).Gap(12).Padding(13, 15).Radius(10).Background(t.Accent.Alpha(0.045)).Children(func() {
			ui.Icon(c, formatIcon(m.Kind)).FontSize(27).TextColor(ui.Hex(formats[m.Kind].color))
			ui.Column(c).Gap(4).Grow(1).MinWidth(0).Children(func() {
				ui.Text(c, m.Filename()).FontSize(12).FontWeight(600).SingleLine().Description("文件名预览")
				ui.Text(c, formats[m.Kind].detail).FontSize(10).TextColor(t.TextMuted)
			})
			ui.Text(c, "即将创建").FontSize(10).TextColor(t.TextMuted)
		})
		if m.Error != "" {
			ui.Text(c, m.Error).FontSize(12).TextColor(t.Danger).Selectable().Label("错误信息")
		} else if m.LastCreated != "" {
			ui.Row(c).Gap(8).Children(func() {
				ui.Icon(c, checkIcon).FontSize(15).TextColor(t.Success)
				ui.Text(c, "已创建 "+filepath.Base(m.LastCreated)).FontSize(11).SingleLine().Grow(1)
				if ui.Button(c, "在 Finder 中显示").FontSize(10).Clicked() && m.Actions.Reveal != nil {
					m.Actions.Reveal(m.LastCreated)
				}
			})
		}
	})
	ui.Column(c).Padding(15, 28, 18).Gap(12).BorderWidth(1, 0, 0, 0).BorderColor(t.Border.Alpha(0.5)).Background(panelColor(t)).Children(func() {
		ui.Row(c).Gap(10).Children(func() {
			if ui.Checkbox(c.Key("open-after"), &m.OpenAfter, "创建后打开文件").FontSize(11).Changed() {
				m.savePreferences()
			}
			ui.Spacer(c)
			label := "创建文件"
			if m.Busy {
				label = "正在创建…"
			}
			if ui.PrimaryButton(c.Key("create"), "").Label("创建文件").Disabled(m.Busy || strings.TrimSpace(m.Name) == "").Height(34).Padding(0, 15).Gap(7).Children(func() {
				ui.Icon(c, plusIcon).FontSize(15)
				ui.Text(c, label).FontSize(12).FontWeight(600)
				ui.Text(c, "⌘↵").FontSize(10).Opacity(0.7).Margin(0, 0, 0, 10)
			}).Clicked() {
				m.Create()
			}
		})
		ui.Row(c).Gap(5).Children(func() {
			ui.Icon(c, lockIcon).FontSize(11).TextColor(t.TextMuted)
			ui.Text(c, "同名自动编号，原文件始终保留。").FontSize(10).TextColor(t.TextMuted)
		})
	})
	if c.Shortcut(ui.Cmd, ui.KeyEnter) {
		m.Create()
	}
}

func (m *Model) formatCard(c *ui.Context, kind string) {
	t, f := c.Theme(), formats[kind]
	selected := m.Kind == kind
	card := ui.ButtonBase(c.Key("format-"+kind)).Label(f.title).Selected(selected).TextColor(t.Text).Disabled(m.Busy).
		Height(64).Padding(12).Gap(10).Radius(9).Border(1, t.Border.Alpha(0.65)).Background(panelColor(t))
	if selected {
		card.Border(1.5, t.Accent.Alpha(0.85)).Background(t.Accent.Alpha(0.065))
	} else if card.Hovered() {
		card.Background(t.SurfaceHover)
	}
	card.Transition(ui.ElementTransition{Colors: true, Duration: 120 * time.Millisecond})
	card.Children(func() {
		ui.Icon(c, formatIcon(kind)).FontSize(25).TextColor(ui.Hex(f.color)).Shrink(0)
		ui.Column(c).Gap(5).Grow(1).MinWidth(0).Children(func() {
			ui.Text(c, f.title).FontSize(11).FontWeight(550).SingleLine()
			ui.Text(c, "."+kind).FontSize(10).Font("monospace").TextColor(t.TextMuted)
		})
		if selected {
			ui.Icon(c, checkIcon).FontSize(12).TextColor(t.Accent).Shrink(0)
		}
	})
	card.OnClick(func() { m.SelectKind(kind) })
}

func panelColor(t *ui.Theme) ui.Color {
	if t.Dark {
		return ui.Hex("#26272B")
	}
	return ui.Hex("#FFFFFF")
}

func useNativeVibrancy(available bool, appearance int) bool {
	return available && appearance == 0
}

func (m *Model) savePreferences() {
	if m.Actions.PreferencesChanged != nil {
		m.Actions.PreferencesChanged()
	}
}
