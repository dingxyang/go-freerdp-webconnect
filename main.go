package main

// Wails 桌面应用入口
// 启动 Wails 窗口 + 本地 WebSocket 桥接

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed frontend/dist
var assets embed.FS

// 标准程序块
var appName string = "gofreerdp" // 应用名称
var appVer string = "0.0.2"      // 应用版本
var IsBeta string                // 是否为 Beta 版本，由构建注入
var BuildTime string             // 构建时间，由构建注入

func main() {
	app := NewApp()

	err := wails.Run(&options.App{
		Title:  "FreeRDP WebConnect",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup:  app.startup,
		OnShutdown: app.shutdown,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
