/**
 * @ Description: 业务命令统一注册入口（传输无关）
 *
 * 注册时机：server.OnInit 中调用 RegisterAll(dispatch.Main)，先于 HTTP/RPC 对外提供服务。
 */

package handler

import dispatchproto "golib/dispatch/proto"

// RegisterAll 注册全部业务命令到派发器。重复注册返回 ErrDuplicate。
func RegisterAll(d *dispatchproto.Dispatcher) error {
	registers := []func(*dispatchproto.Dispatcher) error{
		registerRobotNotice,
	}
	for _, fn := range registers {
		if err := fn(d); err != nil {
			return err
		}
	}
	return nil
}
