// PoW 核心逻辑：穷举 nonce，直到哈希值满足难度目标。

package main

import (
	"crypto/sha256"
	"strconv"
)

// Result 是一次挖矿的结果。
type Result struct {
	Nonce    uint64   // 命中目标的 nonce
	Hash     [32]byte // sha256(昵称 + nonce)
	Content  string   // 被哈希的内容：昵称 + nonce
	Attempts uint64   // 从起始 nonce 到命中，共尝试的次数
}

// Mine 从 start 开始逐次把 nonce 加 1，直到 sha256(name+nonce) 的十六进制表示
// 以 difficulty 个 0 开头，返回该结果。
func Mine(name string, start uint64, difficulty int) Result {
	for nonce := start; ; nonce++ {
		content := name + strconv.FormatUint(nonce, 10)
		sum := sha256.Sum256([]byte(content))
		if hasLeadingZeroHex(sum[:], difficulty) {
			return Result{
				Nonce:    nonce,
				Hash:     sum,
				Content:  content,
				Attempts: nonce - start + 1,
			}
		}
	}
}

// hasLeadingZeroHex 判断哈希的十六进制表示是否以 n 个 0 开头。
// 每 2 个 0 对应一个完整的 0 字节；n 为奇数时再看下一个字节的高 4 位。
func hasLeadingZeroHex(hash []byte, n int) bool {
	for i := 0; i < n/2; i++ {
		if hash[i] != 0 {
			return false
		}
	}
	return n%2 == 0 || hash[n/2]>>4 == 0
}
