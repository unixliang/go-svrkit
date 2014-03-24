package svrkit

import (
	"fmt"
	"runtime"
	"sync"
)

func Run(maxProcs int, logPrefix string, logPriority int) {

	NewLogger(logPrefix, logPriority)

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
			for {
				bytes := make([]byte, MAX_PKG_LEN)
				n, addr, err := local.Conn.ReadFromUDP(bytes)
				if err != nil {
					panic(err)
				}

				seq := <-seqAlloc

				req := bytes[0:n]
				go local.Logic(seq, local.Conn, req, addr)
			}
		}()
	}
	wg.Wait()
}
