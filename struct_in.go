package go_udp_svrkit

import (
	"net"
)

type Local_t struct {
	Conn	*net.UDPConn
	Logic func(*net.UDPConn, []byte, *net.UDPAddr)
}

type Remote_t struct {
	Conn	*net.UDPConn
	PackReq	func(uint32, interface{}) ([]byte, error)
	UnpackRsp	func([]byte) (uint32, interface{}, error)
}

type Ctx_t struct {
	Remote		*Remote_t
	Seq         uint32
	Addr        *net.UDPAddr
	In          interface{}
	Out         interface{}
	Err         error
	Recv        chan *Ctx_t
}

