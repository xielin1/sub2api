#!/usr/bin/env bash
# 本地构建镜像的快速脚本，避免在命令行反复输入构建参数。

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

# 1. 构建前拒绝 macOS 元数据，避免 AppleDouble 文件进入迁移资源。
METADATA_FILE="$(find "${REPO_ROOT}" -path "${REPO_ROOT}/.git" -prune -o -type f \( -name '._*' -o -name '.DS_Store' \) -print -quit)"
if [[ -n "${METADATA_FILE}" ]]; then
    echo "发现 macOS 元数据文件，停止构建：${METADATA_FILE}" >&2
    exit 1
fi

# 2. 禁止 macOS 在 Docker 读取文件时额外生成 AppleDouble 元数据。
COPYFILE_DISABLE=1 docker build -t sub2api:latest \
    --build-arg GOPROXY=https://goproxy.cn,direct \
    --build-arg GOSUMDB=sum.golang.google.cn \
    -f "${REPO_ROOT}/Dockerfile" \
    "${REPO_ROOT}"
