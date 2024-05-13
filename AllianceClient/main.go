package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	ip := "127.0.0.1"
	port := 12345
	addr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		log.Fatalln("resolve addr error: ", err)
	}

	conn, err := net.DialTCP("tcp4", nil, addr)
	if err != nil {
		log.Fatalln("dial failed: ", err)
	}

	var input string
	var length int
	reader := bufio.NewReader(os.Stdin)
	go func() {
		for {
			input, _ = reader.ReadString('\n')
			input = strings.TrimSpace(input)
			length = len(input)
			buf := new(bytes.Buffer)
			err := binary.Write(buf, binary.LittleEndian, int32(length))
			if err != nil {
				fmt.Println("binary.Write failed:", err)
			}
			_, err = conn.Write(append(buf.Bytes(), []byte(input)...))
		}
	}()

	for {
		length := make([]byte, 4) // 长度的字节数固定为4
		if _, err := io.ReadFull(conn, length); err != nil {
			return
		}
		realLength := binary.LittleEndian.Uint32(length)
		data := make([]byte, realLength)
		if _, err := io.ReadFull(conn, data); err != nil {
			return
		}
		content := string(data)
		fmt.Println("===================receive message:===================")
		fmt.Println(content)
	}

}
