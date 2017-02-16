package kvserver

import (
	"fmt"
	"sync"

	"github.com/yafngzh/easyKV/kvclnt"
	"github.com/yafngzh/easyKV/msg"
)

type KVServer struct {
	Lock  sync.RWMutex
	KVMap map[interface{}]interface{}
}

func NewKVServer() *KVServer {
	s := &KVServer{}
	s.KVMap = make(map[interface{}]interface{})
	return s
}

func (s *KVServer) HandleMsg(recvMsg interface{}, kvclnt kvclnt.KVClnt) (interface{}, error) {
	if decodeMsg, ok := recvMsg.(*msg.Msg); ok {
		if decodeMsg.Type == msg.MSG_TYPE_GET {
			return s.Get(decodeMsg.Key), nil
		} else if decodeMsg.Type == msg.MSG_TYPE_SET {
			return s.Set(decodeMsg.Key, decodeMsg.Val), nil
		} else {
			return nil, fmt.Errorf("类型错误 ", decodeMsg.Type)
		}
	} else {
		return nil, fmt.Errorf("消息结构错误 ")
	}
}
func (s *KVServer) Get(k interface{}) interface{} {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	if val, ok := s.KVMap[k]; ok != false {
		return val
	}
	return nil
}

func (s *KVServer) Set(k interface{}, v interface{}) bool {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	if val, ok := s.KVMap[k]; !ok {
		s.KVMap[k] = v
	} else if val != v {
		s.KVMap[k] = v
	} else {
		return false
	}
	return true
}
