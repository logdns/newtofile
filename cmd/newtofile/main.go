package main

import (
	"bytes"
	"flag"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"slices"
	"time"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"

	"newtofile/internal/appui"
	"newtofile/internal/filebridge"
	"newtofile/internal/preferences"
)

type application struct {
	model           *appui.Model
	client          filebridge.Client
	window          *mygo.Window
	tray            *mygo.Tray
	preferencesPath string
}

func main() {
	preview := flag.String("preview", "", "Render the native UI to a PNG without opening a window")
	dark := flag.Bool("dark", false, "Use dark appearance for the preview")
	settings := flag.Bool("settings", false, "Show the settings page in the preview")
	appearance := flag.Int("appearance", 0, "Preview appearance: 0 system, 1 light, 2 dark")
	vibrancy := flag.Bool("vibrancy", false, "Simulate native material availability in the preview")
	flag.Parse()
	mygo.App.SetName("NewToFile")
	mygo.App.SetVersion(appui.Version)
	desktop, err := mygo.App.Path(mygo.PathDesktop)
	if err != nil {
		log.Fatal(err)
	}
	executable, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	app := &application{client: filebridge.Client{Executable: filepath.Join(filepath.Dir(executable), "NewToFileBridge")}}
	if *preview != "" {
		// Previews use only local inspection; they do not load or save user preferences.
		if *appearance < 0 || *appearance > 2 {
			log.Fatal("invalid preview appearance")
		}
		app.model = appui.New(appui.Preferences{Directory: "/Users/you/Desktop", Appearance: *appearance})
		app.detectNow()
		if *settings {
			app.model.Page = "settings"
		}
		tester := ui.NewTester(app.model.View, 820, 740)
		tester.SetScale(2)
		tester.SetDark(*dark)
		tester.SetVibrancy(*vibrancy)
		time.Sleep(150 * time.Millisecond)
		tester.Frame()
		if err := writePreview(*preview, tester); err != nil {
			log.Fatal(err)
		}
		return
	}
	if !mygo.App.RequestSingleInstanceLock() {
		return
	}
	dataDirectory, err := mygo.App.Path(mygo.PathUserData)
	if err != nil {
		log.Fatal(err)
	}
	app.preferencesPath = filepath.Join(dataDirectory, "preferences.json")
	defaults := appui.Preferences{Directory: desktop}
	if _, err := os.Stat(app.preferencesPath); os.IsNotExist(err) {
		if legacy, err := app.client.Call(filebridge.Request{Operation: "legacyPreferences"}); err == nil {
			if filepath.IsAbs(legacy.Directory) {
				defaults.Directory = legacy.Directory
			}
			defaults.OpenAfter = legacy.OpenAfter
		}
	}
	prefs, err := preferences.Load(app.preferencesPath, defaults)
	if err != nil {
		log.Printf("Using default preferences: %v", err)
	}
	app.model = appui.New(prefs)
	app.model.StandardDirectories = map[string]string{}
	for key, name := range map[string]mygo.PathName{"Desktop": mygo.PathDesktop, "Documents": mygo.PathDocuments, "Downloads": mygo.PathDownloads} {
		if path, err := mygo.App.Path(name); err == nil {
			app.model.StandardDirectories[key] = path
		}
	}
	app.detectNow()
	app.model.Actions = appui.Actions{Create: app.create, Refresh: app.refresh, ChooseDirectory: app.chooseDirectory,
		OpenSettings: app.openSettings, Reveal: mygo.Shell.ShowItemInFolder, PreferencesChanged: app.savePreferences}
	mygo.App.SetActivationPolicy(mygo.ActivationPolicyAccessory)
	mygo.App.SetMenu(mygo.NewMenu([]*mygo.MenuItem{
		{Role: mygo.RoleAppMenu},
		{Label: "文件", Submenu: []*mygo.MenuItem{
			{Label: "新建文件…", Accelerator: "Cmd+N", Click: func(*mygo.MenuItem, *mygo.Window) { app.model.Page = "create"; app.showWindow() }},
			{Label: "创建文件", Accelerator: "Cmd+Return", Click: func(*mygo.MenuItem, *mygo.Window) {
				if app.window != nil && app.model.Page == "create" {
					app.model.Create()
					app.invalidate()
				}
			}},
			{Label: "设置与扩展…", Accelerator: "Cmd+,", Click: func(*mygo.MenuItem, *mygo.Window) { app.model.Page = "settings"; app.showWindow() }},
			mygo.Separator(), {Role: mygo.RoleClose, Label: "关闭窗口"},
		}},
		{Role: mygo.RoleEditMenu, Label: "编辑"},
		{Role: mygo.RoleWindowMenu, Label: "窗口"},
	}))
	mygo.App.OnOpenURL(func(raw string) {
		request, err := filebridge.ParseURL(raw)
		if err != nil {
			log.Print(err)
			return
		}
		if app.model.Busy {
			app.model.Error = "正在创建文件，请完成后再发起新的创建请求。"
		} else {
			app.model.Accept(request)
		}
		app.showWindow()
	})
	mygo.App.OnSecondInstance(func(_ []string, _ string) { app.showWindow() })
	mygo.App.OnActivate(func(visible bool) {
		if !visible {
			app.showWindow()
		}
	})
	mygo.App.OnWindowAllClosed(func() { mygo.App.SetActivationPolicy(mygo.ActivationPolicyAccessory) })
	mygo.App.WhenReady(func() {
		app.createTray()
		app.showWindow()
	})
	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}

func (a *application) showWindow() {
	mygo.App.SetActivationPolicy(mygo.ActivationPolicyRegular)
	if a.window == nil || a.window.IsDestroyed() {
		a.window = mygo.NewWindow(mygo.WindowOptions{
			Title: "NewToFile", Width: 820, Height: 740, MinWidth: 760, MinHeight: 680,
			TitleBarStyle: mygo.TitleBarHiddenInset, TrafficLightPosition: &mygo.Point{X: 22, Y: 21},
			Vibrancy: mygo.VibrancySidebar, Content: ui.View(a.model.View),
			StateKey: "main-mygo-v1", BackgroundColor: "light-dark(#F5F5F7, #1C1C1E)",
		})
	}
	a.window.Show()
	a.window.Focus()
	a.invalidate()
}

func (a *application) invalidate() {
	if a.window != nil {
		a.window.Invalidate()
	}
}

func (a *application) detectNow() {
	response, err := a.client.Call(filebridge.Request{Operation: "detect"})
	if err != nil {
		a.model.Office = "应用检测暂不可用"
		a.model.Error = err.Error()
		return
	}
	a.model.Kinds, a.model.Office = response.Kinds, response.Office
}

func (a *application) refresh() {
	if a.model.Detecting || a.model.Busy {
		return
	}
	a.model.Detecting = true
	go func() {
		response, err := a.client.Call(filebridge.Request{Operation: "detect"})
		mygo.RunOnMain(func() {
			a.model.Detecting = false
			if err != nil {
				a.model.Error = err.Error()
			} else {
				a.model.Kinds, a.model.Office, a.model.Error = response.Kinds, response.Office, ""
				if !slices.Contains(a.model.Kinds, a.model.Kind) {
					a.model.Kind = "txt"
				}
			}
			a.invalidate()
		})
	}()
}

func (a *application) create(request filebridge.Request) {
	if a.model.Busy {
		return
	}
	a.model.Busy = true
	openAfter := a.model.OpenAfter
	go func() {
		response, err := a.client.Call(request)
		mygo.RunOnMain(func() {
			a.model.Busy = false
			if err != nil {
				a.model.Error = err.Error()
			} else if response.Path == "" {
				a.model.Error = "文件服务未返回新文件路径。"
			} else {
				a.model.LastCreated, a.model.Error = response.Path, ""
				a.savePreferences()
				if openAfter {
					if err := mygo.Shell.OpenPath(response.Path); err != nil {
						a.model.Error = "文件已创建，但无法打开：" + err.Error()
					}
				} else {
					mygo.Shell.ShowItemInFolder(response.Path)
				}
			}
			a.invalidate()
		})
	}()
}

func (a *application) chooseDirectory() {
	if a.model.Busy {
		return
	}
	paths, err := mygo.Dialog.Open(mygo.OpenDialogOptions{Parent: a.window, Title: "选择新文件的位置", ButtonLabel: "选择文件夹",
		DefaultPath: a.model.Directory, Directory: true, CreateDirectories: true})
	if err != nil {
		a.model.Error = err.Error()
	} else if len(paths) > 0 {
		a.model.Directory, a.model.Error, a.model.LastCreated = paths[0], "", ""
	}
	a.invalidate()
}

func (a *application) openSettings() {
	if err := mygo.Shell.OpenExternal("x-apple.systempreferences:com.apple.ExtensionsPreferences"); err != nil {
		a.model.Error = "无法打开系统设置，请手动搜索「Finder 扩展」。"
	}
}

func (a *application) savePreferences() {
	if err := preferences.Save(a.preferencesPath, a.model.Preferences); err != nil {
		a.model.Error = "无法保存偏好设置：" + err.Error()
	}
}

func (a *application) createTray() {
	var data bytes.Buffer
	svg := ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round"><path d="M13 3H6a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h7M13 3v6h6V9l-6-6M18 13v8M14 17h8"/></svg>`))
	image := ui.Render(func(c *ui.Context) {
		c.Root().Background(ui.Transparent)
		ui.Box(c).Fill().Center().Children(func() { ui.Icon(c, svg).FontSize(20).TextColor(ui.Hex("#000000")) })
	}, 22, 22, 2)
	if err := png.Encode(&data, image); err != nil {
		log.Print(err)
		return
	}
	tray, err := mygo.NewTray(mygo.TrayOptions{Icon: data.Bytes(), IconIsTemplate: true, ToolTip: "NewToFile · 随处新建文件",
		Menu: mygo.NewMenu([]*mygo.MenuItem{
			{Label: "新建文件…", Click: func(*mygo.MenuItem, *mygo.Window) { a.model.Page = "create"; a.showWindow() }},
			{Label: "设置与扩展…", Click: func(*mygo.MenuItem, *mygo.Window) { a.model.Page = "settings"; a.showWindow() }},
			mygo.Separator(), {Label: "退出 NewToFile", Click: func(*mygo.MenuItem, *mygo.Window) { mygo.App.Quit() }},
		})})
	if err != nil {
		log.Print(err)
		return
	}
	a.tray = tray
}

func writePreview(path string, tester *ui.Tester) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return png.Encode(file, tester.Image())
}
