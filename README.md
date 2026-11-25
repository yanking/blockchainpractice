# blockchainpractice

区块链练习。每个练习单独一个文件夹，互不依赖，各自带 `go.mod`。

## 练习列表

| 练习 | 目录 | 运行方式 |
| --- | --- | --- |
| PoW 工作量证明 | `pow/` | `cd pow && go run . -name <昵称> [-nonce <起始 nonce>] [-difficulty <难度列表>]` |

## 1. PoW 工作量证明

对 `昵称 + nonce` 做 sha256，不断递增 nonce，直到哈希值以指定个数的 0 开头。
每个难度目标打印耗时、被哈希的内容和哈希值；有多个目标时依次完成，后一个从前一个命中的 nonce 之后继续。

- `main.go`：程序入口，解析参数、计时、打印结果
- `pow.go`：核心逻辑 `Mine(昵称, 起始 nonce, 难度)`
- `pow_test.go`、`main_test.go`：单元测试

### 参数

- `-name`：昵称，必填
- `-nonce`：起始 nonce，默认 0
- `-difficulty`：难度目标，即哈希值开头 0 的个数，范围 1 到 64；多个用逗号分隔，默认 `4,5`

```bash
cd pow
go run . -name wangyan -nonce 0              # 默认难度 4,5
go run . -name wangyan -difficulty 6         # 只找 6 个 0 开头
go run . -name wangyan -difficulty 4,5,6     # 依次完成 4、5、6
go test
```

难度每加 1，平均尝试次数变为 16 倍。下表耗时按本机约 1200 万次/秒估算，实际耗时随机波动很大：

| 难度 | 平均尝试次数 | 平均耗时 |
| --- | --- | --- |
| 4 | 6.5 万 | 约 0.005 秒 |
| 5 | 104 万 | 约 0.1 秒 |
| 6 | 1677 万 | 约 1.5 秒 |
| 7 | 2.7 亿 | 约 20 秒 |
| 8 | 43 亿 | 约 6 分钟 |
