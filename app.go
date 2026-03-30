package main

// Wails 绑定对象
// 提供给前端调用的 Go 方法

import (
	"context"
	"fmt"
	"sync/atomic"
)

// App 是 Wails 绑定对象，暴露给前端 JS 调用
type App struct {
	ctx    context.Context
	wsPort int32 // WebSocket 桥接端口（atomic）
}

// NewApp 创建 App 实例
func NewApp() *App {
	return &App{}
}

// startup 在 Wails 应用启动时调用
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// 启动本地 WebSocket 桥接服务（随机端口）
	port, err := StartWSBridge("127.0.0.1:0")
	if err != nil {
		fmt.Println("Failed to start WebSocket bridge:", err)
		return
	}
	atomic.StoreInt32(&a.wsPort, int32(port))
}

// shutdown 在 Wails 应用关闭时调用
func (a *App) shutdown(ctx context.Context) {
	fmt.Println("Application shutting down")
}

// Connect 由前端调用，注册 RDP 连接参数并返回 WebSocket 地址
func (a *App) Connect(
	host, user, pass string,
	port, width, height int,
	perf, fntlm int,
	nowallp, nowdrag, nomani, notheme, nonla, notls bool,
) string {
	settings := &rdpConnectionSettings{
		hostname: &host,
		username: &user,
		password: &pass,
		width:    width,
		height:   height,
		port:     port,
		perf:     perf,
		fntlm:    fntlm,
		nowallp:  nowallp,
		nowdrag:  nowdrag,
		nomani:   nomani,
		notheme:  notheme,
		nonla:    nonla,
		notls:    notls,
	}

	token := RegisterConnection(settings)
	wsPort := atomic.LoadInt32(&a.wsPort)

	return fmt.Sprintf("ws://127.0.0.1:%d/ws?token=%s&dtsize=%dx%d", wsPort, token, width, height)
}

// GetVersion 返回应用版本和 FreeRDP 版本
func (a *App) GetVersion() map[string]string {
	return map[string]string{
		"app":     appVer,
		"freerdp": GetFreeRDPVersion(),
	}
}
