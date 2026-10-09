# lol-teammate-helper

英雄联盟选人阶段的队友助手：通过本地 League Client (LCU) 接口读取选人会话，实时展示队友的近期排位战绩与段位。

基于 [Wails](https://wails.io)（Go）+ Vue 3 + TypeScript + Tailwind。

## 功能

- 自动检测运行中的英雄联盟客户端（端口 / token / 大区），客户端重启后自动重连
- 监听选人会话（`lol-champ-select`），展示队友近 5 场战绩与当前英雄
- 个人资料与排位段位总览

## 本地安装与启动（Windows）

前置条件：

| 工具 | 版本 | 说明 |
| --- | --- | --- |
| Go | 1.23+ | https://go.dev/dl/ |
| Node.js | 22 LTS | 自带 npm |
| Wails CLI | v2.10.2 | `go install github.com/wailsapp/wails/v2/cmd/wails@v2.10.2` |
| WebView2 | - | Windows 11 自带，Windows 10 需要安装 |

```bash
git clone https://github.com/MrKangM/lol-teammate-helper.git
cd lol-teammate-helper
git checkout claude/gallant-lamport-jdjbpz   # PR 分支，合并后用 main

cd frontend && npm ci && cd ..               # 前端依赖
wails doctor                                  # 可选：检查环境
wails dev                                     # 开发模式（热更新，终端直接输出日志）
```

`wails build` 会生成 `build/bin/lol-teammate-helper.exe`。

使用方法：先登录并保持英雄联盟客户端运行，再启动本软件（顺序反过来也可以，软件会自动重连）。进入英雄选择后自动显示队友；进入加载界面后可切换到“对手”。

## 日志

- `wails dev`：日志直接打印在运行命令的终端里。
- 打包后的 exe 没有控制台，日志写入 `%AppData%\lol-teammate-helper\logs\app.log`（Linux/macOS 为 `~/.config/lol-teammate-helper/logs/app.log`），超过 5 MB 会重新开始。
- 默认只记录 INFO 及以上。需要看每一次 LCU 请求时，设置环境变量 `LTH_LOG_LEVEL=debug` 再启动：

```powershell
$env:LTH_LOG_LEVEL = "debug"; wails dev
```

### 客户端检测不到

日志出现 `league client not detected` 时，后面的 `err=` 就是原因：

- `league client is not running`：没找到 `LeagueClientUx.exe` 进程。
- `league client found but its command line is unreadable`：客户端以管理员权限运行，普通权限读不到它的启动参数。**用管理员身份打开终端再运行 `wails dev`**，或者设置环境变量 `LTH_LOL_DIR` 指向包含客户端 `lockfile` 文件的目录（通常是客户端安装目录下的 `LeagueClient` 文件夹），程序会改读 lockfile。

排查时关注这几类日志：`connected to league client`（已连上客户端）、`gameflow phase`（游戏阶段）、`snapshot published`（已向界面推送数据）、以及 `failed` / `Warn` 级别的报错。

测试（后端纯逻辑部分）：

```bash
go test ./internal/...
```

## 构建

```bash
wails build
```

## 目录结构

| 路径 | 说明 |
| --- | --- |
| `app.go` / `main.go` | Wails 入口与前端绑定 |
| `internal/lcu` | 探测客户端进程并解析凭据 |
| `internal/connector` | WebSocket 连接、订阅、断线重连 |
| `internal/dispatch` | 事件分发，组装选人快照并推送给前端 |
| `internal/service` | 战绩 / 英雄数据查询与缓存 |
| `frontend/` | Vue 前端 |

## 免责声明

本工具只读取客户端在本机暴露的官方 LCU 接口，不读写游戏内存、不修改客户端。与 Riot Games 无关联。
