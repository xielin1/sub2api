#!/usr/bin/env bash
# 直接从 Wei-Shaw/sub2api 的 main 分支同步，保留 mine 上已有的全部补丁。

set -euo pipefail

cd "$(dirname "$0")"

# 1. 只在 mine 的干净工作区同步，避免合并覆盖未提交的改动。
if [ "$(git branch --show-current)" != "mine" ] || [ -n "$(git status --porcelain)" ]; then
    echo "请在干净的 mine 分支运行同步脚本"
    git status --short
    exit 1
fi

# 2. 从官方上游取最新 main；不再拉取或重置 klno。
git fetch upstream main

# 3. 合并官方更新，保留已有补丁的提交历史；冲突时停下供人工处理。
git merge --no-edit upstream/main

# 4. 展示同步后的来源和当前状态；推送仍由使用者自行决定。
git log -1 --oneline upstream/main
git status -sb
