# 开发记录

## 2026-09-26：版本命名统一为「版本+x序号」

- 规则：同一上游版本上每次发布递增 x 序号（如 `v0.2.8-x1`、`v0.2.8-x2`），合入上游新版本后从 `x1` 重新开始；不再使用 klno 后缀。
- git tag：本地 31 个、远程 30 个 klno tag 按原序号改名为 x 格式（如 `v0.2.7-klno.4` → `v0.2.7-x4`，`v0.2.7-klno.4-x1` → `v0.2.7-x5`，`backup-local-klno` → `backup-local-x`，仅本地），附注说明与指向提交保持不变，远程新 tag 校验一致后删除旧 tag；新增 `v0.2.7-x6` 指向 `1eef66379`（线上2曾运行的镜像）。GitHub 上 3 个 Release 已改挂新 tag 并改名。改名前的全部 tag 引用备份在本机 `/tmp/all-tags-backup-20260926.txt`。
- server-2 镜像：删除 `0.2.8-klno.2-trust`、`0.2.8-klno.1-b21c41b76` 标签（镜像本体由 `rollback-before-theme`、`rollback-before-trust` 保留），`0.2.7-x2` 改名为 `0.2.7-x6`；各发布备份目录里的 compose 已改为引用仍存在的标签，回退命令不受影响。`deploy/` 下更早的 compose 历史备份文件未改动，其中的旧镜像名已不存在。

## 2026-09-26：全站暖色主题与招募入口全局展示，更新线上2

- 提交：`ee89e8545` 首页信任区块、公开状态页与 mdai 品牌（此前未提交、已在线上的 trust 版本内容）；`917d90241` 全站配色由青色改为首页陶土暖色（Tailwind 主色/灰阶/深色背景及硬编码青色）；`d1b577607` 招募入口登录后全局展示并缩小为暖色样式。平台与分组识别色（如 DeepSeek 青色）、支付品牌色有意保留。
- `frontend/pnpm-lock.yaml` 的工作区改动（误删 overrides，会导致 frozen 安装失败）未提交，仍留在工作区。
- 发布：本地构建前端并以 `-tags embed` 交叉编译 Linux amd64 程序，打包 Dockerfile（FROM `sub2api:rollback-before-theme`）上传；`server-2` 构建 `sub2api:0.2.8-x1`（版本号改为「版本+x序号」命名，不再使用 klno 后缀），仅重建 sub2api 容器。备份及原 Compose 在 `/opt/sub2api/backups/release-d1b577607`（`pg_restore --list` 1216 条），保留 `sub2api:rollback-before-theme`。无数据库迁移，回退直接切回旧镜像即可。
- 验证：容器 healthy、重启 0、启动后 5 分钟无 ERROR/FATAL 日志；公网 `/health`、首页、登录页、`/api/v1/public/stats` 均 200，公开配置版本 `0.2.8-x1`，未登录招募接口 401，线上 CSS 含新主色 `#d66b4d`。前端全量测试除既有的 `EditAccountModal.grokMediaEligibility` 3 项外通过。未登录后台逐页检查暖色与招募入口遮挡情况。
- 部署期间本机经手机热点（中国移动）直连 server-2 的 SSH 握手频繁超时，服务器侧 sshd/防火墙/fail2ban 正常；更换网络后直连恢复。

## 2026-09-24：mdai 主域名与品牌迁移

- 仅迁移中转站及文档、CPR、vm2api、gpt2api；用户确认 git、monitor、stock 等其他系统保持原样。GoDaddy 名称服务器改为 finley/melina.ns.cloudflare.com，Cloudflare 新区域配置主站、www、docs、cpr、vm2api、gpt2api 的代理解析；新域名 HTTPS 已验证。
- 线上站点名称、API 基址、文档地址、充值通知链接及条款品牌改为 mdai；保留远日点图标。GitHub OAuth 应用改名 mdai、主页改为新域名，新增新回调并保留旧回调。注册、邮箱验证、GitHub 开启，LinuxDO 保持关闭。Brevo 域名已认证，新发件人 mdai <noreply@mdai.life> 已验证，测试邮件 Delivered。
- 旧主站已知网页 GET/HEAD 入口使用 302 保留路径和查询参数；旧 GitHub 登录起点使用 307 到新域名，避免 state Cookie 跨域。旧 API、静态文件与既有回调不做统一跳转。旧文档 302 到新文档；CPR/vm2api 控制台入口兼容跳转。新 www 和 gpt2api 网页统一回主站。
- 无应用重启、无数据库迁移。Nginx 每次配置校验后平滑重载，应用保持 2026-09-24T04:29:54.533335037Z 启动时间、restarts=0、healthy。配置备份及待发布副本位于 server-2:/opt/sub2api/backups/mdai-migration-20260924；旧证书和旧解析保留。
- 前端源码更新默认品牌和首页终端标题；线上复用 data/public/assets 覆盖 index-BGURCFT1.js、HomeView-CidyzkHj.js 的旧品牌文字，无后端重建。原资源位于上述备份目录 static/。仅定向清除新域名的对应资源缓存；部分浏览器所经 CDN 节点仍读到旧首页终端文字，需后续复核缓存传播，不能宣称所有节点已一致。
- 文档源目录 /Users/zhenxing.lin/Desktop/projects/quantix-docs 更新品牌及链接并构建，部署到 /var/www/docs.mdai.life；旧目录保留。源码备份 /tmp/quantix-docs-before-mdai-20260924。前端类型检查及改动文件 lint、文档构建通过。公网主站和文档 200、旧网页跳转路径参数保留、旧 API 未授权请求返回 401 且无跳转、新 OAuth 起点返回 GitHub 且回调 mdai.life。未完成真人 OAuth 授权闭环、未执行付费模型验收，未暂存或提交工作区。


## 2026-09-24：首页信任区块与线上2公开状态/模型广场

- 代码：新增匿名接口 `/api/v1/public/stats`（仅累计请求数与 Token 数，复用仪表盘缓存）、`/api/v1/public/announcements`（最近 5 条无定向的生效公告）；`/channel-monitors` 只读接口移出登录分组改为按 IP 限流的公开接口。`/monitor` 路由对访客公开，未登录时用模型广场导航条替代后台布局。首页新增累计数据、7 天平均可用率、模型价格（取模型广场，按分组倍率换算每百万 Token 价，最多 8 个）、痛点文案与标签、最近更新。
- 线上2配置：新建公开分组 `openai pro`(42) 并绑定 CPR-pro；`openai`(38)、`openai pro20 号池`(39) 改为公开并补描述；新建 4 个仅列模型、价格留空的渠道（计费模型来源 upstream，不改变原计费），模型广场改为免登录并填写说明；删除停用的旧监控，为 ds/openai/号池/不降智/openai pro 建 5 分钟探测监控（管理员各分组密钥经 `https://quantix.life`）；发布公告 1 条。
- CPR-pro 账号状态为 error（OAuth token 已失效），因此 `openai pro` 分组暂改回专属、其监控已停用；重新授权后需把分组改为公开并启用监控 6。`gpt-不降智` 分组当前倍率为 3，但描述写的是 0.35 倍率，模型广场按 3 倍展示价格，需站长确认。
- 发布：备份与原 Compose 在 `/opt/sub2api/backups/release-trust-20260924`（`pg_restore --list` 1216 条），保留 `sub2api:rollback-before-trust`；镜像 `sub2api:0.2.8-klno.2-trust` 由未提交工作区（基于 `c7806030d`）构建，仅重建 sub2api 容器。无数据库迁移，回退直接切回旧镜像即可。
- 验证：后端 `go build`、routes/handler/service 相关测试通过；前端类型检查、改动文件 lint、i18n 检查、views/router/components 测试通过（`EditAccountModal.grokMediaEligibility` 3 项失败为改动前既有）。线上容器 healthy、重启 0；公网匿名访问四个接口均 200，首页/状态页 1280 与 390 宽度无报错、无横向溢出。

## 2026-09-24：提交工作区并更新线上2

- 按需求分批提交：`a1e33868a` 调整自定义购买入口顺序，`b21c41b76` 增加登录后销售招募；本次同时发布此前已提交的图片工作台与官方 0.2.8 合并内容，未推送远程仓库。
- 以 `b21c41b76983911adf9e9d45f37674a4a4b5e98f` 构建 Linux amd64 程序并嵌入已验证的前端产物；在 `server-2` 复用原镜像运行依赖，发布 `sub2api:0.2.8-klno.1-b21c41b76`，仅重建 sub2api 容器。
- 更新前完成 PostgreSQL 自定义格式备份并通过 `pg_restore --list` 检查，备份及原 Compose 配置保存在 `/opt/sub2api/backups/release-b21c41b76`；保留 `sub2api:rollback-before-b21c41b76` 镜像。新迁移涉及分组计价 JSON 转换，若回退应用，需要同步核对并恢复对应数据语义，不能只切旧镜像。
- 验证通过：容器 healthy、重启次数 0、源站健康接口 200；公网首页、入口静态资源和招募配图 200，公开配置返回版本 `0.2.8-klno.1` 且不包含招募字段；未登录访问招募与图片工作台配置均为 401，新增迁移字段存在。
- 销售招募保持默认关闭；启用前需填写真实微信并确认合作规则。本次未修改线上业务配置，也未进行真实上游生成请求或线上配置保存验收。

## 2026-09-24：登录后销售伙伴招募

- 系统设置的站点设置新增招募开关、介绍首页文案与配图、佣金档位、合作权益、合作规则、微信号/二维码和预览，复用现有 settings JSON 存储；示例默认关闭，不新增佣金结算逻辑。
- 入口只在登录后工作台显示；默认首页、登录页和注册页隐藏。招募配置通过 JWT 鉴权接口 `/api/v1/settings/sales-recruitment` 获取，不进入公开配置或 HTML 注入，退出登录清除内容。
- 使用站内 `/v1/images/generations` 请求 `gpt-image-2` 生成原创配图，发布资源为 `frontend/src/assets/images/sales-partner.jpg`，同目录保留提示词；真实微信二维码由站长上传。
- 验证通过：后端设置相关现有测试、`go build ./...`；前端类型检查、改动文件 lint、83 项现有测试及生产构建。浏览器验证三页中文、390px 窄屏、暗色、复制和 Escape；本地模拟登录/接口验证匿名零请求、登录展示、首页隐藏及退出清理。未做线上配置保存验收。
- 已随上述线上2发布完成提交与部署；启用前需填写真实微信并确认示例佣金规则。

## 2026-09-24：接入图片创作工作台

- 根据用户反馈将模型改为密钥驱动的下拉框：自动以所选密钥请求 `/v1/models` 并筛选图片模型，默认选择第一项；切换密钥清空旧选项，取消旧请求。加载、空列表和失败分别提示，失败可重试，不再接受手工模型 ID；提交和历史回填均遵守当前可选模型范围。
- 基于 `mine` 的 `3788e0361`，在独立 worktree `sub2api-klno-image-workspace`、分支 `feat/image-workspace` 实现 `/image-workspace`；复用用户布局、密钥、模型、计费与现有图片网关，提供生成、上传/拖拽/粘贴改图、预览、下载及继续修改，支持中英文、深色和移动端。
- 已登录配置接口只返回现有异步存储开关；启用时复用异步任务，按用户在浏览器保存 24 小时内的任务元数据，不保存密钥和图片；未启用时使用同步接口并提示刷新会丢失结果。增加按任务原密钥鉴权的图片下载接口，复用现有下载器的超时和大小限制。
- 查询失败停止自动轮询，手动恢复仅 GET；提交不自动重试。历史损坏或初始化失败时保留原记录，防止清理定时器覆盖未恢复的任务。
- 真实 Quantix 验收：`openai` 分组完成生成和改图，均 HTTP 200、1254×1254，下载和继续修改通过。CPR 仅接受 JSON 改图，OpenAI 改图使用 `images[].image_url` data URL；Grok 保留 multipart。用户授权后将 `openai`、`openai pro20 号池`、`gpt-不降智` 三个 OpenAI 分组的图片生成开关全部开启并保留。
- 当前站点异步图片存储关闭；真实验收通过本地 Vite 代理访问 Quantix，仅未部署的配置接口按实际状态返回关闭。异步刷新、轮询失败后手动恢复、失败任务、过期清理、损坏记录和初始化失败跨清理周期保留历史，以浏览器拦截响应验证；未宣称对象存储已做真实验收。
- 验证通过：前端类型检查、改动文件 ESLint、6 项协议及语言键测试、生产构建；后端图片任务相关 handler/service/routes 测试及 `go build ./...`；390px 页面无横向溢出。生产构建仍有现有大包和混合导入提示。
- 剩余外部限制：`gpt-不降智` 上游账号「天才少年中转」对生图返回 403，网关表现为 502（运维错误 46850/46851），开启分组权限不能修复上游拒绝。用户已授权提交并确认合入当前开发主线 `mine`，不推送或部署。

## 2026-09-23：保留现有补丁，直接合并官方上游

- `mine` 直接合并 `Wei-Shaw/sub2api` 的 `upstream/main`（本次为 `a3eb7ef30`）；移除本地 `kin` 远程，保留已有提交和 `klno` 历史分支。
- `sync.sh` 和 `sync-upstream.yml` 改用普通合并，不再重建或重放 `klno`。自动同步仅在验证通过后推送 `mine`，不自动打标签或发布。定时工作流需要存在于远端默认分支才会生效；本次仅完成本地合并。
- 解决 Codex 身份、CPR、额度与转发链路冲突：保留账号隔离和 turn-state 功能，并兼容官方的心跳重试、调度倍率、邀请接口和新服务生命周期。
- metadata 仅将原始 Unicode 转成 JSON 转义，保留字段顺序、未知字段和数值精度，同时满足官方的 HTTP 头安全约束。
- 验证通过：`go build ./...`、`go test -tags unit ./...`、前端生产构建、lint、431 项关键测试，以及同步工作流的 YAML 和 Shell 语法检查。依赖注入代码已通过 Wire 重新生成；未执行部署或真实上游请求。
