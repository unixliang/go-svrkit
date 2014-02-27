package go_udp_svrkit

import (
	"fmt"
	"net"
	"time"
)

func init() {
	go RunSeqAlloc()
	//send
	go func() {
		for {
			func() {
				defer func() {
					if err := recover(); err != nil {
						fmt.Println(err)
					}
				}()
				for {
					ctx := <-send
					if ctx.Remote.PackReq == nil {
						continue
					}
					req, err := ctx.Remote.PackReq(ctx.Seq, ctx.In)
					if err != nil {
						panic(err)
					}
					if ctx.Remote.Conn == nil {
						err := fmt.Errorf("ctx.Remote.Conn not exist")
						panic(err)
					}
					_, err = ctx.Remote.Conn.WriteToUDP(req, ctx.Addr)
					if err != nil {
						panic(err)
					}
				}
			}()
		}
	}()
}

func RunSeqAlloc() {
	var seq uint32
	seq = 0
	for {
		seq++
		seq = (seq & CTXS_SIZE)
		seqAlloc <- seq
	}
}

func RegisterRemote(remoteName string,
	packReq func(uint32, interface{}) ([]byte, error),
	unpackRsp func([]byte) (uint32, interface{}, error)) error {

	if packReq == nil {
		return fmt.Errorf("packReq is nil")
	}

	if unpackRsp == nil {
		return fmt.Errorf("unpackRsp is nil")
	}

	addr, err := net.ResolveUDPAddr("udp", "0.0.0.0:0")
	if err != nil {
		return err
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}
	remote := new(Remote_t)
	remote.Conn = conn
	remote.PackReq = packReq
	remote.UnpackRsp = unpackRsp
	mpRemotes[remoteName] = remote

	//recv
	go func() {
		for {
			func() {
				defer func() {
					if err := recover(); err != nil {
						fmt.Println(err)
					}
				}()
				bytes := make([]byte, MAX_PKG_LEN)
				for {
					n, _, err := conn.ReadFromUDP(bytes)
					if err != nil {
						panic(err)
					}

					rsp := bytes[0:n]

					if unpackRsp == nil {
						continue
					}

					seq, out, err := unpackRsp(rsp)
					if err != nil {
						panic(err)
					}

					ctxs_m[seq].RLock()
					ctx := ctxs[seq]
					ctxs_m[seq].RUnlock()
					if ctx == nil {
						err := fmt.Errorf("ctxs[seq] not exist. seq: %v", seq)
						fmt.Println(err)
					} else {
						ctx.Out = out
						ctx.Err = nil
						ctx.Recv <- ctx
					}
				}
			}()
		}
	}()

	return nil
}

func RunRemote(remoteName string, addr *net.UDPAddr, in interface{}, timeoutMs int) (interface{}, error) {

	remote := mpRemotes[remoteName]
	if remote == nil {
		return nil, fmt.Errorf("mpRemotes[%v] not exise.", remote)
	}

	seq := <-seqAlloc

	ctx := new(Ctx_t)
	ctxs_m[seq].Lock()
	ctxs[seq] = ctx
	ctxs_m[seq].Unlock()
	ctxs_m[seq].RLock()
	defer func() {
		ctxs_m[seq].RUnlock()
		ctxs_m[seq].Lock()
		ctxs[seq] = nil
		ctxs_m[seq].Unlock()
	}()

	ctx.Remote = remote
	ctx.Seq = seq
	ctx.Addr = addr
	ctx.In = in
	ctx.Out = nil
	ctx.Err = nil
	recv := make(chan *Ctx_t, 1)
	ctx.Recv = recv

	send <- ctx
	select {
	case ctx = <-recv:
	case <-time.After(time.Millisecond * time.Duration(timeoutMs)):
		return nil, fmt.Errorf("time out")
	}

	if ctx.Err != nil {
		return nil, ctx.Err
	}
	return ctx.Out, nil

}
