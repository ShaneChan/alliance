package main

import (
	_ "AllianceServer/mgo"
	"AllianceServer/predefine"
	"AllianceServer/user"
	"log"
	"net"
)

func main() {
	ipAddr, err := net.ResolveTCPAddr("tcp4", predefine.Cfg["ipaddress"])

	if err != nil {
		log.Fatal("resolve tcp address error: ", err)
	}

	listener, err := net.ListenTCP("tcp4", ipAddr)
	if err != nil {
		log.Fatalln("listener error: ", err)
	}

	log.Println("starting listening...")
	log.Println("listen ok on", ipAddr)

	// 主协程监听，来了新连接分发到新的协程去处理
	for {
		// 开启监听
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Fatalln("accept error: ", err)
		}
		log.Println("accept new connection, remote address: ", conn.RemoteAddr().String())

		// 请求路由到新的协程去处理
		go user.NewConnection(conn).DealConnection()
	}
}
