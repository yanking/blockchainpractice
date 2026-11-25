// RSA 练习：生成密钥对，用私钥对满足 PoW（哈希值 4 个 0 开头）的「昵称 + nonce」签名，再用公钥验证。
package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

// 要签名的「昵称 + nonce」需满足的 PoW 难度：哈希值 4 个 0 开头。
const difficulty = 4

func main() {
	name := flag.String("name", "", "昵称")
	nonce := flag.Uint64("nonce", 0, "起始 nonce")
	flag.Parse()

	if *name == "" {
		fmt.Fprintln(os.Stderr, "请通过 -name 指定昵称，例如: go run . -name wangyan -nonce 0")
		os.Exit(2)
	}
	fmt.Printf("昵称: %s  起始 nonce: %d\n\n", *name, *nonce)

	// 1. 生成公私钥对
	priv, err := GenerateKey()
	if err != nil {
		fail(err)
	}
	pubPEM, err := PublicPEM(&priv.PublicKey)
	if err != nil {
		fail(err)
	}
	privPEM, err := PrivatePEM(priv)
	if err != nil {
		fail(err)
	}
	fmt.Printf("[1] 生成 RSA 密钥对（%d 位）\n", keyBits)
	fmt.Printf("公钥:\n%s", pubPEM)
	fmt.Printf("私钥:\n%s", privPEM)
	fmt.Println()

	// 2. PoW：找出要签名的内容
	res := Mine(*name, *nonce, difficulty)
	fmt.Printf("[2] PoW：寻找 %d 个 0 开头的 sha256(昵称 + nonce)\n", difficulty)
	fmt.Printf("内容: %s\n", res.Content)
	fmt.Printf("哈希: %x\n\n", res.Hash)

	// 3. 用私钥签名
	sig, err := Sign(priv, res.Content)
	if err != nil {
		fail(err)
	}
	fmt.Println("[3] 用私钥签名（SHA-256 + RSA PKCS#1 v1.5）")
	fmt.Printf("签名: %x\n\n", sig)

	// 4. 用公钥验证
	fmt.Println("[4] 用公钥验证")
	fmt.Printf("原内容 %s: %s\n", res.Content, verdict(Verify(&priv.PublicKey, res.Content, sig)))

	// 把 nonce 改掉，再用同一个签名验证，应当不通过
	tampered := *name + strconv.FormatUint(res.Nonce+1, 10)
	fmt.Printf("内容改成 %s: %s\n", tampered, verdict(Verify(&priv.PublicKey, tampered, sig)))
}

// verdict 把验证结果转成一句话。
func verdict(err error) string {
	if err != nil {
		return "验证不通过"
	}
	return "验证通过"
}

// fail 打印错误并以状态码 2 退出。
func fail(err error) {
	fmt.Fprintln(os.Stderr, "错误:", err)
	os.Exit(2)
}
