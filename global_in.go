package go_udp_svrkit

import (
	"sync"
)

const (
	MAX_PKG_LEN = 10000
	CTXS_SIZE   = 1048575 //2^20 - 1
)

var mpRemotes = make(map[string]*Remote_t)
var liLocal = []*Local_t{}

var recvFromListenConn = make(chan []byte)
var recvFromSendConn = make(chan []byte)
var send = make(chan *Ctx_t)

var ctxs = make([]*Ctx_t, CTXS_SIZE+1)
var ctxs_m = make([]sync.RWMutex, CTXS_SIZE+1)

var seqAlloc = make(chan uint32)
