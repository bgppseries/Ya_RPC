package server

import (
	"log"
	"reflect"
)

type Server interface {
	Register(string, interface{})
	Run()
	Close()
}
type RPCServer struct {
	listener Listener
}

func NewRPCServer(ip string, port int) *RPCServer {
	return &RPCServer{
		listener: NewRPCListener(ip, port),
	}
}
func (svr *RPCServer) Run() {
	go svr.listener.Run()
}
func (svr *RPCServer) Close() {
	if svr.listener != nil {
		svr.listener.Close()
	}
}

// Register 服务注册
func (svr *RPCServer) Register(class interface{}) {
	name := reflect.Indirect(reflect.ValueOf(class)).Type().Name()
	svr.RegisterName(name, class)
}
func (svr *RPCServer) RegisterName(name string, class interface{}) {
	handler := &RPCServerHandler{class: reflect.ValueOf(class)}
	svr.listener.SetHandler(name, handler)
	log.Printf("%s registered success!\n", name)
}
func (l *RPCListener) SetHandler(name string, handler Handler) {
	if _, ok := l.Handlers[name]; ok {
		log.Printf("%s is registered!\n", name)
		return
	}
	log.Println(name, "set handler success")
	l.Handlers[name] = handler
}
