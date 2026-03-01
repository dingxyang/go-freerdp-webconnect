#!/bin/bash

# 设置库路径
PROJECT_ROOT=$(dirname "$(readlink -f "$0")")
FREERDP_INSTALL="${PROJECT_ROOT}/install"

export LD_LIBRARY_PATH="${FREERDP_INSTALL}/lib:${FREERDP_INSTALL}/lib/x86_64-linux-gnu:${LD_LIBRARY_PATH}"

# 运行程序
exec "${PROJECT_ROOT}/go-freerdp-webconnect" "$@"
