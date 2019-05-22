package server

import (
	"fmt"
	"net"
	"strings"
)

type Listener struct {
	host     string
	port     string
	listener net.Listener
}

func NewListener(host string, port string) *Listener {
	listener, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		return nil
	}
	return &Listener{host, port, listener}
}

type RequestHandler struct {
	conn net.Conn
}

func (listener *Listener) Accept() *RequestHandler {
	conn, err := listener.listener.Accept()
	if err != nil {
		fmt.Println(err.Error())
	}
	return newRequestHandler(conn)
}

func newRequestHandler(conn net.Conn) *RequestHandler {
	return &RequestHandler{conn}
}

func (handler *RequestHandler) Send(data []byte) {
	handler.conn.Write(data)
}

func (handler *RequestHandler) Receive() []byte {
	byteMsg := make([]byte, 2048)
	read, err := handler.conn.Read(byteMsg)
	if err != nil {
		return nil
	}
	return byteMsg[:read]
}

func (handler *RequestHandler) GetHost() string {
	return strings.Split(handler.conn.RemoteAddr().String(), ":")[0]
}
