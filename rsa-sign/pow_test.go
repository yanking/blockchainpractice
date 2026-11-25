package main

import (
	"crypto/sha256"
	"strconv"
	"testing"
)

func TestHasLeadingZeroHex(t *testing.T) {
	tests := []struct {
		name string
		hash []byte
		n    int
		want bool
	}{
		{"000010 满足 4 个 0", []byte{0x00, 0x00, 0x10}, 4, true},
		{"000010 第 5 位非 0", []byte{0x00, 0x00, 0x10}, 5, false},
		{"000001 满足 5 个 0", []byte{0x00, 0x00, 0x01}, 5, true},
		{"000001 不满足 6 个 0", []byte{0x00, 0x00, 0x01}, 6, false},
		{"000f 满足 3 个 0", []byte{0x00, 0x0f}, 3, true},
		{"000f 不满足 4 个 0", []byte{0x00, 0x0f}, 4, false},
		{"首位为 1", []byte{0x10}, 1, false},
	}
	for _, tt := range tests {
		if got := hasLeadingZeroHex(tt.hash, tt.n); got != tt.want {
			t.Errorf("hasLeadingZeroHex(%x, %d) = %v, want %v", tt.hash, tt.n, got, tt.want)
		}
	}
}

func TestMine(t *testing.T) {
	const name = "satoshi"
	const start uint64 = 7
	res := Mine(name, start, 4)

	if !hasLeadingZeroHex(res.Hash[:], 4) {
		t.Errorf("hash %x 不满足 4 个 0 开头", res.Hash)
	}
	if want := sha256.Sum256([]byte(name + strconv.FormatUint(res.Nonce, 10))); res.Hash != want {
		t.Errorf("hash 与 sha256(%s+%d) 不一致", name, res.Nonce)
	}
	if res.Content != name+strconv.FormatUint(res.Nonce, 10) {
		t.Errorf("Content = %q, nonce = %d", res.Content, res.Nonce)
	}

	// 结果必须是起点之后第一个满足目标的 nonce。
	for nonce := start; nonce < res.Nonce; nonce++ {
		sum := sha256.Sum256([]byte(name + strconv.FormatUint(nonce, 10)))
		if hasLeadingZeroHex(sum[:], 4) {
			t.Fatalf("nonce %d 已满足目标，Mine 却返回了 %d", nonce, res.Nonce)
		}
	}
	if want := res.Nonce - start + 1; res.Attempts != want {
		t.Errorf("Attempts = %d, want %d", res.Attempts, want)
	}
}
