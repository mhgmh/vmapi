# vmapi 项目开发指南

> 本文档记录项目环境配置、常见坑点和注意事项，供开发与运维参考。

## 一、项目基本信息

| 项目 | 说明 |
|------|------|
| **本仓库** | [mhgmh/vmapi](https://github.com/mhgmh/vmapi) |
| **上游参考** | Wei-Shaw/sub2api（LGPL-3.0） |
| **技术栈** | Go 后端 (Ent ORM + Gin) + Vue3 前端 (pnpm) |
| **数据库** | PostgreSQL 15/16 + Redis 7+ |
| **包管理** | 后端: go modules，前端: **pnpm**（不是 npm） |
| **特色增强** | Seedance 视频平台（`platform=seedance`）+ vmapi 品牌 |

## 二、本地环境配置

### PostgreSQL

| 配置项 | 建议值 |
|--------|--------|
| 端口 | 5432 |
| 数据库 | `vmapi` 或 `sub2api` |
| 用户/密码 | 与 compose / config 保持一致 |

### Redis

| 配置项 | 建议值 |
|--------|--------|
| 端口 | 6379 |
| 密码 | 按环境配置 |

### 开发工具

```bash
# golangci-lint v2.7
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.7

# pnpm（前端包管理）
npm install -g pnpm
```

## 三、常用命令

```bash
# 后端单元测试
cd backend && go test -tags=unit ./...

# 后端集成测试
cd backend && go test -tags=integration ./...

# 代码质量检查
cd backend && golangci-lint run ./...

# 前端依赖安装（必须用 pnpm）
cd frontend && pnpm install

# 前端开发 / 构建
cd frontend && pnpm run dev
cd frontend && pnpm run build
```

## 四、Docker 本地实例

```bash
# 仓库根目录
docker compose -f deploy/docker-compose.local-instance.yml up -d --build
```

默认本地入口可参考部署配置（常见为 `http://127.0.0.1:1800` 或 `8080`）。

## 五、常见坑点

### 1. pnpm-lock.yaml 必须同步提交

`package.json` 改依赖后，务必同步更新并提交 `pnpm-lock.yaml`，否则 `pnpm install --frozen-lockfile` 会失败。

### 2. Google OAuth 凭据不要硬编码

公开仓库已去掉内置 Google OAuth Client Secret。如需 Antigravity / Gemini CLI 相关能力，请通过环境变量配置：

- `ANTIGRAVITY_OAUTH_CLIENT_SECRET`
- `GEMINI_CLI_OAUTH_CLIENT_SECRET`

### 3. Seedance 路径注意单数 video

客户端应调用 **你自己的 vmapi 网关**，例如：

- `POST https://你的域名/v1/video/generations`
- `GET https://你的域名/v1/video/generations/:task_id`
- `POST https://你的域名/v1/assets/uploads`

注意是 `video`（单数），不是 Grok 的 `/videos`。

上游厂商地址写在 Seedance **账号** 的 `credentials.base_url` 里；对外文档不要把上游域名写成 vmapi 默认入口。

## 六、相关文档

- [README.md](README.md) / [README_CN.md](README_CN.md)
- [docs/](docs/) 运营与支付文档
- [deploy/](deploy/) Docker 与安装脚本
