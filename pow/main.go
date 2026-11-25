// PoW 练习：不断递增 nonce，寻找十六进制表示以 N 个 0 开头的 sha256(昵称 + nonce)。
package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	name := flag.String("name", "", "昵称")
	nonce := flag.Uint64("nonce", 0, "起始 nonce")
	difficulty := flag.String("difficulty", "4,5", "难度目标：哈希值开头 0 的个数（1-64），多个用逗号分隔，依次完成")
	flag.Parse()

	if *name == "" {
		fmt.Fprintln(os.Stderr, "请通过 -name 指定昵称，例如: go run . -name wangyan -nonce 0")
		os.Exit(2)
	}
	difficulties, err := parseDifficulties(*difficulty)
	if err != nil {
		fmt.Fprintf(os.Stderr, "参数 -difficulty 无效: %v\n", err)
		os.Exit(2)
	}

	fmt.Printf("昵称: %s  起始 nonce: %d\n\n", *name, *nonce)

	// 后一个目标从前一个结果之后继续递增 nonce。
	start := *nonce
	for i, target := range difficulties {
		begin := time.Now()
		res := Mine(*name, start, target)
		elapsed := time.Since(begin)

		fmt.Printf("目标: %d 个 0 开头\n", target)
		fmt.Printf("耗时: %v（尝试 %d 次）\n", elapsed, res.Attempts)
		fmt.Printf("内容: %s\n", res.Content)
		fmt.Printf("哈希: %x\n", res.Hash)
		if i < len(difficulties)-1 {
			fmt.Println()
		}

		start = res.Nonce + 1
	}
}

// parseDifficulties 解析逗号分隔的难度列表，如 "4,5"。
func parseDifficulties(s string) ([]int, error) {
	var difficulties []int
	for _, field := range strings.Split(s, ",") {
		n, err := strconv.Atoi(strings.TrimSpace(field))
		if err != nil {
			return nil, fmt.Errorf("%q 不是整数", field)
		}
		if n < 1 || n > 64 {
			return nil, fmt.Errorf("%d 超出范围，应为 1 到 64", n)
		}
		difficulties = append(difficulties, n)
	}
	return difficulties, nil
}
