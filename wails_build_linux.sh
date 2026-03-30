#!/usr/bin/env bash
set -euo pipefail

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
FREERDP_INSTALL="${PROJECT_ROOT}/install"
BUILD_BIN="${PROJECT_ROOT}/build/bin"

check_cmd() {
    if ! command -v "$1" >/dev/null 2>&1; then
        echo "ERROR: missing command '$1'"
        exit 1
    fi
}

check_cmd go
check_cmd node
check_cmd npm
check_cmd wails

for required in libfreerdp3.so libfreerdp-client3.so libwinpr3.so; do
    if ! find "${FREERDP_INSTALL}" -type f -name "${required}*" -print -quit 2>/dev/null | grep -q .; then
        echo "ERROR: ${required} not found under ${FREERDP_INSTALL}"
        echo "Please run: ./build_linux.sh"
        exit 1
    fi
done

LIB_PATHS=()
for dir in \
    "${FREERDP_INSTALL}/lib" \
    "${FREERDP_INSTALL}/lib64" \
    "${FREERDP_INSTALL}/lib/x86_64-linux-gnu" \
    "${FREERDP_INSTALL}/lib/aarch64-linux-gnu"; do
    [[ -d "${dir}" ]] && LIB_PATHS+=("${dir}")
done

if [[ ${#LIB_PATHS[@]} -eq 0 ]]; then
    echo "ERROR: FreeRDP runtime libraries not found in ${FREERDP_INSTALL}"
    echo "Please run: ./build_linux.sh"
    exit 1
fi

CGO_LDFLAGS_EXTRA=("-Wl,-rpath,\$ORIGIN")
for dir in "${LIB_PATHS[@]}"; do
    CGO_LDFLAGS_EXTRA+=("-L${dir}" "-Wl,-rpath,${dir}")
done

export LD_LIBRARY_PATH="$(IFS=:; echo "${LIB_PATHS[*]}"):${LD_LIBRARY_PATH:-}"
export CGO_LDFLAGS="${CGO_LDFLAGS_EXTRA[*]} ${CGO_LDFLAGS:-}"

cd "${PROJECT_ROOT}"
echo "Using FreeRDP lib paths: $(IFS=:; echo "${LIB_PATHS[*]}")"
echo "Building Wails package (Linux)..."
wails build -clean "$@"

mkdir -p "${BUILD_BIN}"
while IFS= read -r -d '' lib; do
    cp -f "${lib}" "${BUILD_BIN}/"
done < <(find "${FREERDP_INSTALL}" -type f \( -name 'libfreerdp*.so*' -o -name 'libwinpr*.so*' \) -print0)

echo "Done: ${BUILD_BIN}"
