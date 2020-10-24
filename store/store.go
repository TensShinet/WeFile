// 存储接口设计
package store

import "io"

type File interface {
	io.ReadWriteCloser
	Remove() error
	Location() string
	TotalHash() string
	SamplingHash() string
	Size() int64
}

type Chunk interface {
	io.ReadWriteCloser
	Remove() error
	RemoveAll() error
}

type FileLimit struct {
	MaxSize int64 // 最大大小限制
}

type ChunkLimit struct {
	MaxSize int
}

// Store 接口设计
// TempFile TempChunk 都是暂存文件 通常是本地存储
// SaveFile 将一个暂存的文件存入 store 中，可以是 oss 也可以是 ceph 等
// MergeChunksToSave 将多个 chunk 合并到 store 中
// MergeChunksToSomewhere 将多个 chunk 合并到其他的地方
// GetChunksPath 获取 MergeChunks 需要的参数
type Store interface {
	TempFile(fileLimit FileLimit) (File, error)
	TempChunk(identifier string, chunkIndex int, chunkLimit ChunkLimit) (Chunk, error)
	GetChunksPath(identifier string) string
	SaveFile(file File) (File, error)
	MergeChunksToSave(chunksPath string, chunkNum int, fileLimit FileLimit) (File, error)
	MergeChunksToSomewhere(path string, chunksPath string, chunkNum int, fileLimit FileLimit) (File, error)
}
