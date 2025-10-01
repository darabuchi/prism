package node

import (
	"context"
	"time"

	"github.com/lazygophers/lrpc/middleware/xerror"
)

// DelayICMP 测试 ICMP 延迟（ping）
//
// 注意：大多数代理协议不支持 ICMP 协议，因为 ICMP 是网络层协议，
// 而代理通常工作在传输层（TCP/UDP）。
//
// 此方法提供了一个统一的接口，但实际上可能无法通过代理进行 ICMP 测试。
// 对于不支持 ICMP 的节点类型，将返回 ErrNotSupported 错误。
//
// 未来可能的实现方式：
//   - 对于支持的协议（如 WireGuard），使用系统 ping 命令
//   - 对于其他协议，返回不支持错误
//
// 参数:
//   - ctx: 上下文，用于控制超时和取消
//   - addr: 目标地址（IP 或域名）
//
// 返回:
//   - time.Duration: ICMP 延迟时间（往返时间）
//   - error: 测试失败或不支持时返回错误
//
// 示例:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//	delay, err := node.DelayICMP(ctx, "8.8.8.8")
//	if err != nil {
//	    log.Errorf("ICMP 延迟测试失败: %v", err)
//	    return
//	}
//	log.Infof("ICMP 延迟: %v", delay)
func (n *Node) DelayICMP(ctx context.Context, addr string) (time.Duration, error) {
	if addr == "" {
		return 0, xerror.New(xerror.ErrSystemError, "addr is required")
	}

	// 检查节点类型是否支持 ICMP
	// 目前大多数代理协议不支持 ICMP
	adapterType := n.adapter.Type().String()
	switch adapterType {
	case "WireGuard":
		// WireGuard 可能支持 ICMP，但需要特殊实现
		// TODO: 实现 WireGuard 的 ICMP 测试
		return 0, xerror.New(xerror.ErrSystemError, "ICMP test for WireGuard is not implemented yet")
	default:
		// 其他协议不支持 ICMP
		return 0, xerror.New(xerror.ErrSystemError, "ICMP is not supported by this proxy protocol")
	}
}
