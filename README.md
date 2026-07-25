<div align="center">

<img src="assets/logo.svg" alt="vmapi Logo" width="128" />

# vmapi

[![Go](https://img.shields.io/badge/Go-1.25.7-00ADD8.svg)](https://golang.org/)
[![Vue](https://img.shields.io/badge/Vue-3.4+-4FC08D.svg)](https://vuejs.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-336791.svg)](https://www.postgresql.org/)
[![Redis](https://img.shields.io/badge/Redis-7+-DC382D.svg)](https://redis.io/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED.svg)](https://www.docker.com/)

**AI API Gateway focused on multi-provider routing, billing, and video generation**

English | [中文](README_CN.md) | [日本語](README_JA.md)

</div>

## Overview

**vmapi** is an AI API gateway platform for distributing and managing multi-provider AI quotas.

Users call upstream services through platform-issued API Keys. The gateway handles authentication, billing, load balancing, sticky routing, and request forwarding.

This repository is the **vmapi** fork, with first-class **Seedance / video generation** support and custom branding.

## Features

- **Multi-Account Management** — OAuth and API Key upstream accounts
- **API Key Distribution** — issue and manage user API Keys
- **Precise Billing** — token-level usage tracking and cost calculation
- **Smart Scheduling** — intelligent account selection with sticky sessions
- **Concurrency & Rate Limits** — per-user / per-account controls
- **Built-in Payment** — EasyPay, Alipay, WeChat Pay, Stripe ([Payment Guide](docs/PAYMENT.md))
- **Admin Dashboard** — web console for ops, monitoring, and billing
- **Composite Groups** — multi-provider model routing ([Operator Guide](docs/COMPOSITE_GROUPS.md))
- **Seedance Video Gateway** — native video generation / status / asset upload paths
- **Grok / OpenAI / Claude / Gemini / Antigravity** — multi-platform gateway support

## Seedance Video Support

vmapi adds an independent `seedance` platform (not mounted under OpenAI/Grok).

**Client-facing entrypoint** is your own vmapi deployment, not the upstream vendor host:

```text
https://YOUR_VMAPI_HOST/v1/video/generations
https://YOUR_VMAPI_HOST/v1/video/generations/:task_id
https://YOUR_VMAPI_HOST/v1/assets/uploads
```

Local example: `http://127.0.0.1:1800` (or the port you publish).

| Item | Detail |
|------|--------|
| Platform | `seedance` |
| Create task | `POST /v1/video/generations` |
| Query task | `GET /v1/video/generations/:task_id` |
| Upload asset | `POST /v1/assets/uploads` |
| Account type | API Key only |
| Upstream base URL | configured per Seedance account (`credentials.base_url`); leave empty only if you intentionally use the built-in vendor fallback |
| Client auth | `Authorization: Bearer <vmapi_api_key>` |
| Billing | charged when task is accepted (`video` = CNY/sec, `per_request` = CNY/call) |

## Tech Stack

| Component | Technology |
|-----------|------------|
| Backend | Go 1.25.7, Gin, Ent |
| Frontend | Vue 3.4+, Vite 5+, TailwindCSS |
| Database | PostgreSQL 15+ |
| Cache / Queue | Redis 7+ |

## Quick Start (Docker Compose)

```bash
# clone
git clone https://github.com/mhgmh/vmapi.git
cd vmapi

# prepare env (example)
cp deploy/.env.example deploy/.env
# edit secrets: JWT_SECRET / TOTP_ENCRYPTION_KEY / POSTGRES_PASSWORD / admin password

# start
cd deploy
docker compose up -d

# logs
docker compose logs -f sub2api
```

Open the setup wizard:

```text
http://YOUR_SERVER_IP:8080
```

### Local instance compose (dev)

```bash
# from repo root
docker compose -f deploy/docker-compose.local-instance.yml up -d --build
```

## Nginx Reverse Proxy Note

If Nginx sits in front of vmapi and you use sticky session headers (e.g. Codex CLI `session_id`), enable:

```nginx
underscores_in_headers on;
```

Nginx drops underscore headers by default, which breaks sticky routing.

## Development

```bash
# backend unit tests
cd backend && go test -tags=unit ./...

# frontend
cd frontend
pnpm install
pnpm run dev
pnpm run build
```

See [DEV_GUIDE.md](DEV_GUIDE.md) for environment notes.

## Project Structure

```text
vmapi/
├── backend/                  # Go backend
│   ├── cmd/server/           # entrypoint
│   ├── internal/             # config / model / service / handler / gateway
│   └── resources/
├── frontend/                 # Vue 3 admin + user console
│   └── src/
│       ├── api/
│       ├── stores/
│       ├── views/
│       └── components/
├── assets/                   # brand assets (logo.svg / logo.png)
├── deploy/                   # Docker / install scripts
└── docs/                     # operator docs
```

## Important Notice

- Using this project may violate upstream providers' terms of service. Review them before use; all risk is on the operator/user.
- Comply with local laws and regulations. Unlawful use is prohibited.
- Provided for technical learning and self-hosted operation. Authors are not liable for account bans, outages, data loss, or other damages.
- Commercial operation based on this project is solely the operator's responsibility.

## License

This project is licensed under the [GNU Lesser General Public License v3.0](LICENSE) (or later).

Upstream project: [Wei-Shaw/sub2api](https://github.com/Wei-Shaw/sub2api) (LGPL-3.0). vmapi is a branded fork with additional Seedance video gateway work.

---

<div align="center">

**vmapi** · Multi-provider AI API gateway

</div>
