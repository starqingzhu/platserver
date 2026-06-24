#!/bin/bash
# filepath: platserver/scripts/start.sh

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
BIN_DIR="$PROJECT_DIR/bin"
LOGS_DIR="$PROJECT_DIR/logs"
PID_FILE="$BIN_DIR/platserver.pid"

case "$(uname -s)" in
    CYGWIN*|MINGW*|MSYS*)
        EXEC_FILE="$BIN_DIR/platserver.exe"
        ;;
    *)
        EXEC_FILE="$BIN_DIR/platserver"
        ;;
esac

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log_info()  { echo -e "${GREEN}[INFO]${NC} $1"; }
log_warn()  { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }

pre_start_check() {
    if [ ! -f "$EXEC_FILE" ]; then
        log_error "可执行文件不存在: $EXEC_FILE"
        log_info "请先运行编译脚本: ./scripts/build.sh"
        exit 1
    fi

    chmod +x "$EXEC_FILE"

    if [ -f "$PID_FILE" ] && kill -0 $(cat "$PID_FILE") 2>/dev/null; then
        log_warn "PlatServer 已在运行中，PID: $(cat $PID_FILE)"
        exit 1
    fi

    mkdir -p "$LOGS_DIR"
    mkdir -p "$BIN_DIR"

    if [ ! -f "$BIN_DIR/.devops.yaml" ]; then
        log_error "配置文件不存在: $BIN_DIR/.devops.yaml"
        log_info "请先运行编译脚本: ./scripts/build.sh"
        exit 1
    fi

    log_info "配置文件: bin/.devops.yaml"
}

start_server() {
    log_info "正在启动 PlatServer..."

    cd "$BIN_DIR"
    nohup "$EXEC_FILE" > "$LOGS_DIR/startup.log" 2>&1 &
    local pid=$!
    echo $pid > "$PID_FILE"

    sleep 3
    if kill -0 $pid 2>/dev/null; then
        log_info "PlatServer 启动成功"
        log_info "PID: $pid"
        log_info "日志目录: $LOGS_DIR/"
        log_info "配置文件: $BIN_DIR/.devops.yaml"
    else
        log_error "PlatServer 启动失败"
        log_error "查看启动日志: $LOGS_DIR/startup.log"
        rm -f "$PID_FILE"
        exit 1
    fi
}

main() {
    log_info "=== PlatServer 启动 ==="
    pre_start_check
    start_server
    log_info "=== 启动完成 ==="
}

main "$@"
