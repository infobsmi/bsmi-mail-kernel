package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/json-iterator/go"
	"gopkg.in/natefinch/npipe.v2"
	"net"
)



func newTestServer() *rpc.Server {
	server := rpc.NewServer()
	if err := server.RegisterName("test", new(testService)); err != nil {
		panic(err)
	}
	return server
}



type testService struct{}

type echoArgs struct {
	S string
}

type echoResult struct {
	Method string
	String string
	Int    int
	Args   *echoArgs
}

func (s *testService) Echo(str string, i int, args *echoArgs) echoResult {
	return echoResult{"test_echo", str, i, args}
}
func main() {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	data := 1
	marshal, err := json.Marshal(&data)
	if err != nil {
		return
	}
	fmt.Printf("json: %+v", marshal)
	fmt.Println("Hello world")


	ln, err := npipe.Listen(`\\.\pipe\bsmi_mail_kernel`)
	if err != nil {
		// handle error
	}
	server := newTestServer()
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			continue
		}

		go server.ServeCodec(rpc.NewCodec(conn), 0)
	}

}

func handleConnection(conn net.Conn) {

}
