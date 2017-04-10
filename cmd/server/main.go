package main

import (
	"encoding/gob"
	"flag"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/yafngzh/easyKV/kvclnt"
	"github.com/yafngzh/easyKV/kvserver"
	"github.com/yafngzh/easyKV/msg"
)

var (
	svr *kvserver.KVServer
)

func handleConn(c net.Conn) {
	var exitClose = true
	defer func() {
		if exitClose {
			log.Printf("c exit,RemoteAddr %+v", c.RemoteAddr)
			c.Close()
		}
	}()
	for {
		var (
			clnt kvclnt.KVClnt
		)
		t1 := time.Now().UnixNano()
		dec := gob.NewDecoder(c)
		curMsg := &msg.Msg{}
		err := dec.Decode(curMsg)
		if err != nil {
			if err == io.EOF {
				log.Printf("[read %v] 对方关闭连接\n", c.RemoteAddr())
			} else {
				log.Printf("read error %v from %v, now exit\n", err, c.RemoteAddr())
			}
			break
		}
		//log.Printf("receive msg %+v\n", curMsg)
		clnt.CreatTime = curMsg.CreatTime
		addr := c.RemoteAddr().String()
		addrSlice := strings.Split(addr, ":")
		if len(addrSlice) != 2 {
			log.Printf("格式错误 %s \n", addr)
			break
		}
		clnt.IP = addrSlice[0]
		clnt.Port, err = strconv.Atoi(addrSlice[1])
		if err != nil {
			log.Printf("%v\n", err)
			break
		}
		clnt.PID = curMsg.PID
		clnt.Conn = c

		resp, err := svr.HandleMsg(curMsg, clnt)
		if err != nil {
			log.Fatalf("handle failed %v", err)
			break
		}
		encoder := gob.NewEncoder(c)
		err = encoder.Encode(resp)
		if err != nil {
			log.Println("send error %v", err)
		}

		t2 := time.Now().UnixNano()
		log.Printf("time elapsed %d ms\n", (t2-t1)/1000)
		exitClose = false
		break
	}
}
func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	keySize := flag.Int("keysize", 8, "")
	valSize := flag.Int("valsize", 8, "")
	totalSize := flag.Int("totalsize", 8, "总共支持的数量")
	filename := flag.String("filename", "/data/appdata/easyKV/persistentmap", "")

	flag.Parse()

	svr = kvserver.NewKVServer(keySize, valSize, totalSize, filename)
	defer {

	}
	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Printf("Listen error %v", err)
		return
	}

	for {
		c, err := l.Accept()
		if err != nil {
			log.Printf("Accept error, %v", err)
			continue
		}
		go handleConn(c)
	}
}
