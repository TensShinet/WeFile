package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

func Digest256(data []byte) string {
	s := sha256.New()
	s.Write(data)
	return hex.EncodeToString(s.Sum(nil))
}

func ParseInt(s string) (res int, err error) {
	res, err = strconv.Atoi(s)
	return
}

func ParseInt64(s string) (res int64, err error) {
	res, err = strconv.ParseInt(s, 10, 64)
	return
}

func ParseInt64ToString(s int64) string {
	return strconv.FormatInt(s, 16)
}

func ParseIntToString(s int) string {
	return strconv.Itoa(s)
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
