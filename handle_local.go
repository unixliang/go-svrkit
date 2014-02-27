package go_udp_svrkit

import (
	"fmt"
	"net"
)

func RegisterLocal(addrStr string,
	logic func(*net.UDPConn, []byte, *net.UDPAddr)) error {

	if logic == nil {
		return fmt.Errorf("logic is nil")
	}

	addr, err := net.ResolveUDPAddr("udp", addrStr)
	if err != nil {
		return err
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}

	local := new(Local_t)
	local.Conn = conn
	local.Logic = logic

	liLocal = append(liLocal, local)

	return nil
}
