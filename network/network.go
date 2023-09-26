package network

import (
	"encoding/binary"
	"io"
	"log"
	"net"
	"tuuna/monitor/model"
)

func SendPacket(conn net.Conn, b []byte) {
	_, err := conn.Write(b)
	if err != nil {
		panic(err)
	}
}

func RecvPacket(conn net.Conn, ch chan *model.Response) {
	for {
		var offset uint32 = 0
		b := make([]byte, 4096)

		n, err := conn.Read(b)
		if err != nil {
			if err == io.EOF {
				log.Println("disconnected from server")
				return
			}
			log.Fatal("read Error: " + err.Error())
			return
		}

		if n <= 0 {
			return
		}

		if offset < uint32(n) {
			res := parseResponse(b, &offset)
			ch <- res
		}
	}
}

func parseResponse(b []byte, offset *uint32) *model.Response {
	p := binary.LittleEndian.Uint32(b[*offset:4])
	l := binary.LittleEndian.Uint32(b[*offset+4 : 8])

	res := &model.Response{Protocol: p, Length: l, Data: b[*offset+8 : 8+l]}

	*offset += 4 + 4 + l
	return res
}

func StartConnect() net.Conn {
	conn, err := net.Dial("tcp", "127.0.0.1:9988")
	if err != nil {
		log.Fatal("dial Error: " + err.Error())
		return nil
	}

	return conn
}
