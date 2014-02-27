package main

import (
	"bytes"
	"encoding/binary"
)


func PackAddReq(seq uint32, in interface{}) ([]byte, error) {
	addIn := in.(*AddIn)

	p := new(bytes.Buffer)
	binary.Write(p, binary.BigEndian, seq)
	binary.Write(p, binary.BigEndian, addIn.A)
	binary.Write(p, binary.BigEndian, addIn.B)

	return p.Bytes(), nil
}

func UnpackAddRsp(pkg []byte) (uint32, interface{}, error) {
	var seq uint32
	var addOut AddOut

	p := bytes.NewBuffer(pkg)
	binary.Read(p, binary.BigEndian, &seq)
	binary.Read(p, binary.BigEndian, &addOut.Sum)

	return seq, &addOut, nil
}
