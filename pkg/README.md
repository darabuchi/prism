# 公共包

## 说明

可被外部项目引用的公共包，高度独立，无内部依赖。

## 子包

### parser/ - 订阅解析器
- `clash.go`: Clash YAML 解析
- `v2ray.go`: V2Ray JSON 解析
- `surge.go`: Surge 配置解析

### protocol/ - 协议实现
- `socks5/`: SOCKS5 协议
- `http/`: HTTP 代理协议
- `vmess/`: VMess 协议

## 设计原则

1. **独立性**: 不依赖 internal 包
2. **文档性**: 完善的 GoDoc 文档
3. **测试性**: 完整的单元测试
4. **稳定性**: API 保持向后兼容
