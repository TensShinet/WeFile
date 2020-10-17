package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
)

func Digest256(data []byte) string {
	s := sha256.New()
	s.Write(data)
	return hex.EncodeToString(s.Sum(nil))
}

// 将一个流的内容读完做 sha256 的 hash
func StreamDigest256(src io.Reader) (string, error) {
	t := sha256.New()
	if _, err := io.Copy(t, src); err != nil {
		return "", err
	}
	return hex.EncodeToString(t.Sum(nil)), nil
}
