package kvclnt

import "net"

type KVClnt struct {
	CreatTime uint64
	IP        string
	Port      int
	PID       int
	Conn      net.Conn
}

func (clnt *KVClnt) Equal(clnt_another *KVClnt) bool {
	if clnt.IP == clnt_another.IP && clnt.Port == clnt_another.Port && clnt.CreatTime == clnt_another.CreatTime {
		return true
	}
	return false
}
