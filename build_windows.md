# Windows 11 编译指南

本文档描述在 Windows 11 上编译 go-freerdp-webconnect 的完整步骤。

## 环境要求

| 组件 | 版本 | 说明 |
|------|------|------|
| Windows 11 | 22H2+ | 含 IoT LTSC 2024 |
| MSYS2 + MinGW-w64 | GCC 15.x | C 编译器 + CMake + make |
| Go | 1.24.1 | CGO 构建，需 GOTOOLCHAIN=go1.24.1 |
| Git | 任意 | 用于拉取源码 |

---

## 第一步：安装 MSYS2

### 推荐方式（清华镜像 sfx 包）

1. 从清华镜像下载 MSYS2 sfx 安装包：
   ```
   https://mirrors.tuna.tsinghua.edu.cn/msys2/distrib/x86_64/
   ```
   下载最新的 `msys2-x86_64-*.exe` 并安装到 `C:\DevDisk\DevTools\msys64`。

2. 打开 **MSYS2 MinGW64** 终端，安装所需工具链：
   ```bash
   pacman -S mingw-w64-x86_64-gcc \
             mingw-w64-x86_64-cmake \
             mingw-w64-x86_64-make \
             git
   ```

3. 验证安装：
   ```bash
   gcc --version    # 应显示 GCC 15.x
   cmake --version  # 应显示 CMake 4.x
   ```

> **注意**：后续所有构建操作均在 **MSYS2 MinGW64** 终端中执行。

---

## 第二步：安装 Go 工具链

1. 从官网下载 Go 1.24.1 Windows 安装包并安装（默认路径 `C:\Go`）。
2. 确保 `go.exe` 在 PATH 中可用：
   ```bash
   go version   # 应显示 go1.24.1
   ```

---

## 第三步：获取项目源码

```bash
git clone https://github.com/yourname/go-freerdp-webconnect.git
cd go-freerdp-webconnect

# 初始化 FreeRDP 子模块
git submodule update --init --recursive
```

FreeRDP 源码位于 `src/FreeRDP/`。

---

## 第四步：一键构建

在项目根目录执行构建脚本（首次运行会自动编译 FreeRDP，约需 5~15 分钟）：

```bash
./build_windows.sh
```

脚本完成后输出：
```
=== 构建成功! ===
```

### 构建产物

| 文件 | 说明 |
|------|------|
| `gofreerdp-windows.exe` | 主程序（约 20 MB，含静态 Go 运行时） |
| `install/bin/libfreerdp3.dll` | FreeRDP 核心库 |
| `install/bin/libfreerdp-client3.dll` | FreeRDP 客户端库 |
| `install/bin/libwinpr3.dll` | WinPR 运行时库 |

> **注意**：运行时三个 DLL 必须与 `.exe` 在同一目录，或其所在目录已加入 PATH。

---

## 第五步：验证构建结果

```bash
./test_windows.sh
```

全部通过时输出：
```
测试结果: ✅ 10 通过  ❌ 0 失败
=== 所有测试通过! ===
```

测试项目包括：
- 可执行文件存在性及大小
- 三个 FreeRDP DLL 存在性
- `--help` 输出正常
- `--version` 输出正常
- HTTP 服务在 56788 端口启动成功
- `GET /` 返回 HTTP 200
- `GET /api/version` 返回正确 JSON

---

## 第六步：运行服务

```bash
./run_windows.sh -h <RDP服务器地址> -u <用户名> -p <密码>
```

完整参数：

```
-h, --host     RDP 服务器地址（必填）
-P, --port     RDP 服务器端口（默认: 53389）
-u, --user     用户名（必填）
-p, --pass     密码（必填）
-l, --listen   HTTP 监听端口（默认: 54455）
```

服务启动后，浏览器访问：
```
http://localhost:54455/index-debug.html
```

---

## 构建脚本详解

### build_windows.sh 工作流程

```
[1/2] 编译 FreeRDP（仅首次执行）
  ├── cmake 配置（MinGW Makefiles）
  ├── mingw32-make -j<CPU核数>
  └── mingw32-make install → install/

[2/2] 编译 Go 项目
  ├── CGO_ENABLED=1
  ├── CC=<MSYS2>/mingw64/bin/gcc.exe
  ├── GOTOOLCHAIN=go1.24.1
  └── go build -o gofreerdp-windows.exe .
```

### CGO 编译参数（rdp.go）

```c
#cgo windows CFLAGS: -I${SRCDIR}/install/include/freerdp3 \
                     -I${SRCDIR}/install/include/winpr3 \
                     -D__STDC_NO_THREADS__=1
#cgo windows LDFLAGS: -L${SRCDIR}/install/bin \
                      -lfreerdp3 -lfreerdp-client3 -lwinpr3
```

### FreeRDP CMake 关键参数说明

| 参数 | 原因 |
|------|------|
| `-D__STDC_NO_THREADS__=1` | MinGW 缺少 C11 `threads.h`，绕过线程头文件检测 |
| `-Wno-incompatible-pointer-types` | 忽略 SSPI 回调函数指针类型不匹配警告 |
| `-DWITH_SSE2=OFF -DWITH_SIMD=OFF` | 避免 AVX intrinsics 内联汇编编译错误 |
| `-DUSE_UNWIND=OFF` | 避免依赖 Linux-only 的 `dlfcn.h` |
| `TEMP=/c/Temp TMP=/c/Temp` | 修复 GCC 在某些 Windows 路径下的临时文件权限问题 |
| `-G "MinGW Makefiles"` | 使用 MinGW 的 `mingw32-make`，而非 Visual Studio |
| `-DWITH_CLIENT=OFF -DWITH_SERVER=OFF` | 不编译 FreeRDP 命令行客户端和服务端，只编译库 |

---

## 常见问题

### Q: `gcc.exe` 找不到

确保 MSYS2 已安装 `mingw-w64-x86_64-gcc`，且脚本中路径与实际安装路径一致：
```bash
# 默认路径
MSYS64="/c/DevDisk/DevTools/msys64"
```
如果安装到其他路径，修改 `build_windows.sh` 第 15 行的 `MSYS64` 变量。

### Q: FreeRDP 编译出现 `threads.h: No such file`

已通过 `-D__STDC_NO_THREADS__=1` 解决。如遇此错误，确认 cmake 命令中包含该参数。

### Q: `go build` 失败，提示找不到 FreeRDP 库

确认 `install/bin/` 目录下存在三个 DLL 文件。如不存在，删除 `install/` 目录后重新运行 `build_windows.sh`。

### Q: 运行时提示 DLL 缺失

`run_windows.sh` 脚本已自动将 `install/bin/` 和 MinGW 的 `bin/` 加入 PATH。若直接运行 `.exe`，需手动将上述路径加入系统或用户 PATH，或将 DLL 复制到 `.exe` 同目录。

### Q: Go 模块下载失败

设置代理后重试（`GOPROXY=off` 时需确保 `go.sum` 已存在）：
```bash
export GOPROXY=https://goproxy.io
go mod download
```

---

## 目录结构

```
go-freerdp-webconnect/
├── src/FreeRDP/          # FreeRDP 源码（git submodule）
├── build/freerdp-windows/ # CMake 构建中间文件（gitignore）
├── install/              # FreeRDP 编译安装目录
│   ├── bin/              # DLL 文件
│   ├── include/          # 头文件
│   └── lib/              # 静态库（可选）
├── rdp.go                # CGO 核心代码
├── build_windows.sh      # 构建脚本
├── run_windows.sh        # 运行脚本
├── test_windows.sh       # 测试脚本
└── gofreerdp-windows.exe # 构建产物
```
