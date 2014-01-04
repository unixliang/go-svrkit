package main

import (
	"net"
)

type Ctx struct {
	ServiceName string
	Seq         uint32
	Addr        *net.UDPAddr
	In          interface{}
	Out         interface{}
	Err         error
	Recv        chan *Ctx
}

type Pack struct{}
type Unpack struct{}
