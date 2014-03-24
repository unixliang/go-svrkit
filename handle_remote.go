package svrkit

import (
	"fmt"
	"net"
	"time"
)

func init() {
	go RunSeqAlloc()
	//send
	for i := 0; i < 8; i++ {
		go func() {
			for {
				func() {
					var ctx *Ctx_t
					defer func() {
						if err := recover(); err != nil {
							ctx.Err = err.(error)
							ctx.Recv <- ctx
						}
					}()
					for {
						ctx = <-send
						if ctx.Remote.PackReq == nil {
							panic(fmt.Errorf("ctx.Remote.PackReq == nil"))
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
}

func RunSeqAlloc() {
	var seq uint32
	seq = 0
	for {
		seq++
		if seq == 0 {
			seq++
		}
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
	for i := 0; i < 8; i++ {
		go func() {
			for {
				func() {
					var ctx *Ctx_t
					ctx = nil
					defer func() {
						if err := recover(); err != nil {
							if ctx != nil {
								ctx.Out = nil
								ctx.Err = err.(error)
								ctx.Recv <- ctx
							} else {
								Err(err)
							}
						}
					}()
					bytes := make([]byte, MAX_PKG_LEN)
					for {
						ctx = nil

						n, _, err := conn.ReadFromUDP(bytes)
						if err != nil {
							panic(err)
						}

						rsp := bytes[0:n]

						if unpackRsp == nil {
							panic(fmt.Errorf("unpackRsp == nil"))
						}

						seq, out, err := unpackRsp(rsp)
						if seq > 0 {
							ctxs_m[seq].RLock()
							ctx = ctxs[seq]
							ctxs_m[seq].RUnlock()
						}
						if err != nil {
							panic(err)
						}
						if ctx == nil {
							err := fmt.Errorf("ctxs[seq] not exist. seq: %v", seq)
							panic(err)
						} else {
							ctx.Out = out
							ctx.Err = nil
							ctx.Recv <- ctx
						}
					}
				}()
			}
		}()
	}

	return nil
}

func RunRemote(seq uint32, remoteName string, addr *net.UDPAddr, in interface{}, timeoutMs int) (interface{}, error) {

	remote := mpRemotes[remoteName]
	if remote == nil {
		return nil, fmt.Errorf("mpRemotes[%v] not exise.", remote)
	}

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
