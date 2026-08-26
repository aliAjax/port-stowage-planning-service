# port-stowage-planning-service__006

基于 Go 实现的港口集装箱配载计划 Web 项目，一款后端服务与配载控制台，完成船舶靠泊与岸桥指派、堆场箱区分配、潮汐窗口排程、配载求解、作业指令租约、绑扎校核与闸口验箱等港口作业功能。
## 构建镜像

请从**仓库根目录**执行；`benzhi.Dockerfile`、`build_benzhi_docker.sh`、`BENZHI_README.md` 均固定在该目录：

```bash
./build_benzhi_docker.sh <image-name> [linux/amd64|linux/arm64]
```

## 标准命令

```bash
go build ./...     # 编译
go run ./cmd/planner   # 启动
go test ./...      # 测试（如有）
```

## 环境

- 基础镜像: golang:1.22
- Go 模块目录: `.`
- 依赖已在镜像构建阶段预下载，容器内离线可用。
- 容器内工作目录: `/app`
