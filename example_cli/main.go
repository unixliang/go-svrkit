package main

import(
		"net"
		"fmt"
		"bytes"
		"encoding/binary"
	  )

func main() {
	var sum uint32

	bs := make([]byte, 128)

	conn, err := net.ListenUDP("udp", nil)
	if err != nil {
		panic(err)
	}

	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:3000")

	for i := 0; i < 100; i++ {
		p := new(bytes.Buffer)
		binary.Write(p, binary.BigEndian, uint32(i))
		binary.Write(p, binary.BigEndian, uint32(i + 1))

		fmt.Println("send: ", i, " + ", i + 1, " = ")

		_, err = conn.WriteToUDP(p.Bytes(), addr)
		if err != nil {
			panic(err)
		}

		n, _, err := conn.ReadFromUDP(bs)
		if err != nil {
			panic(err)
		}
		p = bytes.NewBuffer(bs[0 : n])
		binary.Read(p, binary.BigEndian, &sum)
		fmt.Println("recv: sum = ", sum)
	}
}
