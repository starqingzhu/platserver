#!/bin/bash
# filepath: platserver/scripts/restart.sh

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log_info()  { echo -e "${GREEN}[INFO]${NC} $1"; }
log_warn()  { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }

main() {
    log_info "=== PlatServer 重启 ==="

    if "$SCRIPT_DIR/status.sh" >/dev/null 2>&1; then
        log_info "PlatServer 正在运行，即将停止..."
        "$SCRIPT_DIR/stop.sh"
    else
        log_info "PlatServer 未运行"
    fi

    sleep 3

    log_info "启动 PlatServer..."
    "$SCRIPT_DIR/start.sh"

    sleep 2
    if "$SCRIPT_DIR/status.sh" >/dev/null 2>&1; then
        log_info "PlatServer 重启成功"
    else
        log_error "PlatServer 重启失败"
        exit 1
    fi
}

main "$@"
