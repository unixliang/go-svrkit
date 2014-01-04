package main

import(
		"net"
		"bytes"
		"encoding/binary"
	  )

func main() {
	var seq uint32
	var a uint32
	var b uint32
	var sum uint32

	bs := make([]byte, 128)

	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:4000")

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		panic(err)
	}


	for {
		n, fromAddr, err := conn.ReadFromUDP(bs)
		if err != nil {
			panic(err)
		}
		p := bytes.NewBuffer(bs[0 : n])
		binary.Read(p, binary.BigEndian, &seq)
		binary.Read(p, binary.BigEndian, &a)
		binary.Read(p, binary.BigEndian, &b)

		sum = a + b

		p = new(bytes.Buffer)
		binary.Write(p, binary.BigEndian, seq)
		binary.Write(p, binary.BigEndian, sum)

		_, err = conn.WriteToUDP(p.Bytes(), fromAddr)
		if err != nil {
			panic(err)
		}
	}
}
