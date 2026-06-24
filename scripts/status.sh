#!/bin/bash
# filepath: platserver/scripts/status.sh

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
BIN_DIR="$PROJECT_DIR/bin"
LOGS_DIR="$PROJECT_DIR/logs"
PID_FILE="$BIN_DIR/platserver.pid"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log_info()   { echo -e "${GREEN}[INFO]${NC} $1"; }
log_warn()   { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_error()  { echo -e "${RED}[ERROR]${NC} $1"; }
log_detail() { echo -e "${BLUE}[DETAIL]${NC} $1"; }

check_process_status() {
    if [ ! -f "$PID_FILE" ]; then
        log_error "PlatServer 未运行 (PID文件不存在)"
        return 1
    fi

    local pid=$(cat "$PID_FILE")
    if ! kill -0 "$pid" 2>/dev/null; then
        log_error "PlatServer 未运行 (PID已失效)"
        rm -f "$PID_FILE"
        return 1
    fi

    log_info "PlatServer 正在运行"
    log_detail "PID: $pid"

    if command -v ps >/dev/null 2>&1; then
        local info=$(ps -p "$pid" -o pid,ppid,user,pcpu,pmem,etime,cmd --no-headers 2>/dev/null || true)
        [ -n "$info" ] && log_detail "进程信息: $info"
    fi

    return 0
}

check_config_file() {
    log_info "=== 配置文件 ==="
    if [ -f "$BIN_DIR/.devops.yaml" ]; then
        log_detail "配置文件: $BIN_DIR/.devops.yaml"
        local cluster=$(grep "^cluster:" "$BIN_DIR/.devops.yaml" | awk '{print $2}')
        local http_port=$(grep "http_listen_addr:" "$BIN_DIR/.devops.yaml" | awk -F: '{print $NF}' | tr -d ' ' | cut -d'#' -f1)
        local rpc_port=$(grep "address:" "$BIN_DIR/.devops.yaml" | head -1 | awk -F: '{print $NF}' | tr -d ' "' | cut -d'#' -f1)
        local pprof_port=$(grep -A2 "^pprof:" "$BIN_DIR/.devops.yaml" | grep "addr:" | awk -F: '{print $NF}' | tr -d ' "' | cut -d'#' -f1)
        log_detail "集群: $cluster"
        log_detail "HTTP端口: $http_port"
        log_detail "RPC端口: $rpc_port"
        [ -n "$pprof_port" ] && log_detail "Pprof端口: $pprof_port"
    else
        log_warn "配置文件不存在: $BIN_DIR/.devops.yaml"
    fi
}

check_network_ports() {
    log_info "=== 网络端口 ==="
    local http_port rpc_port pprof_port port_pattern
    if [ -f "$BIN_DIR/.devops.yaml" ]; then
        http_port=$(grep "http_listen_addr:" "$BIN_DIR/.devops.yaml" | awk -F: '{print $NF}' | tr -d ' ' | cut -d'#' -f1)
        rpc_port=$(grep "address:" "$BIN_DIR/.devops.yaml" | head -1 | awk -F: '{print $NF}' | tr -d ' "' | cut -d'#' -f1)
        pprof_port=$(grep -A2 "^pprof:" "$BIN_DIR/.devops.yaml" | grep "addr:" | awk -F: '{print $NF}' | tr -d ' "' | cut -d'#' -f1)
    fi
    port_pattern=""
    [ -n "$http_port" ]  && port_pattern=":$http_port"
    [ -n "$rpc_port" ]   && port_pattern="$port_pattern|:$rpc_port"
    [ -n "$pprof_port" ] && port_pattern="$port_pattern|:$pprof_port"
    port_pattern=${port_pattern#|}
    [ -z "$port_pattern" ] && { log_warn "无法从配置文件读取端口信息"; return; }

    if command -v ss >/dev/null 2>&1; then
        log_detail "监听端口:"
        ss -tlnp 2>/dev/null | grep -E "($port_pattern)" || log_warn "未找到PlatServer监听端口"
    elif command -v netstat >/dev/null 2>&1; then
        log_detail "监听端口:"
        netstat -an 2>/dev/null | grep -E "($port_pattern)" | grep -i "LISTEN" || log_warn "未找到PlatServer监听端口"
    else
        log_warn "netstat 和 ss 均不可用，无法检查端口"
    fi
}

check_log_files() {
    log_info "=== 日志文件 ==="
    if [ -d "$LOGS_DIR" ]; then
        log_detail "日志目录: $LOGS_DIR"
        local logs=$(ls -lh "$LOGS_DIR"/*.log 2>/dev/null | tail -5)
        [ -n "$logs" ] && echo "$logs" || log_warn "日志目录为空"
    else
        log_warn "日志目录不存在: $LOGS_DIR"
    fi
}

main() {
    echo ""
    log_info "=== PlatServer 状态检查 ==="
    echo ""

    check_process_status
    local status=$?

    echo ""
    check_config_file

    if [ $status -eq 0 ]; then
        echo ""
        check_network_ports
        echo ""
        check_log_files
    fi

    echo ""
    log_info "=== 状态检查完成 ==="
    echo ""

    exit $status
}

main "$@"
