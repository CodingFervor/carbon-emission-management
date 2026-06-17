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

### Tech Stack
- Go 1.22 + Gin (HTTP framework)
- PostgreSQL 16 (primary store) + Redis 7 (cache/queue)
- `database/sql` + `lib/pq` (no ORM, raw SQL control)
- `golang-jwt/v5` for auth, `log/slog` for structured JSON logging
- Docker Compose for one-command infrastructure

### Project Structure
```
carbon-emission-management/
├── cmd/api/                 # Entry point, route registration, handler stubs
├── internal/
│   ├── config/              # YAML config + DSN()/Addr() builders
│   ├── database/            # PostgreSQL connection pool
│   ├── cache/               # Redis client
│   ├── model/               # Domain models (10 entities)
│   ├── service/             # Service context + health checks
│   └── middleware/          # Auth, RBAC, CORS
├── pkg/                     # Reusable helpers: response, logger, jwt
├── configs/config.yaml      # Runtime configuration
├── sql/init.sql             # Full schema + seed data
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

### 技术栈
- Go 1.22 + Gin（HTTP 框架）
- PostgreSQL 16（主存储）+ Redis 7（缓存/队列）
- `database/sql` + `lib/pq`（无 ORM，原生 SQL）
- `golang-jwt/v5` 鉴权，`log/slog` 结构化 JSON 日志
- Docker Compose 一键部署基础设施

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

API 端点详见上方英文版的表格。

### 构建与测试
```bash
make build      # 编译二进制到 ./bin/
make test       # 运行所有测试
make docker-up  # 启动基础设施容器
make lint       # golangci-lint 检查
```

### 开源许可
MIT — 详见 [LICENSE](LICENSE)。
