package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	svr "go_udp_svrkit"

)


func Logic(conn *net.UDPConn, req []byte, fromAddr *net.UDPAddr) {
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

	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:3001")

	res, err := svr.RunRemote("Add", addr, &addIn, 1000)
	if err != nil {
		panic(fmt.Errorf("Add err: %v", err))
	}

	addOut := res.(*AddOut)

	p = new(bytes.Buffer)
	binary.Write(p, binary.BigEndian, addOut.Sum)

	rsp := p.Bytes()

	_, err = conn.WriteToUDP(rsp, fromAddr)
	if err != nil {
		panic(err)
	}
}

func main() {
	var err error

	err = svr.RegisterRemote("Add", PackAddReq, UnpackAddRsp)
	if err != nil {
		panic(err)
	}

	err = svr.RegisterLocal("127.0.0.1:3000", Logic)
	if err != nil {
		panic(err)
	}

	svr.Run(10)

	return

}
