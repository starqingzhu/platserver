# PlatServer 多游戏平台化评估报告

## 1. 结论

**PlatServer 无需进行多游戏平台化改造。**

PlatServer 是纯粹的**进程级告警/通知网关**，不持有任何游戏业务数据，也不维护玩家状态。其职责是将内部告警通过企业微信/机器人推送给运维人员，天然与 gameId 无关。

---

## 2. 服务定位

PlatServer 当前仅提供两项 RPC 接口：

| 接口 | 用途 | 数据范围 |
|------|------|---------|
| `S2SPlayer2RobotNotice` | 玩家触发告警（如异常充值、作弊检测） | 单条通知消息 |
| `S2SGmRobotNotice` | GM 操作告警（如大额补发、敏感配置变更） | 单条通知消息 |

通知内容包含 `userId`、`bizType`、`body`，运维人员通过 `bizType` 即可识别业务场景。不同游戏的告警可以推送到**同一个** webhook，也可以按 `bizType` 路由到不同的 webhook（当前已实现）。

---

## 3. 为何无需 gameId 隔离

### 3.1 无持久化数据

PlatServer 内部无任何 Redis/MongoDB 读写：

```go
// platserver/internal/module/notify/notifier.go
func (n *RobotNotifier) send(userId string, bizType, body string) error {
    msgType, payload := getBuilder(bizType).Build(userId, bizType, body)
    // 直接发 HTTP POST 到 webhook，无缓存、无数据库
}
```

### 3.2 无玩家会话状态

PlatServer 不管理玩家在线状态、不维护玩家映射关系、不存在"同一 userId 跨游戏冲突"问题。

### 3.3 通知是运维行为，非游戏行为

接收通知的是**运维人员/开发团队**，不是玩家。运维人员关心的是"哪里出了问题"，这个问题本身已经通过 `bizType` 和 `body` 描述清楚，无需额外用 gameId 隔离通知通道。

### 3.4 Webhook 配置是进程级

```go
// platserver/internal/module/notify/config.go
type NotifyConfig struct {
    WebhookURL string `json:"webhookUrl"` // 单个进程一个 webhook
}
```

若未来需要不同游戏的告警发送到不同群，应在**配置层**按 `bizType` 路由 webhook（当前已支持），而非按 `gameId` 隔离进程。

---

## 4. 可选增强（非必须）

虽然 PlatServer 无需改造，但出于**日志可追溯性**考虑，建议在通知消息中附加 `gameId`：

### 4.1 GameServer 调用侧（已有 metadata）

GameServer 调用 PlatServer 时已经携带 `x-game-id` metadata：

```go
// gameserver/internal/router/rpc/platserver/client.go
func (h *ClientHandler) SendRobotNotice(userId int64, bizType, body string) {
    ctx := WithGameId(context.Background(), globalGameId) // 若已有此封装
    _, err := ssmsg.NewPlatServerClient(conn).S2SPlayer2RobotNotice(ctx, &ssmsg.PBS2SRobotNoticeRequest{...})
}
```

### 4.2 PlatServer 日志增强

PlatServer 可在 RPC Server 中提取 gameId 并记录到日志，方便排查：

```go
// platserver/internal/router/rpc/plat/server.go
func (h *ServerHandler) S2SPlayer2RobotNotice(ctx context.Context, req *pb.PBS2SRobotNoticeRequest) (*pb.PBS2RobotNoticeResponse, error) {
    gameId := extractGameId(ctx) // 从 metadata 提取
    zaplog.LoggerSugar.Infof("[PlatServer] S2SPlayer2RobotNotice userId=%d gameId=%d bizType=%s", req.UserId, gameId, req.BizType)
    // ... 其余逻辑不变
}
```

> 此为**日志增强**，不涉及数据隔离，不影响业务逻辑，可随时独立上线。

---

## 5. 与其他服务的对比

| 服务 | 是否需改造 | 原因 |
|------|-----------|------|
| **GameServer** | ✅ 是 | 持有玩家会话、游戏数据，必须隔离 |
| **PayServer** | ✅ 是 | 持有订单数据，同一 userId 跨游戏订单必须隔离 |
| **SocialServer** | ✅ 是 | 持有排行榜数据，同名活动跨游戏必须隔离 |
| **GlobalServer** | ✅ 是 | LBS 路由和客户端配置必须按 gameId 隔离 |
| **PlatServer** | ❌ 否 | 无状态通知网关，无数据隔离需求 |

---

## 6. 总结

PlatServer 作为无状态告警推送服务，在多游戏平台架构中保持现状即可。它的"租户"是**接收通知的运维团队**，不是**游戏**。若未来业务扩展（如按游戏分 webhook），通过配置层 `bizType → webhook` 映射即可满足，无需引入 gameId 概念到核心流程。
