package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
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

func RandomString(n int) string {
	randBytes := make([]byte, n/2)
	_, _ = rand.Read(randBytes)
	return fmt.Sprintf("%x", randBytes)
}

// 通过 os.O_EXCL 原子创建一个 file
func CreateFile(path string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, err
	}
	if file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0644); err != nil {
		return nil, err
	} else {
		return file, nil
	}
}
