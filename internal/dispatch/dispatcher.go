// Package dispatch 持有 platserver 的全局派发器实例。
package dispatch

import dispatchproto "golib/dispatch/proto"

var Main = dispatchproto.New()
