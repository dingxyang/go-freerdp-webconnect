#!/bin/bash

# 设置库路径
PROJECT_ROOT=$(dirname "$(readlink -f "$0")")
FREERDP_INSTALL="${PROJECT_ROOT}/install"

export LD_LIBRARY_PATH="${FREERDP_INSTALL}/lib:${FREERDP_INSTALL}/lib/x86_64-linux-gnu:${LD_LIBRARY_PATH}"

# 编译
cd "${PROJECT_ROOT}"
go build -o gofreerdp . || exit 1

# 运行程序
pkill gofreerdp 2>/dev/null || true
exec "${PROJECT_ROOT}/gofreerdp" "$@"
