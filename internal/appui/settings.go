package appui

import "github.com/egoist/mygo/ui"

const (
	GitHubURL  = "https://github.com/logdns/newtofile"
	WebsiteURL = "https://www.xinai.de/"
)

func (m *Model) settings(c *ui.Context) {
	t := c.Theme()
	ui.Scroll(c.Key("settings-scroll")).Grow(1).Padding(27, 28).Gap(24).Children(func() {
		ui.Column(c).Gap(6).Children(func() {
			ui.Text(c, "恰到好处的小工具").FontSize(23).FontWeight(650)
			ui.Text(c, "融入 Finder，也融入你的工作方式。").FontSize(12).TextColor(t.TextMuted)
		})
		ui.Column(c).Gap(11).Children(func() {
			ui.Text(c, "Finder 右键扩展").FontSize(12).FontWeight(600)
			ui.Column(c).Padding(18).Gap(15).Radius(11).Background(panelColor(t)).Border(1, t.Border.Alpha(0.65)).Children(func() {
				ui.Row(c).Gap(12).Children(func() {
					ui.Box(c).Size(38, 38).Radius(10).Background(t.Accent.Alpha(0.09)).Center().Children(func() { ui.Icon(c, desktopIcon).FontSize(22).TextColor(t.Accent) })
					ui.Column(c).Gap(5).Grow(1).Children(func() {
						ui.Text(c, "让新建文件，触手可及").FontSize(13).FontWeight(600)
						ui.Text(c, "在桌面或文件夹的空白处右键即可使用。").FontSize(11).TextColor(t.TextMuted)
					})
				})
				ui.Divider(c).Background(t.Border.Alpha(0.6))
				for i, step := range []string{"将 NewToFile 放入「应用程序」并打开一次。", "打开系统设置，在「Finder 扩展」中启用 NewToFile。", "回到 Finder，右键选择「新建文件」。"} {
					ui.Row(c).Gap(10).Children(func() {
						ui.Textf(c, "%d", i+1).FontSize(10).TextColor(t.TextMuted).Size(20, 20).TextAlign(ui.Center).Background(t.Surface).Radius(6)
						ui.Text(c, step).FontSize(11).TextColor(t.TextMuted).Grow(1)
					})
				}
				ui.Row(c).Justify(ui.End).Children(func() {
					if ui.Button(c, "打开系统设置").FontSize(11).Gap(5).Children(func() { ui.Icon(c, arrowIcon).FontSize(12) }).Clicked() && m.Actions.OpenSettings != nil {
						m.Actions.OpenSettings()
					}
				})
			})
		})
		ui.Column(c).Gap(11).Children(func() {
			ui.Text(c, "偏好设置").FontSize(12).FontWeight(600)
			ui.Column(c).Radius(11).Background(panelColor(t)).Border(1, t.Border.Alpha(0.65)).Children(func() {
				ui.Row(c).Padding(15).Gap(12).Children(func() {
					ui.Column(c).Gap(4).Grow(1).Children(func() {
						ui.Text(c, "创建后打开文件").FontSize(12)
						ui.Text(c, "使用系统默认应用打开新文件。").FontSize(10).TextColor(t.TextMuted)
					})
					if ui.Switch(c.Key("settings-open-after"), &m.OpenAfter).Label("切换创建后打开文件").Changed() {
						m.savePreferences()
					}
				})
				ui.Divider(c).Background(t.Border.Alpha(0.6)).Margin(0, 15)
				ui.Row(c).Padding(15).Gap(12).Children(func() {
					ui.Text(c, "外观").FontSize(12).Grow(1)
					if ui.Segmented(c.Key("appearance"), &m.Appearance, "跟随系统", "浅色", "深色").Label("界面外观").FontSize(11).Changed() {
						m.savePreferences()
					}
				})
			})
		})
		ui.Column(c).Gap(11).Children(func() {
			ui.Text(c, "办公应用").FontSize(12).FontWeight(600)
			ui.Row(c).Padding(15).Gap(10).Radius(11).Background(panelColor(t)).Border(1, t.Border.Alpha(0.65)).Children(func() {
				ui.Icon(c, officeIcon).FontSize(19).TextColor(t.TextMuted)
				ui.Text(c, m.Office).FontSize(11).Grow(1)
				if ui.Button(c, "重新检测").FontSize(11).Disabled(m.Detecting).Clicked() && m.Actions.Refresh != nil {
					m.Actions.Refresh()
				}
			})
		})
		ui.Column(c).Gap(11).Children(func() {
			ui.Text(c, "关于").FontSize(12).FontWeight(600)
			ui.Column(c).Padding(15).Gap(10).Radius(11).Background(panelColor(t)).Border(1, t.Border.Alpha(0.65)).Children(func() {
				ui.Row(c).Gap(10).Children(func() {
					ui.Box(c).Size(34, 34).Radius(9).Background(t.Accent.Alpha(0.09)).Center().Children(func() {
						ui.Icon(c, documentIcon).FontSize(19).TextColor(t.Accent)
					})
					ui.Column(c).Gap(3).Grow(1).Children(func() {
						ui.Text(c, "NewToFile").FontSize(13).FontWeight(600)
						ui.Text(c, "原生 macOS 文件创建工具 · "+Version).FontSize(10).TextColor(t.TextMuted)
					})
				})
				ui.Divider(c).Background(t.Border.Alpha(0.6))
				ui.Row(c).Gap(14).Children(func() {
					ui.Text(c, "项目链接").FontSize(11).TextColor(t.TextMuted)
					ui.Link(c, "GitHub 项目", GitHubURL).FontSize(11).TextColor(t.Accent).Tooltip("打开 GitHub 项目主页")
					ui.Link(c, "官方网站", WebsiteURL).FontSize(11).TextColor(t.Accent).Tooltip("打开 xinai.de 官方网站")
				})
			})
		})
		if m.Error != "" {
			ui.Text(c, m.Error).FontSize(11).TextColor(t.Danger).Selectable()
		}
		ui.Row(c).Gap(7).Children(func() {
			ui.Icon(c, lockIcon).FontSize(13).TextColor(t.TextMuted)
			ui.Text(c, "NewToFile "+Version+" · MyGo Native UI · 文件与路径不离开本机").FontSize(10).TextColor(t.TextMuted)
		})
	})
}
