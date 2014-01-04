package main

import (
	"fmt"
	"net"
	"reflect"
	"time"
)


func RunSeqAlloc() {
	var seq uint32
	seq = 0
	for {
		seq++
		seq = (seq & CTXS_SIZE)
		seqAlloc <- seq
	}
}

func RegisterService(serviceName string) error {

	if !reflect.ValueOf(&pack).MethodByName(serviceName).IsValid() {
		return fmt.Errorf("no method name pack.%v. check service.go", serviceName)
	}
	if !reflect.ValueOf(&unpack).MethodByName(serviceName).IsValid() {
		return fmt.Errorf("no method name unpack.%v. check service.go", serviceName)
	}

	addr, err := net.ResolveUDPAddr("udp", "0.0.0.0:0")
	if err != nil {
		return err
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}
	mpConn[serviceName] = conn

	//recv
	go func() {
		for {
			func() {
				defer func() {
					if err := recover(); err != nil {
						fmt.Println(err)
					}
				}()
				unpackFunc := reflect.ValueOf(&unpack).MethodByName(serviceName)
				if !unpackFunc.IsValid() {
					panic(fmt.Errorf("no method name unpack.%v", serviceName))
				}
				for {
					bytes := make([]byte, MAX_PKG_LEN)
					n, _, err := conn.ReadFromUDP(bytes)
					if err != nil {
						fmt.Println(err)
						continue
					}

					rsp := bytes[0:n]

					res := unpackFunc.Call([]reflect.Value{reflect.ValueOf(rsp)})
					if !res[2].IsNil() {
						err = res[2].Interface().(error)
						fmt.Println(err)
						continue
					}
					seq := res[0].Interface().(uint32)
					out := res[1].Interface()

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

	if !haveRegistered {
		haveRegistered = true

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
						packFunc := reflect.ValueOf(&pack).MethodByName(ctx.ServiceName)
						if !packFunc.IsValid() {
							err := fmt.Errorf("no method name pack.%v", serviceName)
							fmt.Println(err)
							continue
						}
						res := packFunc.Call([]reflect.Value{reflect.ValueOf(ctx.Seq), reflect.ValueOf(ctx.In)})
						if !res[1].IsNil() {
							err = res[1].Interface().(error)
							fmt.Println(err)
							continue
						}
						req := res[0].Interface().([]byte)
						conn := mpConn[ctx.ServiceName]
						if conn == nil {
							err := fmt.Errorf("mpConn[%v] not exist", ctx.ServiceName)
							fmt.Println(err)
							panic(err)
						}
						_, err = conn.WriteToUDP(req, ctx.Addr)
						if err != nil {
							fmt.Println(err)
							continue
						}
					}
				}()
			}
		}()
	}
	return nil
}


func RunService(serviceName string, addr *net.UDPAddr, in interface{}, timeoutMs int) (interface{}, error) {
	seq := <-seqAlloc

	ctx := new(Ctx)
	ctxs_m[seq].Lock()
	ctxs[seq] = ctx
	ctxs_m[seq].Unlock()
	ctxs_m[seq].RLock()
	defer func() {
		ctxs_m[seq].RUnlock()
		ctxs_m[seq].Lock()
		ctxs[seq] = nil
		ctxs_m[seq].Unlock()
	} ()

	ctx.ServiceName = serviceName
	ctx.Seq = seq
	ctx.Addr = addr
	ctx.In = in
	ctx.Out = nil
	ctx.Err = nil
	recv := make(chan *Ctx, 1)
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
