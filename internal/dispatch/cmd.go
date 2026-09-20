/**
 * @ Description: 业务命令名常量（跨传输层共享的唯一标识）
 *
 * 同一命令可同时由 RPC / HTTP 触发，传输层只负责把「命令名 + 载荷」送到
 * dispatch.Dispatcher，因此命令名在这里统一定义，避免各处手写字符串。
 */

package dispatch

const (
	CmdRobotNoticePlayer = "robotNotice.player" // 玩家维度机器人通知（RPC S2SPlayer2RobotNotice + HTTP GM 触发）
	CmdRobotNoticeGM     = "robotNotice.gm"      // GM 维度机器人通知（RPC S2SGmRobotNotice）
)
