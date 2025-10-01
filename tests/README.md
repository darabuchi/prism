# 测试目录

## 说明

集成测试和专项测试。

## 子目录

### integration/ - 集成测试
完整的端到端测试，测试多个组件协同工作。

### stress/ - 压力测试
性能测试和负载测试。

### testutil/ - 测试工具
测试辅助函数和数据工厂。

## 单元测试

单元测试放在对应包的 `*_test.go` 文件中。

## 运行测试

```bash
# 运行所有测试
make test

# 运行单元测试
make test-unit

# 运行集成测试
make test-integration

# 生成覆盖率报告
make test-coverage
```

## 参考文档

- [测试策略](../docs/详细设计/测试策略.md)
