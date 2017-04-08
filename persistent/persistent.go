package persistent

import (
	"bytes"
	"encoding/gob"
	"os"

	"github.com/edsrzf/mmap-go"
)

type Map struct {
	mem       mmap.MMap
	memData   map[interface{}]interface{}
	totalSize int
	keySize   int
	valSize   int
	index     int
}

func NewMap(keySize, valSize, totalSize int, filename string) (MMap, error) {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)

	if err != nil {
		return nil, err
	}
	buf := make([]byte, keySize+valSize, totalSize)
	writeSize, err := f.Write(buf)
	if err != nil {
		return nil, err
	}

	mem, err := mmap.Map(f,  0, 0)
	if err != nil {
		return nil, err
	}
	memData = make(map[interface{}]interface{}, totalSize/(keySize+valSize))
	return mem, err
}

func (m *Map) Get(key interface{}) (interface{}, error) {

}

// todo: 返回接口 应有 告诉 是否 成功更新
func (m *Map) Set(key, val interface{}) error {
	memData[key] = val
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return err
	}
	for i := index; i < m.keySize; ++i {
		m.mem[i] = buf.Bytes()
	}
}
