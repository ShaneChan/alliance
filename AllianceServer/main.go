package main

import (
	_ "AllianceServer/mgo"
	"AllianceServer/user"
	"fmt"
	"log"
	"net"
)

func main() {
	ip := "127.0.0.1"
	port := 12345
	addr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		log.Fatal("resolve tcp address error: ", err)
	}

	listener, err := net.ListenTCP("tcp4", addr)
	if err != nil {
		log.Fatalln("listener error: ", err)
	}

	log.Println("listen ok on ", addr)
	for true { // 主协程监听，来了新连接分发到新的协程去处理
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Fatalln("accept error: ", err)
		}
		log.Println("accept new connection, remote address: ", conn.RemoteAddr().String())

		go user.NewConnection(conn).DealConnection()
	}
}
