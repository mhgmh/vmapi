<div align="center">

<img src="assets/logo.svg" alt="vmapi Logo" width="128" />

# vmapi

[![Go](https://img.shields.io/badge/Go-1.25.7-00ADD8.svg)](https://golang.org/)
[![Vue](https://img.shields.io/badge/Vue-3.4+-4FC08D.svg)](https://vuejs.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-336791.svg)](https://www.postgresql.org/)
[![Redis](https://img.shields.io/badge/Redis-7+-DC382D.svg)](https://redis.io/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED.svg)](https://www.docker.com/)

**面向多厂商路由、计费与视频生成的 AI API 网关**

[English](README.md) | 中文 | [日本語](README_JA.md)

</div>

## 项目概述

**vmapi** 是一个 AI API 网关平台，用于分发和管理多厂商 AI 配额。

用户通过平台下发的 API Key 调用上游服务；网关负责鉴权、计费、负载均衡、粘性路由和请求转发。

本仓库是 **vmapi** 分支版本，重点增强了 **Seedance / 视频生成** 能力，并使用独立品牌素材。

## 核心功能

- **多账号管理** — 支持 OAuth、API Key 上游账号
- **API Key 分发** — 为用户生成和管理 API Key
- **精确计费** — Token 级用量追踪与成本计算
- **智能调度** — 智能选号 + 粘性会话
- **并发与限流** — 用户级 / 账号级控制
- **内置支付** — EasyPay、支付宝、微信、Stripe（[配置指南](docs/PAYMENT_CN.md)）
- **管理后台** — Web 控制台：运维、监控、账单
- **组合分组** — 多厂商模型路由（[运营指南](docs/COMPOSITE_GROUPS.md)）
- **Seedance 视频网关** — 原生视频生成 / 状态查询 / 素材上传
- **Grok / OpenAI / Claude / Gemini / Antigravity** — 多平台网关支持

## Seedance 视频能力

vmapi 增加了独立的 `seedance` 平台（不挂在 OpenAI/Grok 下）：

| 项目 | 说明 |
|------|------|
| 平台标识 | `seedance` |
| 创建任务 | `POST /v1/video/generations` |
| 查询任务 | `GET /v1/video/generations/:task_id` |
| 上传素材 | `POST /v1/assets/uploads` |
| 账号类型 | 仅 API Key |
| 默认 base_url | `https://api.7tai.cc/v1` |
| 鉴权 | `Authorization: Bearer <token>` |
| 计费 | 任务受理（返回 task_id）时计费；`video`=元/秒，`per_request`=元/次 |

## 技术栈

| 组件 | 技术 |
|------|------|
| 后端 | Go 1.25.7, Gin, Ent |
| 前端 | Vue 3.4+, Vite 5+, TailwindCSS |
| 数据库 | PostgreSQL 15+ |
| 缓存 / 队列 | Redis 7+ |

## 快速开始（Docker Compose）

```bash
# 克隆
git clone https://github.com/mhgmh/vmapi.git
cd vmapi

# 准备环境变量
cp deploy/.env.example deploy/.env
# 编辑密钥：JWT_SECRET / TOTP_ENCRYPTION_KEY / POSTGRES_PASSWORD / 管理员密码

# 启动
cd deploy
docker compose up -d

# 日志
docker compose logs -f sub2api
```

浏览器打开安装向导：

```text
http://你的服务器IP:8080
```

### 本地开发实例

```bash
# 仓库根目录
docker compose -f deploy/docker-compose.local-instance.yml up -d --build
```

## Nginx 反向代理注意

若 Nginx 前置 vmapi，并使用粘性会话请求头（如 Codex CLI 的 `session_id`），请开启：

```nginx
underscores_in_headers on;
```

Nginx 默认会丢弃含下划线的请求头，导致粘性路由失效。

## 开发

```bash
# 后端单测
cd backend && go test -tags=unit ./...

# 前端
cd frontend
pnpm install
pnpm run dev
pnpm run build
```

更多环境说明见 [DEV_GUIDE.md](DEV_GUIDE.md)。

## 项目结构

```text
vmapi/
├── backend/                  # Go 后端
│   ├── cmd/server/           # 入口
│   ├── internal/             # config / model / service / handler / gateway
│   └── resources/
├── frontend/                 # Vue 3 管理端 + 用户端
│   └── src/
│       ├── api/
│       ├── stores/
│       ├── views/
│       └── components/
├── assets/                   # 品牌素材（logo.svg / logo.png）
├── deploy/                   # Docker / 安装脚本
└── docs/                     # 运营文档
```

## 重要提醒

- 使用本项目可能违反上游服务商条款。请先阅读相关协议，风险由使用者自行承担。
- 请在符合当地法律法规的前提下使用，禁止违法用途。
- 本项目供技术学习与自建部署使用。作者不对封号、中断、数据丢失等损失负责。
- 基于本项目开展的商业运营，责任由运营方自行承担。

## 许可证

本项目基于 [GNU Lesser General Public License v3.0](LICENSE)（或更高版本）授权。

上游项目：[Wei-Shaw/sub2api](https://github.com/Wei-Shaw/sub2api)（LGPL-3.0）。vmapi 是带独立品牌与 Seedance 视频网关增强的分支版本。

---

<div align="center">

**vmapi** · 多厂商 AI API 网关

</div>
