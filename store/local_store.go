// 本地存储
// 可以在 service 和 cmd 中使用
// 在 service 中可以将 file 写入磁盘，路径是 `root`/`fileHash`/`fileHash`
// 在 cmd 中可以使用 MergeChunksToSomewhere 将 chunks 合并放入需要的文件中
package store

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/TensShinet/WeFile/utils"
	"hash"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"syscall"
)

type LocalStore struct {
	root              string
	tempFilesPath     string
	tempChunksPath    string
	samplingChunkSize int
}

type LocalFile struct {
	*os.File
	path         string // 存储位置
	totalHash    string // 使用 sha256 计算 hash
	samplingHash string // 计算抽样 hash
	currentSize  int64
	maxSize      int64 // 最大大小限制
	hashing      hash.Hash
}

type LocalChunk struct {
	*os.File
	path        string
	currentSize int
	maxSize     int
}

const (
	defaultLocalStoreRoot = "/tmp/wefile"
	defaultFileNameLength = 12
)

var (
	ErrFileSizeExceed  = errors.New("file size exceed")
	ErrChunkSizeExceed = errors.New("chunk size exceed")
)

func NewLocalStore(root string, samplingChunkSize int) (Store, error) {

	if root == "" {
		root = defaultLocalStoreRoot
	}

	if err := os.MkdirAll(root, 0755); err != nil {
		return nil, err
	}
	// temp 存储目录
	tempPath := filepath.Join(root, "temp")
	s := &LocalStore{
		root:              root,
		tempFilesPath:     filepath.Join(tempPath, "files"),
		tempChunksPath:    filepath.Join(tempPath, "chunks"),
		samplingChunkSize: samplingChunkSize,
	}
	// 创建文件暂存目录
	if err := os.MkdirAll(s.tempFilesPath, 0755); err != nil {
		return nil, err
	}

	// 创建 chunks 暂存目录
	if err := os.MkdirAll(s.tempChunksPath, 0755); err != nil {
		return nil, err
	}

	return s, nil
}

// 创建一个 file path 是暂存路径 limit 是文件限制如果是读的话暂时不需要这些限制
func (s *LocalStore) TempFile(fileLimit FileLimit) (File, error) {

	var (
		err  error
		path string
		file *os.File
	)

	path = filepath.Join(s.tempFilesPath, utils.RandomString(defaultFileNameLength))
	if file, err = os.Create(path); err != nil {
		return nil, err
	}
	f := &LocalFile{
		File:    file,
		path:    path,
		maxSize: fileLimit.MaxSize,
		hashing: sha256.New(),
	}
	return f, nil
}

// identifier 确定一组 chunk 存的位置 chunkIndex 确定 chunk 的索引
func (s *LocalStore) TempChunk(identifier string, chunkIndex int, chunkLimit ChunkLimit) (Chunk, error) {

	var (
		err  error
		file *os.File
	)
	c := &LocalChunk{
		path:    filepath.Join(s.tempChunksPath, identifier, utils.ParseIntToString(chunkIndex)),
		maxSize: chunkLimit.MaxSize,
	}

	if file, err = utils.CreateFile(c.path); err != nil {
		if err == syscall.EEXIST {
			err = ErrChunkExists
		}
		return nil, err
	}

	c.File = file

	return c, nil
}

func (s *LocalStore) GetChunksPath(identifier string) string {
	return filepath.Join(s.tempChunksPath, identifier)
}

// 将一个 file 写入存储中
func (s *LocalStore) SaveFile(file File) (File, error) {
	var (
		f      *LocalFile
		err    error
		path   string
		osFile *os.File
	)
	path = filepath.Join(s.root, file.TotalHash(), file.TotalHash())
	if osFile, err = utils.CreateFile(path); err != nil {
		if err != syscall.EEXIST {
			return nil, err
		} else {
			// 已经保存的 file 不需要再写
			osFile = nil
		}

	}
	buf := make([]byte, s.samplingChunkSize)
	hashing := sha256.New()
	for {
		n, err := file.Read(buf)
		if err != io.EOF && err != nil {
			return nil, err
		}
		if err == io.EOF {
			break
		} else if n == s.samplingChunkSize {
			hashing.Write(buf[:10])
			hashing.Write(buf[n-10 : n])
		} else {
			hashing.Write(buf[:n])
		}
		if osFile != nil {
			if n, err = osFile.Write(buf[:n]); err != nil {
				return nil, err
			}
		}
	}
	if osFile != nil {
		osFile.Close()
	}

	f = &LocalFile{
		path:         path,
		totalHash:    file.TotalHash(),
		samplingHash: hex.EncodeToString(hashing.Sum(nil)),
		currentSize:  file.Size(),
	}

	return f, nil
}

// 调用这个之前需要检查 chunksNum 是否满足
func (s *LocalStore) MergeChunksToSave(chunksPath string, chunksNum int, fileLimit FileLimit) (File, error) {

	var (
		err      error
		reader   *chunksReader
		tempFile File
	)

	reader = newChunksReader(chunksPath, chunksNum)
	tempFile, err = s.TempFile(fileLimit)
	defer tempFile.Remove()
	if _, err = io.Copy(tempFile, reader); err != nil {
		return nil, err
	}
	if _, err = tempFile.Seek(0, 0); err != nil {
		return nil, err
	}

	return s.SaveFile(tempFile)
}

func (s *LocalStore) MergeChunksToSomewhere(path string, chunksPath string, chunksNum int, fileLimit FileLimit) (File, error) {
	return nil, nil
}

// File 相关
func (f *LocalFile) Size() int64 {
	return f.currentSize
}

func (f *LocalFile) Remove() error {
	_ = f.Close()
	return os.Remove(f.path)
}

func (f *LocalFile) Location() string {
	return f.path
}

func (f *LocalFile) TotalHash() string {
	if f.totalHash == "" {
		f.totalHash = hex.EncodeToString(f.hashing.Sum(nil))
	}
	return f.totalHash
}

func (f *LocalFile) SamplingHash() string {
	return f.samplingHash
}

// 控制文件写 计算全量 hash
func (f *LocalFile) Write(p []byte) (n int, err error) {
	n = len(p)
	if int64(n)+f.currentSize > f.maxSize {
		_ = f.Remove()
		return 0, ErrFileSizeExceed
	}
	f.currentSize += int64(n)
	f.hashing.Write(p)
	return f.File.Write(p)
}

// Chunk 相关
func (c *LocalChunk) Remove() error {
	_ = c.Close()
	return os.Remove(c.path)
}

func (c *LocalChunk) RemoveAll() error {
	_ = c.Close()
	return os.RemoveAll(filepath.Dir(c.path))
}

// 控制 chunk 写
func (c *LocalChunk) Write(p []byte) (n int, err error) {
	n = len(p)
	if n+c.currentSize > c.maxSize {
		_ = c.Remove()
		return 0, ErrChunkSizeExceed
	}
	c.currentSize += n
	return c.File.Write(p)
}

// Chunks Reader 相关

const (
	defaultChunksReaderLen = 1024 * 1024
)

type chunksReader struct {
	path              string   // 存放所有 chunk 的目录
	num               int      // chunk 数量
	currentChunkIndex int      // 当前 chunk 索引 默认从 1 开始
	buf               []byte   // chunk 的缓存
	bufIndex          int      // buf 的指针
	bufLen            int      // buf 实际长度
	file              *os.File //
}

func newChunksReader(path string, num int) *chunksReader {
	return &chunksReader{
		path:              path,
		num:               num,
		currentChunkIndex: 1,
		buf:               make([]byte, defaultChunksReaderLen),
		bufIndex:          0,
		bufLen:            0,
	}
}

func (r *chunksReader) Read(p []byte) (n int, err error) {
	needLen, readLen, restLen := len(p), 0, r.bufLen-r.bufIndex
	if restLen >= needLen {
		copy(p, r.buf[r.bufIndex:r.bufIndex+needLen])
		r.bufIndex += needLen
		return needLen, nil
	}

	copy(p[readLen:readLen+restLen], r.buf[r.bufIndex:r.bufIndex+restLen])
	readLen += restLen

	r.bufIndex, r.bufLen = 0, 0

	if r.file == nil {
		if r.file, err = os.Open(filepath.Join(r.path, strconv.Itoa(r.currentChunkIndex))); err != nil {
			return 0, err
		}
	}
	for r.currentChunkIndex <= r.num {
		n, err = r.file.Read(r.buf)
		// 什么都没有读到
		if err == io.EOF {
			r.file.Close()
			if r.currentChunkIndex == r.num {
				break
			}
			r.currentChunkIndex += 1
			if r.file, err = os.Open(filepath.Join(r.path, strconv.Itoa(r.currentChunkIndex))); err != nil {
				return 0, err
			}
			continue
			//	读错误
		} else if err != nil {
			return 0, err
		}
		restLen, r.bufLen = n, n
		// 新读出来的长度大于需要的长度
		if restLen >= needLen-readLen {
			copy(p[readLen:needLen], r.buf[:needLen-readLen])
			r.bufIndex = needLen - readLen
			return needLen, nil
		}
		copy(p[readLen:readLen+restLen], r.buf[:restLen])
		readLen += restLen
		r.bufIndex, r.bufLen = 0, 0
	}

	return readLen, io.EOF
}
