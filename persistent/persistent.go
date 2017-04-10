package persistent

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"os"

	mmap "github.com/edsrzf/mmap-go"
)

type Map struct {
	mem       mmap.MMap
	memData   map[interface{}]interface{}
	totalSize int
	keySize   int
	valSize   int
	elemNum   int
	index     int
	filename  string
}

func NewMap(keySize, valSize, totalSize int, filename string, isTruncate bool) (*Map, error) {
	persistentMap := &Map{}
	persistentMap.keySize = keySize
	persistentMap.valSize = valSize
	persistentMap.filename = filename
	persistentMap.elemNum = totalSize / (keySize + valSize)
	var (
		f       *os.File
		err     error
		memData map[interface{}]interface{}
		mem     mmap.MMap
	)

	_, err = os.Stat(filename)
	if !os.IsNotExist(err) { // 原有文件存在
		if isTruncate {
			bkFileName := filename + ".bak"
			os.Rename(filename, bkFileName)
			log.Printf("文件存在，做出备份")
		}
	}

	f, err = os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	if isTruncate { //
		buf := make([]byte, totalSize)
		writeSize, err := f.Write(buf)
		if err != nil {
			log.Printf("%v", err)
			return nil, err
		}

		log.Printf("write size=%v", writeSize)

		mem, err = mmap.Map(f, mmap.RDWR, 0)
		if err != nil {
			log.Printf("%v", err)
			return nil, err
		}
		memData = make(map[interface{}]interface{}, persistentMap.elemNum)
	} else {
		log.Printf("%+v \n", f)
		// mem, err = mmap.Map(f, mmap.RDWR, 0)
		mem, err = mmap.MapRegion(f, persistentMap.totalSize, mmap.RDWR, 0, 0)
		if err != nil {
			log.Printf("%+v", err)
			log.Println("%+v", err)
			return nil, err
		}
		memData = make(map[interface{}]interface{}, persistentMap.elemNum)
		index := 0
		for i := 0; i < persistentMap.elemNum; i++ { // 从文件中恢复
			var key, val *interface{}
			var k = bytes.NewBuffer(mem[index : index+persistentMap.keySize])
			dec := gob.NewDecoder(k)
			err = dec.Decode(key)
			if err != nil {
				log.Fatalf("错误 %v", err)
				return nil, err
			}

			var v = bytes.NewBuffer(mem[index+persistentMap.keySize : index+persistentMap.keySize+persistentMap.valSize])
			dec = gob.NewDecoder(v)
			err = dec.Decode(val)
			if err != nil {
				log.Fatalf("错误 %v", err)
				return nil, err
			}

			memData[key] = val
			index += (persistentMap.keySize + persistentMap.valSize)

		}
	}
	persistentMap.memData = memData
	persistentMap.mem = mem
	return persistentMap, nil
}

func (m *Map) Get(key interface{}) (interface{}, error) {
	return m.memData[key], nil
}

// todo: 返回接口 应有 告诉 是否 成功更新
func (m *Map) Set(key, val interface{}) error {
	m.memData[key] = val
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		log.Printf("%v", err)
		return err
	}

	if len(buf.Bytes()) > m.keySize {
		return fmt.Errorf("key的size过大, len(buf.Bytes())=%v, m.KeySize= %v,buf=%+v", len(buf.Bytes()), m.keySize, buf.Bytes())
	}

	if len(buf.Bytes()) > m.valSize {
		log.Printf("%v", err)
		return fmt.Errorf("value的size过大")
	}

	k := buf.Bytes()
	// m.mem[i] = buf.Bytes()
	index := m.index
	copy(m.mem[index:index+m.keySize], k[:])

	buf.Reset()
	valEnc := gob.NewEncoder(&buf)
	err = valEnc.Encode(val)
	if err != nil {
		log.Printf("%v", err)
		return fmt.Errorf("val encode error")
	}

	v := buf.Bytes()
	copy(m.mem[index+m.keySize:index+m.keySize+m.valSize], v[:])
	m.index += (m.keySize + m.valSize)
	log.Printf("m.index=%v", m.index)
	m.Flush()
	return nil
}

func (m *Map) Flush() {
	m.mem.Flush()
}

func (m *Map) Unmap() {
	m.mem.Unmap()
}

/*
func(m *Map) set2Mem(i interface{}){
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return err
	}

	if len(buf.Bytes()) > m.keySize{
		return fmt.Errorf("key的size过大")
	}

	if len(buf.Bytes()) > m.valSize {
		return fmt.Errorf("key的size过大")
	}

	k := buf.Bytes()
	// m.mem[i] = buf.Bytes()
	m.mem[index:index+keySize] = k[:]

}
*/
