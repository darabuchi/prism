package node

import (
	"context"
	"net"
	"time"

	"github.com/lazygophers/lrpc/middleware/xerror"
	"github.com/lazygophers/utils/candy"
	"github.com/metacubex/mihomo/constant"
)

// DelayTCP 测试 TCP 连接延迟
//
// 通过节点代理建立到目标地址的 TCP 连接，测量连接建立的时间
// 这个方法测试的是纯 TCP 连接延迟，不包含应用层协议的开销
//
// 参数:
//   - ctx: 上下文，用于控制超时和取消
//   - addr: 目标地址（格式：host:port）
//
// 返回:
//   - time.Duration: TCP 连接延迟时间
//   - error: 测试失败时返回错误
//
// 示例:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//	delay, err := node.DelayTCP(ctx, "www.google.com:80")
//	if err != nil {
//	    log.Errorf("TCP 延迟测试失败: %v", err)
//	    return
//	}
//	log.Infof("TCP 延迟: %v", delay)
func (n *Node) DelayTCP(ctx context.Context, addr string) (time.Duration, error) {
	if addr == "" {
		return 0, xerror.New(xerror.ErrSystemError, "addr is required")
	}

	// 解析地址
	host, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		return 0, xerror.WrapError(err, xerror.ErrSystemError, "invalid addr format")
	}

	// 解析端口
	port := uint16(candy.ToInt(portStr))
	if port == 0 {
		return 0, xerror.New(xerror.ErrSystemError, "invalid port")
	}

	// 构建 Mihomo 连接元数据
	metadata := &constant.Metadata{
		NetWork: constant.TCP,
		Host:    host,
		DstPort: port,
	}

	// 记录开始时间
	start := time.Now()

	// 使用节点的 DialContext 建立 TCP 连接
	conn, err := n.adapter.DialContext(ctx, metadata)
	if err != nil {
		return 0, xerror.WrapError(err, xerror.ErrSystemError, "TCP connection failed")
	}
	defer conn.Close()

	// 计算延迟
	delay := time.Since(start)

	return delay, nil
}
