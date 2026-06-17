# Carbon Emission Management System

## English | [中文](#中文)

An enterprise-grade carbon emission management platform for measuring, reporting, and reducing greenhouse gas (GHG) emissions across **Scope 1, 2, and 3**, built to support ESG disclosure standards (GHG Protocol, ISO 14064, CDP, CSRD).

### Features
- **Scope 1/2/3 emission tracking** — full coverage of stationary & mobile combustion, purchased electricity/heat, and value-chain emissions
- **Emission factor library** — pluggable factors from IPCC, DEFRA, EPA, and custom sources with validity windows
- **Automated CO2e calculation** — activity data × factor → CO2-equivalent mass, with batch processing
- **Carbon credit / offset management** — track CER/ERU/VER credits and retire them against net emissions
- **Reduction target planning** — set baseline years, target percentages, and monitor progress (on-track / at-risk / achieved)
- **ESG & carbon reporting** — generate disclosure-ready reports by period and standard
- **Multi-organization & multi-facility** — model companies, sites, and vehicle fleets with geolocation
- **Analytics & dashboards** — scope breakdown, trend analysis, baseline comparison, facility-level views
- **Role-based access control** — admin / manager / analyst / viewer with JWT authentication
- **Full audit trail** — every create/update/delete/retire action is logged for compliance

**Ops & management (v2)**
- **Data import/export** — CSV/Excel bulk jobs for any entity
- **Scheduled tasks** — cron-driven recurring jobs with run/pause control
- **Alerts & warnings** — threshold/anomaly/trend alerts with ack & resolve workflow
- **Notifications** — email/sms/webhook dispatch with templates and retries
- **API key management** — programmatic access scoped per organization (X-API-Key header)
- **Webhook subscriptions** — event-driven outbound delivery with failure tracking
- **File attachments** — polymorphic file storage linked to any entity
- **Report exports** — render carbon reports to PDF/Excel
- **Audit rollback** — snapshot-based undo of destructive changes (admin)
- **System settings** — runtime configuration store grouped by category

### Tech Stack
- Go 1.22 + Gin (HTTP framework)
- PostgreSQL 16 (primary store) + Redis 7 (cache/queue)
- `database/sql` + `lib/pq` (no ORM, raw SQL control)
- `golang-jwt/v5` for auth, `log/slog` for structured JSON logging
- Docker Compose for one-command infrastructure

### Project Structure
```
carbon-emission-management/
├── cmd/api/                 # Entry point: config load, DI, graceful shutdown
├── internal/
│   ├── config/              # YAML config + DSN()/Addr() builders
│   ├── database/            # PostgreSQL connection pool
│   ├── cache/               # Redis client
│   ├── model/               # Domain models (20 entities)
│   ├── repository/          # Data access: generic CRUD + per-entity repos + analytics
│   ├── handler/             # HTTP handlers (real DB-backed, replaces stubs)
│   ├── server/              # Gin engine + full route table
│   ├── service/             # Service context + health checks
│   └── middleware/          # Auth (JWT), RBAC, CORS, API-key auth
├── pkg/                     # Reusable helpers: response, logger, jwt, password
├── configs/config.yaml      # Runtime configuration
├── sql/init.sql             # Full schema (20 tables) + seed data
├── Dockerfile               # Multi-stage container build
└── docker-compose.yml       # postgres + redis stack
```

### Quick Start
```bash
# 1. Start PostgreSQL + Redis
docker-compose up -d

# 2. Run the API (auto-loads sql/init.sql on first start)
go run cmd/api/main.go

# 3. Health check
curl http://localhost:8080/health
```

Default admin credentials (seeded): `admin` / `admin123`.

### API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| POST | /api/v1/auth/login | Login and obtain JWT |
| POST | /api/v1/auth/register | Register a new user |
| GET | /api/v1/auth/profile | Current user profile |
| GET/POST | /api/v1/organizations | List / create organizations |
| GET/PUT/DELETE | /api/v1/organizations/:id | Read / update / delete |
| GET/POST | /api/v1/facilities | List / create facilities |
| GET/PUT/DELETE | /api/v1/facilities/:id | Read / update / delete |
| GET | /api/v1/organizations/:id/facilities | Facilities by organization |
| GET/POST | /api/v1/emission-sources | List / create emission sources |
| GET/PUT/DELETE | /api/v1/emission-sources/:id | Read / update / delete |
| GET/POST | /api/v1/emission-factors | List / create emission factors |
| GET/PUT/DELETE | /api/v1/emission-factors/:id | Read / update / delete |
| GET/POST | /api/v1/emission-records | List / create emission records |
| GET/PUT/DELETE | /api/v1/emission-records/:id | Read / update / delete |
| GET | /api/v1/facilities/:id/emission-records | Records by facility |
| POST | /api/v1/emission-records/calculate | Calculate CO2e from activity × factor |
| GET/POST | /api/v1/carbon-credits | List / create carbon credits |
| POST | /api/v1/carbon-credits/:id/retire | Retire a credit against emissions |
| GET/POST | /api/v1/reduction-targets | List / create reduction targets |
| GET/POST | /api/v1/reports | List / create carbon reports |
| POST | /api/v1/reports/:id/generate | Compile a report from records |
| GET | /api/v1/analytics/dashboard | Summary KPIs (total/scope/offsets/net) |
| GET | /api/v1/analytics/by-scope | Emissions grouped by scope |
| GET | /api/v1/analytics/trend | Emissions over time |
| GET | /api/v1/analytics/comparison | Baseline vs current |
| GET | /api/v1/analytics/facility-breakdown | Emissions per facility |
| GET | /api/v1/audit-logs | Audit trail (admin only) |

#### Ops & Management (v2)

| Method | Path | Description |
|--------|------|-------------|
| POST/GET | /api/v1/data/imports | Create / list import jobs (CSV/Excel) |
| POST | /api/v1/data/imports/:id/process | Process an import job |
| GET | /api/v1/data/exports/:entity | Export an entity to CSV |
| GET/POST | /api/v1/scheduled-tasks | List / create cron-driven tasks |
| POST | /api/v1/scheduled-tasks/:id/run | Trigger a task immediately |
| POST | /api/v1/scheduled-tasks/:id/pause | Pause a task |
| GET/POST | /api/v1/alerts | List / create alerts |
| POST | /api/v1/alerts/:id/acknowledge | Acknowledge an alert |
| POST | /api/v1/alerts/:id/resolve | Resolve an alert |
| GET/POST | /api/v1/notifications | List / create notifications |
| POST | /api/v1/notifications/:id/send | Dispatch a notification |
| GET | /api/v1/notifications/templates | Message templates |
| GET/POST | /api/v1/api-keys | List / create API keys |
| POST | /api/v1/api-keys/:id/revoke | Revoke an API key |
| GET/POST | /api/v1/webhooks | List / create webhook subscriptions |
| POST | /api/v1/webhooks/:id/test | Fire a test delivery |
| GET/POST | /api/v1/attachments | List / create attachments |
| DELETE | /api/v1/attachments/:id | Delete an attachment |
| GET | /api/v1/report-exports | List report exports |
| POST | /api/v1/reports/:id/export | Render a report to PDF/Excel |
| GET | /api/v1/rollbacks | List rollback records (admin) |
| POST | /api/v1/audit-logs/:id/rollback | Roll back a change (admin) |
| GET/POST | /api/v1/settings | List / create system settings |
| PUT/DELETE | /api/v1/settings/:id | Update / delete a setting |

### Build & Test
```bash
make build      # compile binary to ./bin/
make test       # run all tests
make docker-up  # start infra containers
make lint       # golangci-lint
```

### License
MIT — see [LICENSE](LICENSE).

---

<a id="中文"></a>
# 碳排放管理系统

基于 Go + Gin + PostgreSQL + Redis 构建的企业级碳排放管理平台，覆盖 **范围 1/2/3** 全口径温室气体排放，支持 ESG 披露标准（GHG Protocol、ISO 14064、CDP、CSRD）。

### 功能特性
- **范围 1/2/3 排放追踪** — 固定/移动燃烧、外购电力热力、价值链排放全覆盖
- **排放因子库** — 支持 IPCC、DEFRA、EPA 及自定义因子，含有效期管理
- **自动 CO2e 计算** — 活动数据 × 排放因子 → CO2 当量，支持批量处理
- **碳信用/抵消管理** — 跟踪 CER/ERU/VER 信用并在净排放中注销
- **减排目标规划** — 设置基准年、减排百分比，监控进度（达标中/有风险/已达成）
- **ESG 与碳报告** — 按期间与披露标准生成报告
- **多组织多设施** — 建模公司、厂区、车队，支持地理位置
- **分析与看板** — 范围分布、趋势分析、基线对比、设施级视图
- **基于角色的权限控制** — admin/manager/analyst/viewer，JWT 鉴权
- **完整审计日志** — 创建/更新/删除/注销等操作全记录

**运营与管理（v2）**
- **数据导入导出** — 任意实体的 CSV/Excel 批量任务
- **定时任务** — cron 驱动的周期任务，支持运行/暂停
- **预警告警** — 阈值/异常/趋势告警，支持确认与解决工作流
- **消息通知** — email/sms/webhook 发送，含模板与重试
- **API 密钥管理** — 按组织隔离的编程访问（X-API-Key 头）
- **Webhook 订阅** — 事件驱动的外发投递，含失败计数
- **文件附件** — 关联任意实体的多态文件存储
- **报告导出** — 碳报告渲染为 PDF/Excel
- **审计回滚** — 基于快照的破坏性变更撤销（管理员）
- **系统设置** — 按 category 分组的运行时配置中心

### 技术栈
- Go 1.22 + Gin（HTTP 框架）
- PostgreSQL 16（主存储）+ Redis 7（缓存/队列）
- `database/sql` + `lib/pq`（无 ORM，原生 SQL）
- `golang-jwt/v5` 鉴权，`log/slog` 结构化 JSON 日志
- Docker Compose 一键部署基础设施

### 项目结构
```
carbon-emission-management/
├── cmd/api/                 # 入口：加载配置、依赖注入、优雅关闭
├── internal/
│   ├── config/              # YAML 配置 + DSN()/Addr() 构造
│   ├── database/            # PostgreSQL 连接池
│   ├── cache/               # Redis 客户端
│   ├── model/               # 领域模型（20 个实体）
│   ├── repository/          # 数据访问层：泛型 CRUD + 各实体 repo + 分析查询
│   ├── handler/             # HTTP 处理器（真实数据库实现，替代 stub）
│   ├── server/              # Gin 引擎 + 完整路由表
│   ├── service/             # 服务上下文 + 健康检查
│   └── middleware/          # 鉴权(JWT)、RBAC、CORS、API Key 鉴权
├── pkg/                     # 可复用助手：response、logger、jwt、password
├── configs/config.yaml      # 运行时配置
├── sql/init.sql             # 完整 schema（20 张表）+ 种子数据
├── Dockerfile               # 多阶段容器构建
└── docker-compose.yml       # postgres + redis 栈
```

### 快速开始
```bash
# 1. 启动 PostgreSQL + Redis
docker-compose up -d

# 2. 运行 API（首次启动自动执行 sql/init.sql）
go run cmd/api/main.go

# 3. 健康检查
curl http://localhost:8080/health
```

默认管理员账号（种子数据）：`admin` / `admin123`。

API 端点详见上方英文版的表格。主要模块：

| 模块 | 说明 |
|------|------|
| 鉴权 | `POST /auth/login`（JWT 登录）、`POST /auth/register`、`GET /auth/profile` |
| 组织/设施 | 组织与设施的完整 CRUD，支持按组织查询设施 |
| 排放源/因子 | 排放源与排放因子管理 |
| 排放记录 | CRUD + `POST /emission-records/calculate`（活动量×因子→CO2e） |
| 碳信用 | CRUD + `POST /carbon-credits/:id/retire`（注销抵消） |
| 减排目标 | 基准年/目标年/达成进度管理 |
| 碳报告 | CRUD + `POST /reports/:id/generate`（聚合 scope1/2/3 生成报告） |
| 分析 | 仪表盘、按范围、趋势、基线对比、设施级分布 |
| 审计 | `GET /audit-logs`、`POST /audit-logs/:id/rollback`（管理员） |
| 运营管理(v2) | 导入导出、定时任务、告警、通知、API Key、Webhook、附件、报告导出、回滚、系统设置 |

所有列表端点支持 `?page=&page_size=` 分页，返回 `{data, total, page, page_size}`。除 `/auth/login`、`/auth/register` 外均需 JWT（或 `X-API-Key` 头）。

### 构建与测试
```bash
make build      # 编译二进制到 ./bin/
make test       # 运行所有测试
make docker-up  # 启动基础设施容器
make lint       # golangci-lint 检查
```

### 开源许可
MIT — 详见 [LICENSE](LICENSE)。
