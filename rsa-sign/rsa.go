// RSA 密钥生成、签名与验证。

package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
)

// keyBits 是 RSA 密钥长度，2048 位是目前通用的标准。
const keyBits = 2048

// GenerateKey 生成一对新的 RSA 公私钥。
func GenerateKey() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, keyBits)
}

// Sign 用私钥对 content 签名：先算 SHA-256，再做 RSA PKCS#1 v1.5 签名。
func Sign(priv *rsa.PrivateKey, content string) ([]byte, error) {
	digest := sha256.Sum256([]byte(content))
	return rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA256, digest[:])
}

// Verify 用公钥检查 content 和签名 sig 是否匹配，匹配时返回 nil。
func Verify(pub *rsa.PublicKey, content string, sig []byte) error {
	digest := sha256.Sum256([]byte(content))
	return rsa.VerifyPKCS1v15(pub, crypto.SHA256, digest[:], sig)
}

// PublicPEM 返回公钥的 PEM 文本（PKIX 格式）。
func PublicPEM(pub *rsa.PublicKey) (string, error) {
	der, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return "", err
	}
	return string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: der})), nil
}

// PrivatePEM 返回私钥的 PEM 文本（PKCS#8 格式）。
func PrivatePEM(priv *rsa.PrivateKey) (string, error) {
	der, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return "", err
	}
	return string(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der})), nil
}
