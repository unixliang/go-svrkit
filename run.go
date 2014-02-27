package go_udp_svrkit

import (
	"runtime"
	"fmt"
	"sync"
	)

func Run(maxProcs int) {

	runtime.GOMAXPROCS(maxProcs)
	
	var wg sync.WaitGroup
	
	for _, local := range liLocal {
		wg.Add(1)
		go func() {
			defer func() {
				if err := recover(); err != nil {
					fmt.Println(err)
				}
			}()
			bytes := make([]byte, MAX_PKG_LEN)
			for {
				n, addr, err := local.Conn.ReadFromUDP(bytes)
				if err != nil {
					panic(err)
				}

				req := bytes[0:n]
				go local.Logic(local.Conn, req, addr)
			}
		} ()
	}
	wg.Wait()
}
