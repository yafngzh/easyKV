package persistent

import (
	"log"
	"testing"
)

func TestReadWrite(t *testing.T) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	keySize, valSize := 16, 16
	totalSize := 10240
	filename := "/tmp/testMap"
	isTruncate := true
	persistentMap, err := NewMap(keySize, valSize, totalSize, filename, isTruncate)
	if err != nil {
		t.Errorf("创建失败 %v", err)
		return
	}

	key := "hello111"
	value := "world"
	err = persistentMap.Set(key, value)

	if err != nil {
		t.Errorf("写入失败 %v", err)
		return
	}

	// t.Error("erro")
	readValue, err := persistentMap.Get(key)
	if err != nil {
		t.Errorf("读取失败 %v", err)
		return
	}

	t.Logf("readValue is %v", readValue)
	var (
		valStr string
		ok     bool
	)
	if valStr, ok = readValue.(string); !ok {
		t.Errorf("值类型错误")
		return
	}

	if valStr != value {
		t.Errorf("值不一致 before=%v after=%v", value, valStr)
		return
	}

	t.Logf("成功!")

}
