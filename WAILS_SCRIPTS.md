# Wails 三端脚本

## 依赖准备脚本

- macOS: `./build_macos.sh`
- Linux: `./build_linux.sh`
- Windows: `build_windows.cmd`

> Linux/Windows 会编译并安装本地 `FreeRDP` 到 `./install`，供 Wails 构建和运行使用。

## 开发启动脚本

- macOS: `./wails_dev_macos.sh`
- Linux: `./wails_dev_linux.sh`
- Windows: `wails_dev_windows.cmd`

## 打包脚本

- macOS: `./wails_build_macos.sh`
- Linux: `./wails_build_linux.sh`
- Windows: `wails_build_windows.cmd`

Linux/Windows 打包脚本会在 `wails build` 之后自动将运行所需 `FreeRDP` 动态库复制到 `build/bin`。

## 已移除旧脚本

- `run_macos.sh`
- `run_linux.sh`
- `run_windows.sh`
- `run_windows.cmd`
- `build_windows.sh`
