package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// 创建应用菜单
	appMenu := menu.NewMenu()

	// 添加设置菜单
	settingsMenu := appMenu.AddSubmenu("设置")
	settingsMenu.AddText("主题设置", keys.CmdOrCtrl("T"), func(_ *menu.CallbackData) {
		// 当菜单项被点击时，发送事件到前端
		if app.ctx != nil {
			runtime.EventsEmit(app.ctx, "open-theme-settings")
		}
	})
	settingsMenu.AddText("自定义颜色", keys.CmdOrCtrl("C"), func(_ *menu.CallbackData) {
		// 当菜单项被点击时，发送事件到前端
		if app.ctx != nil {
			runtime.EventsEmit(app.ctx, "open-color-settings")
		}
	})

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "小助手",
		Width:  600,
		Height: 400,
		//true显示边框
		Frameless: false,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		Menu:             appMenu,
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 0},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
