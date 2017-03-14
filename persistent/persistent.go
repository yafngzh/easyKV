package persistent

import "github.com/golang/exp/mmap"

type Map struct {
	reader mmap.ReaderAt
}

func (m *Map) NewMap(keySize, valSize int) {

}
