package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"testing"
)

const signed = "wangyan37063"

// newKey 生成一把测试用密钥，失败时终止测试。
func newKey(t *testing.T) *rsa.PrivateKey {
	t.Helper()
	priv, err := GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	return priv
}

func TestSignAndVerify(t *testing.T) {
	priv := newKey(t)

	sig, err := Sign(priv, signed)
	if err != nil {
		t.Fatal(err)
	}
	if len(sig) != keyBits/8 {
		t.Errorf("签名长度 = %d 字节, want %d", len(sig), keyBits/8)
	}
	if err := Verify(&priv.PublicKey, signed, sig); err != nil {
		t.Errorf("原内容验证失败: %v", err)
	}
}

func TestVerifyRejects(t *testing.T) {
	priv := newKey(t)
	sig, err := Sign(priv, signed)
	if err != nil {
		t.Fatal(err)
	}
	other := newKey(t)
	// 签名被改动一个比特
	corrupted := append([]byte(nil), sig...)
	corrupted[0] ^= 1

	tests := []struct {
		name    string
		pub     *rsa.PublicKey
		content string
		sig     []byte
	}{
		{"内容被修改", &priv.PublicKey, "wangyan37064", sig},
		{"换一把公钥", &other.PublicKey, signed, sig},
		{"签名被修改", &priv.PublicKey, signed, corrupted},
	}
	for _, tt := range tests {
		if err := Verify(tt.pub, tt.content, tt.sig); err == nil {
			t.Errorf("%s: 验证通过了，应当不通过", tt.name)
		}
	}
}

func TestPEMRoundTrip(t *testing.T) {
	priv := newKey(t)

	privPEM, err := PrivatePEM(priv)
	if err != nil {
		t.Fatal(err)
	}
	block, _ := pem.Decode([]byte(privPEM))
	if block == nil || block.Type != "PRIVATE KEY" {
		t.Fatalf("私钥 PEM 无法解析: %q", privPEM)
	}
	parsed, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		t.Fatal(err)
	}
	if key, ok := parsed.(*rsa.PrivateKey); !ok || !key.Equal(priv) {
		t.Error("私钥 PEM 解析后与原密钥不一致")
	}

	pubPEM, err := PublicPEM(&priv.PublicKey)
	if err != nil {
		t.Fatal(err)
	}
	block, _ = pem.Decode([]byte(pubPEM))
	if block == nil || block.Type != "PUBLIC KEY" {
		t.Fatalf("公钥 PEM 无法解析: %q", pubPEM)
	}
	parsedPub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		t.Fatal(err)
	}
	if pub, ok := parsedPub.(*rsa.PublicKey); !ok || !pub.Equal(&priv.PublicKey) {
		t.Error("公钥 PEM 解析后与原密钥不一致")
	}
}
