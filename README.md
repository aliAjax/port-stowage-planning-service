# Port Stowage Planner

企业级纯 Go 港口集装箱配载与岸桥作业计划服务。服务围绕船舶靠泊、集装箱属性、舱位约束、危险品隔离和作业指令租约构建，默认设备适配器是内存模拟器，不连接真实港机。

## 快速开始

```bash
go test ./...
go vet ./...
go run ./cmd/planner
./scripts/smoke.sh
```

默认监听 `:8080`，可通过 `PORT_ADDR` 和 `PORT_SOLVE_BUDGET_MS` 覆盖配置。`/healthz`、`/readyz` 可用于编排探针。

## API 流程

1. `POST /api/v1/voyages` 创建靠泊航次。
2. `POST /api/v1/containers` 注册箱属性。
3. `POST /api/v1/plans/solve` 运行硬约束剪枝和确定性 first-fit 求解。
4. `POST /api/v1/plans/{id}/publish` 发布计划并生成带租约、序列号和不可变标志的作业指令。
5. `GET /api/v1/plans/{id}/simulate|explain|instructions` 查看仿真、解释和指令。

状态机为 draft -> simulated/review -> published -> frozen，并支持 published/frozen -> rolled_back。内容摘要由航次、版本和按容器 ID 排序的槽位决策计算，避免 map 遍历导致结果漂移。

## 设计要点

- 硬约束：重量上限、尺寸、甲板/舱位、冷藏供电、危险品行隔离、重复槽位剪枝。
- 优化目标：船期时间、翻箱数量、岸桥均衡和危险品风险，超预算时返回最佳已知解标志。
- 执行安全：发布时生成租约 token；模拟适配器检测重复指令和过期租约；已执行指令通过 `Immutable` 保持不可变。
- 目录按 `cmd`、`internal/domain`、`repository`、`solver`、`dispatch`、`simulation`、`transport` 分层，便于后续接入 PostgreSQL、OTel 和真实 gRPC 适配器。

## 生产化路线

当前仓储为内存实现，适合 0-1 验证；生产部署应替换持久化仓储、增加 protobuf 生成代码、OTel exporter、认证和审计存储，并将设备回执通过幂等事件表落库。

## 业务模块

- `internal/berth`：泊位分配与岸桥指派，处理船舶长度/吃水兼容与占用窗口冲突。
- `internal/yard`：堆场箱区分配，处理冷藏插座、危险品隔离与目的港分组。
- `internal/schedule`：潮汐窗口排程，按潮汐表插值计算可用吃水窗口。
- `internal/cargo`：超限箱（OOG）判定与甲板/舱盖约束。
- `internal/inspection`：闸口验箱，校验证书、封条与 VGM 重量偏差。
- `internal/lashing`：绑扎计算，按船舶运动与风力推导绑扎点需求。
- `internal/reefer`：冷藏箱温度监控与报警生命周期。
- `internal/document`：贝位图生成与翻箱统计。
- `webui`：嵌入式配载控制台页面（go:embed 静态资源）。
