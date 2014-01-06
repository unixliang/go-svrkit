package main

import (
	"net"
	"sync"
)

const (
	MAX_PKG_LEN = 10000
	CTXS_SIZE   = 1048575 //2^20 - 1
)

var mpConn = make(map[string]*net.UDPConn)

var listenConn *net.UDPConn

var recvFromListenConn = make(chan []byte)
var recvFromSendConn = make(chan []byte)
var send = make(chan *Ctx)

var pack = Pack{}
var unpack = Unpack{}

var ctxs = make([]*Ctx, CTXS_SIZE+1)
var ctxs_m = make([]sync.RWMutex, CTXS_SIZE+1)

var seqAlloc = make(chan uint32)
