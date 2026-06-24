#!/bin/bash
# filepath: platserver/scripts/stop.sh

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
BIN_DIR="$PROJECT_DIR/bin"
PID_FILE="$BIN_DIR/platserver.pid"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log_info()  { echo -e "${GREEN}[INFO]${NC} $1"; }
log_warn()  { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }

stop_server() {
    if [ ! -f "$PID_FILE" ]; then
        log_warn "PlatServer 未运行 (PID文件不存在)"
        return 0
    fi

    local pid=$(cat "$PID_FILE")
    if ! kill -0 "$pid" 2>/dev/null; then
        log_warn "PlatServer 未运行 (PID已失效)"
        rm -f "$PID_FILE"
        return 0
    fi

    log_info "正在停止 PlatServer (PID: $pid)..."
    kill -INT "$pid" 2>/dev/null || true

    log_info "等待进程优雅退出..."
    local count=0
    while kill -0 "$pid" 2>/dev/null; do
        sleep 1
        count=$((count + 1))
        printf "."
        if [ $((count % 30)) -eq 0 ]; then
            echo ""
            log_info "已等待 ${count} 秒，进程仍在处理中..."
        fi
    done
    echo ""

    rm -f "$PID_FILE"
    log_info "PlatServer 已优雅退出 (总耗时: ${count} 秒)"
}

main() {
    log_info "=== PlatServer 停止 ==="
    stop_server
    log_info "=== 停止完成 ==="
}

main "$@"
