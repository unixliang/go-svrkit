package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)


func RunLogic(req []byte, fromAddr *net.UDPAddr) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Print(err)
		}
	}()

	var A, B uint32

	p := bytes.NewBuffer(req)
	binary.Read(p, binary.BigEndian, &A)
	binary.Read(p, binary.BigEndian, &B)

	var addIn AddIn

	addIn.A = A
	addIn.B = B

	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:4000")

	res, err := RunService("Add", addr, &addIn, 1000)
	if err != nil {
		panic(fmt.Errorf("Add err: %v", err))
	}

	addOut := res.(*AddOut)

	p = new(bytes.Buffer)
	binary.Write(p, binary.BigEndian, addOut.Sum)

	rsp := p.Bytes()

	_, err = listenConn.WriteToUDP(rsp, fromAddr)
	if err != nil {
		panic(err)
	}
}

func main() {

	err := Init("127.0.0.1", 3000)
	if err != nil {
		panic(err)
	}

	err = RegisterService("Add")
	if err != nil {
		panic(err)
	}

	for {
		bytes := make([]byte, MAX_PKG_LEN)
		n, addr, err := listenConn.ReadFromUDP(bytes)
		if err != nil {
			panic(err)
		}

		req := bytes[0:n]
		go RunLogic(req, addr)
	}

	return

}
