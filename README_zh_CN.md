# Flow Go SDK [![GoDoc](https://godoc.org/github.com/onflow/flow-go-sdk?status.svg)](https://godoc.org/github.com/onflow/flow-go-sdk)

Flow Go SDK 提供相关开发包帮助 Golang 开发者完成在 Flow network 上进行应用的构建和开发。

*注意: 这个 SDK 通过 [Flow Emulator](https://docs.onflow.org/devtools/emulator/)完成实现，并且可以被用于本地开发。*

## [英文](/README.md) | [中文（简体）](#)

## 什么是 Flow ?

Flow is a new blockchain for open worlds. Read more about it [here](https://github.com/onflow/flow).

## 目录

- [开始](#getting-started)
  - [安装](#installing)
  - [生成密钥](#generating-keys)
    - [支持的曲线](#supported-curves)
  - [创建一个帐户](#creating-an-account)
  - [签名一个交易](#signing-a-transaction)
    - [Flow中的签名是如何在工作的](#how-signatures-work-in-flow)
      - [一人一签](#single-party-single-signature)
      - [多人签名](#single-party-multiple-signatures)
      - [多方](#multiple-parties)
      - [多方参与，两个自动执行器](#multiple-parties-two-authorizers)
      - [多方参与，多个签名](#multiple-parties-multiple-signatures)
  - [发起一笔交易](#sending-a-transaction)
  - [查询交易结果](#querying-transaction-results)
  - [查询块](#querying-blocks)
  - [执行一个脚本](#executing-a-script)
  - [查询事件](#querying-events)
    - [事件查询格式](#event-query-format)
    - [事件的结果](#event-results)
  - [查询账户](#querying-accounts)
- [例子](#examples)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## 开始

### 安装

开始使用 SDK, 首先安装 Go 1.13+ 版本， 并且运行 `go get`:

```sh
go get github.com/onflow/flow-go-sdk
```

### 生成密钥

Flow 使用 [ECDSA](https://en.wikipedia.org/wiki/Elliptic_Curve_Digital_Signature_Algorithm) 
去控制用户的账户权限。 每一个密钥都使用了 `SHA2-256` 或者 `SHA3-256` 哈希算法实现了。

这里有关于如何生成 ECDSA，采用P-256 (secp256r1) 曲线及私钥的相关方法:

```go
import "github.com/onflow/flow-go-sdk/crypto"

// 种子 短语
// 注意: 这只是一个例子, 请使用安全的方式随机生成种子
seed := []byte("elephant ears space cowboy octopus rodeo potato cannon pineapple")

privateKey, err := crypto.GeneratePrivateKey(crypto.ECDSA_P256, seed)
```
这个私钥可以被编码成  bytes 类型 (i.e. for storage):

```go
encPrivateKey := privateKey.Encode()
```

可以为一个私钥创建一个对应的公钥:

```go
publicKey := privateKey.PublicKey()
```

#### 支持的曲线

例子部分使用了 ECDSA 密钥， 采用的曲线算法是 P-256 (secp256r1) 。
Flow 也支持 Bitcoin 和 Ethereum 所使用的 secp256k1 曲线算法

这里展示如何采用secp256k1 曲线生成一个 ECDSA 私钥:

```go
privateKey, err := crypto.GeneratePrivateKey(crypto.ECDSA_secp256k1, seed)
```

这里是支持的签名哈希算法的完整列表:[Flow 签名 & 哈希算法](https://github.com/onflow/flow/blob/master/docs/accounts-and-keys.md#supported-signature--hash-algorithms)

### 创建一个帐户

一旦你完成了 [生成一个密钥对](#generating-keys), 您可以使用它的公钥，创建一个新帐户。


```go
import (
    "github.com/onflow/flow-go-sdk"
    "github.com/onflow/flow-go-sdk/crypto"
    "github.com/onflow/flow-go-sdk/templates"
)

ctx := context.Background()

// 为账户生成一个新的私钥
// 注意: 这只是一个例子, 请使用安全的方式随机生成种子
seed := []byte("elephant ears space cowboy octopus rodeo potato cannon pineapple")
privateKey, _ := crypto.GeneratePrivateKey(crypto.ECDSA_P256, seed)

// 得到公钥
publicKey := privateKey.PublicKey()

// 从公钥中构建一个用户密钥
accountKey := flow.NewAccountKey().
    SetPublicKey(publicKey).
    SetHashAlgo(crypto.SHA3_256).        // SHA3_256 哈希算法生成的密钥对
    SetWeight(flow.AccountKeyWeightThreshold) // 授予这个密钥签名权重

// 生成一个帐户创建脚本
// 这将创建一个帐户，该帐户只有一个公钥，没有代码
script, _ := templates.CreateAccount([]*flow.AccountKey{accountKey}, nil)

// 连接到本地运行的模拟器
c, err := client.New("localhost:3569")
if err != nil {
    panic("failed to connect to emulator")
}

payer, payerKey, payerSigner := examples.ServiceAccount(c)

tx := flow.NewTransaction().
    SetScript(script).
    SetGasLimit(100).
    SetProposalKey(payer, payerKey.Index, payerKey.SequenceNumber).
    SetPayer(payer)

err = tx.SignEnvelope(payer, payerKey.Index, payerSigner)
if err != nil {
    panic("failed to sign transaction")
}

err = c.SendTransaction(ctx, *tx)
if err != nil {
    panic("failed to send transaction")
}

result, err := c.GetTransactionResult(ctx, tx.ID())
if err != nil {
    panic("failed to get transaction result")
}

var myAddress flow.Address

if result.Status == flow.TransactionStatusSealed {
    for _, event := range result.Events {
        if event.Type == flow.EventAccountCreated {
            accountCreatedEvent := flow.AccountCreatedEvent(event)
            myAddress = accountCreatedEvent.Address()
        }
	}
}
```

### 签名一个交易

下面是一个简单的例子使用 `crypto.PrivateKey` 签名一个交易.

```go
import (
    "github.com/onflow/flow-go-sdk"
    "github.com/onflow/flow-go-sdk/crypto"
)

var (
    myAddress    flow.Address
    myAccountKey flow.AccountKey
    myPrivateKey crypto.PrivateKey
)

tx := flow.NewTransaction().
    SetScript([]byte("transaction { execute { log(\"Hello, World!\") } }")).
    SetGasLimit(100).
    SetProposalKey(myAddress, myAccountKey.Index, myAccountKey.SequenceNumber).
    SetPayer(myAddress)
```

交易签名通过 `crypto.Signer` 接口完成.  `crypto.Signer` 最安全、最简单的接口实现方法是 `crypto.InMemorySigner`.

使用这样的签名，可以更安全地生成密钥存储在硬件设备中，如 [HSM](https://en.wikipedia.org/wiki/Hardware_security_module).  `crypto.Signer` 
就是这样的签名接口的简单实现。

```go
// 通过你的私钥构造一个签名器，通过哈希算法完成签名
mySigner := crypto.NewInMemorySigner(myPrivateKey, myAccountKey.HashAlgo)

err := tx.SignEnvelope(myAddress, myAccountKey.Index, mySigner)
if err != nil {
    panic("failed to sign transaction")
}
```

#### Flow 中的签名是如何在工作的

Flow 引入了新的概念，允许在创建和签署事务时具有更大的灵活性。

在尝试下面的示例之前，我们建议您阅读[事务签名文档](https://github.com/onflow/flow/blob/master/docs/accounts-and-keys.md#signing-a-transaction).

---

##### [一人一签](https://github.com/onflow/flow/blob/master/docs/accounts-and-keys.md#single-party-single-signature)

- Proposer, payer, authorizer是同一个account (`0x01`)
- 只有信封(envolope)必须要被签名
- 提案密钥必须有充分的签名权重。

| Account   | Key ID | Weight |
|-----------|--------|--------|
| `0x01`    | 1      | 1.0    |

```go
account1, _ := c.GetAccount(ctx, flow.HexToAddress("01"))

key1 := account1.Keys[0]

// 安全的通过私钥创建一个签名器
key1Signer := getSignerForKey1()

tx := flow.NewTransaction().
    SetScript([]byte(`
        transaction { 
            prepare(signer: AuthAccount) { log(signer.address) }
        }
    `)).
    SetGasLimit(100).
    SetProposalKey(account1.Address, key1.Index, key1.SequenceNumber).
    SetPayer(account1.Address).
    AddAuthorizer(account1.Address)

// 账户1 使用 key 1 进行交易的签名
err := tx.SignEnvelope(account1.Address, key1.Index, key1Signer)
```

[完整的可运行的例子](/examples#single-party-single-signature)

---

##### [多人签名](https://github.com/onflow/flow/blob/master/docs/accounts-and-keys.md#single-party-multiple-signatures)

- Proposer, payer, authorizer是同一个account (`0x01`)

- 只要交易必须签名

- 每个密钥的权重为0.5，因此需要两个签名者。

| Account   | Key ID | Weight |
|-----------|--------|--------|
| `0x01`    | 1      | 0.5    |
| `0x01`    | 2      | 0.5    |

```go
account1, _ := c.GetAccount(ctx, flow.HexToAddress("01"))

key1 := account1.Keys[0]
key2 := account1.Keys[1]

//  安全的通过私钥创建两个签名器
key1Signer := getSignerForKey1()
key2Signer := getSignerForKey2()

tx := flow.NewTransaction().
    SetScript([]byte(`
        transaction { 
            prepare(signer: AuthAccount) { log(signer.address) }
        }
    `)).
    SetGasLimit(100).
    SetProposalKey(account1.Address, key1.Index, key1.SequenceNumber).
    SetPayer(account1.Address).
    AddAuthorizer(account1.Address)

//  账户1 使用 key 1 进行交易的签名
err := tx.SignEnvelope(account1.Address, key1.Index, key1Signer)

//  账户1 使用 key 2 进行交易的签名
err = tx.SignEnvelope(account1.Address, key2.Index, key2Signer)
```

[完整的可运行的例子](/examples#single-party-multiple-signatures)

---

##### [多方](https://github.com/onflow/flow/blob/master/docs/accounts-and-keys.md#multiple-parties)


- Proposer和authorizer是同一个账号(`0x01`)

- Payer是一个单独的帐户(`0x02`).

- 帐户`0x01`对交易负载(payload)签名。

- 帐户`0x02`对信封(envolope)签名。

- 帐户`0x02`必须最后签名，因为它是付款人。

| Account   | Key ID | Weight |
|-----------|--------|--------|
| `0x01`    | 1      | 1.0    |
| `0x02`    | 3      | 1.0    |

```go
account1, _ := c.GetAccount(ctx, flow.HexToAddress("01"))
account2, _ := c.GetAccount(ctx, flow.HexToAddress("02"))

key1 := account1.Keys[0]
key3 := account2.Keys[0]

// 安全的通过私钥创建两个签名器
key1Signer := getSignerForKey1()
key3Signer := getSignerForKey3()

tx := flow.NewTransaction().
    SetScript([]byte(`
        transaction { 
            prepare(signer: AuthAccount) { log(signer.address) }
        }
    `)).
    SetGasLimit(100).
    SetProposalKey(account1.Address, key1.Index, key1.SequenceNumber).
    SetPayer(account2.Address).
    AddAuthorizer(account1.Address)

// 账户1 使用 key 1 进行交易的签名
err := tx.SignPayload(account1.Address, key1.Index, key1Signer)

// 账户2 使用 key 3 进行交易的签名
// 注意: 付款者总是最后一个签名
err = tx.SignEnvelope(account2.Address, key3.Index, key3Signer)
```

[完整的可运行的例子](/examples#multiple-parties)

---
##### [多方参与，两个授权者(authorizer)](https://github.com/onflow/flow/blob/master/docs/accounts-and-keys.md#multiple-parties)

- Proposer和authorizer是同一个账号(`0x01`)

- Payer是一个单独的帐户(`0x02`)

- 帐户`0x01`对支付签名。

- 帐户`0x02`对信封(envolope)签名。

-帐户`0x02`必须最后签名，因为它是付款人。

-帐户`0x02`也是一个授权器，用于展示如何将两个authaccount包含到一个事务中

| Account   | Key ID | Weight |
|-----------|--------|--------|
| `0x01`    | 1      | 1.0    |
| `0x02`    | 3      | 1.0    |

```go
account1, _ := c.GetAccount(ctx, flow.HexToAddress("01"))
account2, _ := c.GetAccount(ctx, flow.HexToAddress("02"))

key1 := account1.Keys[0]
key3 := account2.Keys[0]

// 安全的通过私钥创建两个签名器
key1Signer := getSignerForKey1()
key3Signer := getSignerForKey3()

tx := flow.NewTransaction().
    SetScript([]byte(`
        transaction {
            prepare(signer1: AuthAccount, signer2: AuthAccount) {
              log(signer.address)
              log(signer2.address)
          }
        }
    `)).
    SetGasLimit(100).
    SetProposalKey(account1.Address, key1.Index, key1.SequenceNumber).
    SetPayer(account2.Address).
    AddAuthorizer(account1.Address).
    AddAuthorizer(account2.Address)

//  账户1 使用 key 1 进行交易的签名
err := tx.SignPayload(account1.Address, key1.Index, key1Signer)

//  账户2 使用 key 3 进行交易的签名
//  注意: 付款者总是最后一个签名
err = tx.SignEnvelope(account2.Address, key3.Index, key3Signer)
```

[完整的可运行的例子](/examples#multiple-parties-two-authorizers)

---

##### [多方参与，多个签名](https://github.com/onflow/flow/blob/master/docs/accounts-and-keys.md#multiple-parties-multiple-signatures)

- Proposer和authorizer是同一个账号(`0x01`)

- Payer是一个单独的帐户(`0x02`)

- 帐户`0x01`对有效付款签名。

- 帐户`0x02`对信封(envolope)签名。

- 帐户`0x02`必须最后签名，因为它是付款人。

- 两个账户都必须签名两次(每个密钥签名一次)。

| Account   | Key ID | Weight |
|-----------|--------|--------|
| `0x01`    | 1      | 0.5    |
| `0x01`    | 2      | 0.5    |
| `0x02`    | 3      | 0.5    |
| `0x02`    | 4      | 0.5    |

```go
account1, _ := c.GetAccount(ctx, flow.HexToAddress("01"))
account2, _ := c.GetAccount(ctx, flow.HexToAddress("02"))

key1 := account1.Keys[0]
key2 := account1.Keys[1]
key3 := account2.Keys[0]
key4 := account2.Keys[1]

// 安全的通过私钥创建4个签名器
key1Signer := getSignerForKey1()
key2Signer := getSignerForKey1()
key3Signer := getSignerForKey3()
key4Signer := getSignerForKey4()

tx := flow.NewTransaction().
    SetScript([]byte(`
        transaction { 
            prepare(signer: AuthAccount) { log(signer.address) }
        }
    `)).
    SetGasLimit(100).
    SetProposalKey(account1.Address, key1.Index, key1.SequenceNumber).
    SetPayer(account2.Address).
    AddAuthorizer(account1.Address)

// 账户1 使用 key 1 进行交易的签名
err := tx.SignPayload(account1.Address, key1.Index, key1Signer)

// 账户1 使用 key 2 进行交易的签名
err = tx.SignPayload(account1.Address, key2.Index, key2Signer)

// 账户2 使用 key 3 进行交易的签名
// 注意: 付款者总是最后一个签名
err = tx.SignEnvelope(account2.Address, key3.Index, key3Signer)

// 账户2 使用 key 4 进行交易的签名
// 注意: 付款者总是最后一个签名
err = tx.SignEnvelope(account2.Address, key4.Index, key4Signer)
```

[完整的可运行的例子](/examples#multiple-parties-multiple-signatures)

### 发起一笔交易

You can submit a transaction to the network using the Access API client.

```go
import "github.com/onflow/flow-go-sdk/client"

// 连接本地服务
c, err := client.New("localhost:3569")
if err != nil {
    panic("failed to connect to emulator")
}

ctx := context.Background()

err = c.SendTransaction(ctx, tx)
if err != nil {
    panic("failed to send transaction")
}
```

### 查询交易结果

提交交易后，您可以通过ID查询交易状态:

```go
result, err := c.GetTransactionResult(ctx, tx.ID())
if err != nil {
    panic("failed to fetch transaction result")
}
```

结果包括一个“Status”字段，该字段将是以下值之一:

- `UNKNOWN` - 该交易尚未被网络看到。

- `PENDING` - 交易尚未包含在一个块中。

- `FINALIZED` - 交易已包含在一个区块中。

- `EXECUTED` - 交易已执行，但结果尚未封包。

- `SEALED` - 事务已经被执行并且结果被打包在一个块中。

```go
if result.Status == flow.TransactionStatusSealed {
  fmt.Println("Transaction is sealed!")
}
```

结果还包含一个“Error”，该“Error”包含失败事务的错误信息。

```go
if result.Error != nil {
    fmt.Printf("Transaction failed with error: %v\n", result.Error)
}
```

### 查询区块

你可以使用' GetLatestBlock '方法来获取最新的打包或未打包的区块:

```go
 
isSealed := true
latestBlock, err := c.GetLatestBlock(ctx, isSealed)
if err != nil {
    panic("failed to fetch latest sealed block")
}

 
isSealed := false
latestBlock, err := c.GetLatestBlock(ctx, isSealed)
if err != nil {
    panic("failed to fetch latest unsealed block")
}
```

一个块包含以下字段:

- `ID` - 块的ID(散列)

- `ParentBlockID` - 链中前一个块的ID。

- `Height` - 链条中区块的高度。

- `collectionguarantee` - 集合中包含的集合列表。


### 执行一个脚本
可以使用“ExecuteScriptAtLatestBlock”方法根据最新的密封执行状态执行只读脚本。

此功能可用于从区块链读取状态。

脚本必须采用以下形式:

- 具有单一返回值的单一`main`函数

这是一个有效脚本的例子:

```
fun main(): Int { return 1 }
```

```go
import "github.com/onflow/cadence"

script := []byte("fun main(): Int { return 1 }")

value, err := c.ExecuteScript(ctx, script)
if err != nil {
    panic("failed to execute script")
}

ID := value.(cadence.Int)

// 转换成  Go int 类型
myID := ID.Int()
```

### 查询事件

你可以查询事件与' GetEventsForHeightRange '函数:

```go
import "github.com/onflow/flow-go-sdk/client"

blocks, err := c.GetEventsForHeightRange(ctx, client.EventRangeQuery{
    Type:       "flow.AccountCreated",
    StartHeight: 10,
    EndHeight:   15,
})
if err != nil {
    panic("failed to query events")
}
```

#### 事件查询格式

一个事件查询包括以下字段:

**Type**

要筛选的事件类型。事件类型由声明它们的帐户和合约命名。

例如, 一个 `Transfer` 被定义在一个 `Token` 合约中，该合约被部署在账户 `0x55555555555555555555` 中， 将会得到一个类型 `A.0x55555555555555555555.Token.Transfer`.

阅读 [编程语言文档](https://docs.onflow.org/cadence/language/events/) 关于 Cadence 语言.

**StartHeight, EndHeight**

要过滤的块。事件将从“StartHeight”到“EndHeight”范围内的块返回。

#### 事件的结果

函数的作用是: 返回按块分组的事件。

每个块包含一个按执行顺序匹配查询的事件列表。

```go
for _, block := range blocks {
    fmt.Printf("Events for block %s:\n", block.BlockID)
    for _, event := range block.Events {
        fmt.Printf(" - %s", event)
    }
}
```

<!--
#### Decoding an Event

TODO: 事件解码示例
-->

### 查询账户

您可以查询帐户的状态用 `GetAccount` 函数:
```go
import "github.com/onflow/flow-go-sdk"

address := flow.HexToAddress("01")

account, err := c.GetAccount(ctx, address)
if err != nil {
    panic("failed to fetch account")
}
```

A `flow.Account` contains the following fields:

- `Address: flow.Address` - The account address.
- `Balance: uint64` - The account balance.
- `Code: []byte` - The code deployed at this account.
- `Keys: []flow.AccountKey` - A list of the public keys associated with this account.

## 例子

[examples](/examples)目录包含使用SDK与控件交互的代码示例

[Flow Emulator](https://docs.onflow.org/devtools/emulator/).
