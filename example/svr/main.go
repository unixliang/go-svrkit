package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	svrkit "go-svrkit"

)


func Logic(seq uint32, conn *net.UDPConn, req []byte, fromAddr *net.UDPAddr) {
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

	res, err := svrkit.RunRemote(seq, "Add", addr, &addIn, 1000)
	if err != nil {
		panic(fmt.Errorf("Add err: %v", err))
	}

	addOut := res.(*AddOut)

	svrkit.Infof("seq: %v a: %v b: %v sum: %v", seq, addIn.A, addIn.B, addOut.Sum)

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

	err = svrkit.RegisterRemote("Add", PackAddReq, UnpackAddRsp)
	if err != nil {
		panic(err)
	}

	err = svrkit.RegisterLocal("127.0.0.1:3000", Logic)
	if err != nil {
		panic(err)
	}

	svrkit.Run(10, "./svr", 7)

	return

}
