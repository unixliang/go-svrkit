package main

import(
		"net"
		"bytes"
		"encoding/binary"
		svrkit "go-svrkit"
	  )

func Logic(seq uint32, conn *net.UDPConn, req []byte, fromAddr *net.UDPAddr) {
	var a uint32
	var b uint32
	var sum uint32

		p := bytes.NewBuffer(req)
		binary.Read(p, binary.BigEndian, &seq)
		binary.Read(p, binary.BigEndian, &a)
		binary.Read(p, binary.BigEndian, &b)

		sum = a + b

		svrkit.Infof("seq: %v a: %v b: %v sum: %v", seq, a, b, sum)

		p = new(bytes.Buffer)
		binary.Write(p, binary.BigEndian, seq)
		binary.Write(p, binary.BigEndian, sum)

		_, err := conn.WriteToUDP(p.Bytes(), fromAddr)
		if err != nil {
			panic(err)
		}
}

func main() {
	var err error

	err = svrkit.RegisterLocal("127.0.0.1:3001", Logic)
	if err != nil {
		panic(err)
	}

	svrkit.Run(10, "./add_svr", 7)

}
