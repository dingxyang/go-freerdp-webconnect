#!/bin/bash
# Go-FreeRDP-WebConnect Windows 构建脚本
# 依赖: MSYS2 (C:/DevDisk/DevTools/msys64) + MinGW-w64 GCC

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${GREEN}=== Go-FreeRDP-WebConnect Windows 构建脚本 ===${NC}"

PROJECT_ROOT="$(cd "$(dirname "$0")" && pwd)"
MSYS64="/c/DevDisk/DevTools/msys64"
MINGW_BIN="$MSYS64/mingw64/bin"
FREERDP_SRC="$PROJECT_ROOT/src/FreeRDP"
FREERDP_BUILD="$PROJECT_ROOT/build/freerdp-windows"
FREERDP_INSTALL="$PROJECT_ROOT/install"

# 检查 MSYS2/MinGW
if [ ! -f "$MINGW_BIN/gcc.exe" ]; then
    echo -e "${RED}错误: 未找到 GCC ($MINGW_BIN/gcc.exe)${NC}"
    echo "请先安装 MSYS2: https://mirrors.tuna.tsinghua.edu.cn/msys2/distrib/x86_64/"
    echo "或运行: pacman -S mingw-w64-x86_64-gcc mingw-w64-x86_64-cmake mingw-w64-x86_64-make"
    exit 1
fi

export PATH="$MINGW_BIN:$PATH"
echo -e "${GREEN}✓ GCC: $(gcc --version | head -1)${NC}"

# 编译 FreeRDP (如果 install 目录不存在)
if [ ! -f "$FREERDP_INSTALL/bin/libfreerdp3.dll" ]; then
    echo ""
    echo -e "${YELLOW}[1/2] 编译 FreeRDP...${NC}"

    export TEMP=/c/Temp
    export TMP=/c/Temp
    mkdir -p /c/Temp "$FREERDP_BUILD" "$FREERDP_INSTALL"

    cd "$FREERDP_BUILD"
    cmake "$FREERDP_SRC" \
        -G "MinGW Makefiles" \
        -DCMAKE_INSTALL_PREFIX="$FREERDP_INSTALL" \
        -DCMAKE_BUILD_TYPE=Release \
        -DCMAKE_C_COMPILER=gcc \
        "-DCMAKE_C_FLAGS=-D__STDC_NO_THREADS__=1 -Wno-incompatible-pointer-types" \
        -DWITH_SSE2=OFF \
        -DWITH_SIMD=OFF \
        -DWITH_CUPS=OFF \
        -DWITH_WAYLAND=OFF \
        -DWITH_PULSE=OFF \
        -DWITH_FFMPEG=OFF \
        -DWITH_SWSCALE=OFF \
        -DWITH_DSP_FFMPEG=OFF \
        -DWITH_FUSE=OFF \
        -DWITH_GSTREAMER_1_0=OFF \
        -DWITH_CLIENT=OFF \
        -DWITH_SERVER=OFF \
        -DBUILD_TESTING=OFF \
        -DCHANNEL_URBDRC=OFF \
        -DWITH_X11=OFF \
        -DWITH_ALSA=OFF \
        -DUSE_UNWIND=OFF \
        -DWITH_OPENSSL=OFF \
        2>&1

    mingw32-make -j$(nproc) 2>&1
    mingw32-make install 2>&1
    echo -e "${GREEN}✓ FreeRDP 编译完成${NC}"
else
    echo -e "${GREEN}✓ FreeRDP 已编译 (跳过)${NC}"
fi

# 编译 Go 项目
echo ""
echo -e "${YELLOW}[2/2] 编译 Go 项目...${NC}"

cd "$PROJECT_ROOT"
export CGO_ENABLED=1
export CC="$MINGW_BIN/gcc.exe"
export GOTOOLCHAIN="go1.24.1"
export GOPROXY="off"
export PATH="$FREERDP_INSTALL/bin:$PATH"

go build -o gofreerdp-windows.exe . 2>&1
echo -e "${GREEN}✓ 编译完成: gofreerdp-windows.exe${NC}"

echo ""
echo -e "${GREEN}=== 构建成功! ===${NC}"
echo ""
echo "使用方法:"
echo "  ./run_windows.sh -h <主机> -u <用户名> -p <密码>"
echo ""
echo "测试:"
echo "  ./test_windows.sh"
