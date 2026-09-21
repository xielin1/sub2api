#!/usr/bin/env bash
# 同步上游 KlN-4096/sub2api 的 klno 分支，并把自己的改动重放到新底座上。
#
# 分支职责：
#   klno  纯净镜像，永不在此提交，只做 reset --hard kin/klno
#   mine  自己的开发分支，每次同步后 rebase 到 klno 之上
#
# KlN 用 rebase 维护 klno，上游历史会被改写，所以这里必须用 reset + rebase，
# 不能用 merge —— merge 过一次后上游再改写历史就会产生重复提交。

set -euo pipefail

cd "$(dirname "$0")"

# 1. 确认工作区干净，避免 rebase 中途因未提交改动失败
if [ -n "$(git status --porcelain)" ]; then
    echo "工作区有未提交改动，先提交或 stash 后再同步："
    git status --short
    exit 1
fi

# 2. 拉取上游最新状态
echo "==> 拉取 kin 远程"
git fetch kin --tags --force

# 3. 把镜像分支对齐到上游，本地不保留任何自己的提交
echo "==> klno 对齐 kin/klno"
git checkout klno
git reset --hard kin/klno

# 4. 把自己的提交重放到新底座上
echo "==> mine rebase 到 klno"
git checkout mine
git rebase klno

# 5. 展示同步结果，确认自己的提交数量没有意外变化
echo
echo "==> 同步完成"
echo "klno: $(git log --oneline -1 klno)"
echo "mine 领先 $(git rev-list --count klno..mine) 个提交："
git log --oneline klno..mine

# 6. 提示推送命令，不自动推送，留给人工确认
echo
echo "确认无误后推送："
echo "  git push --force-with-lease origin mine"
echo "  git push --force-with-lease origin klno"
