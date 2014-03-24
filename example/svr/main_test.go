package main


import (
		"testing"
		"net"
		"fmt"
		"bytes"
		"encoding/binary"
	   )

func Benchmark_Add(b *testing.B) {
	b.StopTimer()

	var sum uint32

	bs := make([]byte, 128)

	conn, err := net.ListenUDP("udp", nil)
	if err != nil {
		panic(err)
	}

	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:3000")

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		p := new(bytes.Buffer)
		binary.Write(p, binary.BigEndian, uint32(i))
		binary.Write(p, binary.BigEndian, uint32(i + 1))


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

		fmt.Println(i, " + ", i + 1, " = ", sum)
	}
}
