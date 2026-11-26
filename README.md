# blockchainpractice

区块链练习。每个练习单独一个文件夹，互不依赖；Go 练习各自带 `go.mod`。

## 练习列表

| 练习 | 目录 | 运行方式 |
| --- | --- | --- |
| PoW 工作量证明 | `pow/` | `cd pow && go run . -name <昵称> [-nonce <起始 nonce>] [-difficulty <难度列表>]` |
| RSA 非对称加密 | `rsa-sign/` | `cd rsa-sign && go run . -name <昵称> [-nonce <起始 nonce>]` |
| Bank 合约 | `bank/` | Remix 编译部署（Solidity 0.8.39），MetaMask 存款、提取 |

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

## 2. RSA 非对称加密

生成一对 RSA 公私钥，用私钥对满足 PoW（哈希值 4 个 0 开头）的 `昵称 + nonce` 签名，再用公钥验证。
输出依次是密钥对（PEM 格式）、PoW 结果、签名和验证结果。

- `main.go`：程序入口，按「生成密钥对 → PoW → 签名 → 验证」的顺序执行并打印
- `rsa.go`：密钥生成、签名（SHA-256 + RSA PKCS#1 v1.5）、验证，以及密钥的 PEM 输出
- `pow.go`、`pow_test.go`：从 `pow/` 原样复制，用来找出要签名的内容
- `rsa_test.go`：单元测试

密钥每次运行时新生成，只存在于内存和输出里，不写入文件。

```bash
cd rsa-sign
go run . -name wangyan -nonce 0
go test
```

## 3. Bank 合约

接收 ETH 存款，记录每个地址的存款金额，维护存款金额前 3 名，只有管理员（部署合约的地址）能提取合约里的全部 ETH。

- `Bank.sol`：Bank 合约，单文件，Solidity 0.8.39，无外部依赖

### 接口

| 接口 | 说明 |
| --- | --- |
| `receive()` | 存款。用 MetaMask 向合约地址转账（不带 data）即可触发，金额累加到该地址名下并更新前 3 名；金额为 0 时 revert |
| `withdraw()` | 管理员提取合约里的全部 ETH，其他地址调用会 revert |
| `admin` | 管理员地址，即部署者 |
| `deposits(地址)` | 该地址的存款金额 |
| `top3(下标)` | 存款金额前 3 名的地址，下标 0 到 2，按金额从高到低；不足 3 人时空位为全 0 地址 |

金额相同的，已在榜上的地址排在前面。

### 运行方式

在 [Remix](https://remix.ethereum.org) 里编译部署，用 MetaMask 操作，钱包切到 Sepolia 测试网即可：

1. 打开 `bank/Bank.sol`，在 Solidity Compiler 里把编译器版本选为 `0.8.39`，点 Compile
2. 在 Deploy & Run Transactions 里把 Environment 选为 `Injected Provider - MetaMask`，连接钱包后点 Deploy
3. 在 MetaMask 里切换账户，向合约地址转账，每笔转账就是一次存款
4. 回到 Remix，在 Deployed Contracts 里点 `top3`（输入 0、1、2）和 `deposits`（输入地址）查看记录
5. 用部署时的账户点 `withdraw` 提取全部 ETH；换别的账户点会 revert
