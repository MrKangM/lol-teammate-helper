# lol-teammate-helper

英雄联盟选人阶段的队友助手：通过本地 League Client (LCU) 接口读取选人会话，实时展示队友的近期排位战绩与段位。

基于 [Wails](https://wails.io)（Go）+ Vue 3 + TypeScript + Tailwind。

## 功能

- 自动检测运行中的英雄联盟客户端（端口 / token / 大区），客户端重启后自动重连
- 监听选人会话（`lol-champ-select`），展示队友近 5 场战绩与当前英雄
- 个人资料与排位段位总览

## 开发

```bash
cd frontend && npm ci && cd ..
wails dev
```

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
