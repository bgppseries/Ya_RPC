package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"ya-rpc/config"
	"ya-rpc/protocol"
)

// 封装listen
// 使用TCP传输

type Listener interface {
	Run()
	SetHandler(string, Handler)
	Close()
}
type RPCListener struct {
	ServiceIp   string
	ServicePort int
	Handlers    map[string]Handler
	nl          net.Listener
}

func NewRPCListener(serviceIp string, servicePort int) *RPCListener {
	return &RPCListener{
		ServiceIp:   serviceIp,
		ServicePort: servicePort,
		Handlers:    make(map[string]Handler)}
}
func (l *RPCListener) Run() {
	addr := fmt.Sprintf("%s:%d", l.ServiceIp, l.ServicePort)
	println(addr)
	//nl, err := net.Listen(config.NET_TRANS_PROTOCOL, addr)
	nl, err := net.Listen("tcp", "0.0.0.0:4545") //tcp
	if err != nil {
		log.Println("TCP server listen err:", err)
		return
	}
	l.nl = nl
	for {
		conn, err := l.nl.Accept()
		if err != nil {
			continue
		}
		log.Println("the server accept from", conn.RemoteAddr())
		go l.handleConn(conn)
	}
}
func (l *RPCListener) Close() {
	if l.nl != nil {
		err := l.nl.Close()
		if err != nil {
			log.Println("RPC listener close err:", err)
			return
		}
	}
}
func (l *RPCListener) handleConn(conn net.Conn) {
	//catch panic
	defer func() {
		if err := recover(); err != nil {
			log.Printf("server %s catch panic err:%s\n", conn.RemoteAddr(), err)
		}
		err := conn.Close()
		if err != nil {
			return
		}
	}()
	for {
		msg, err := l.receiveData(conn) //接受数据并反序列化，得到msg结构化的数据
		println("receive req:", msg.ServiceMethod, msg.ServiceClass, msg.Payload)
		if err != nil || msg == nil {

			return
		}
		coder := config.Codecs[msg.Header.SerializeType()]
		if coder == nil {
			log.Println("server coder defined err", err)
			return
		}
		inArgs := make([]interface{}, 0)
		err = coder.Decode(msg.Payload, &inArgs)
		if err != nil {
			log.Println("the func in Decode error:", err)
			return
		}
		log.Println("receive input:", inArgs[0])
		handler, ok := l.Handlers[msg.ServiceClass]
		if !ok {
			println("????", ok)
			return
		}

		result, err := handler.Handle(msg.ServiceMethod, inArgs)
		println("result:", result)
		encodeRes, err := coder.Encode(result)
		if err != nil {
			println("!!!!!")
			return
		}
		err = l.sendData(conn, encodeRes)
		if err != nil {
			println("&&&&&&")
			return
		}
	}
}
func (l *RPCListener) receiveData(conn net.Conn) (*protocol.RPCMsg, error) {
	msg, err := protocol.Read(conn)
	if err != nil {
		if err != io.EOF { //close
			return nil, err
		}
	}
	return msg, nil
}
func (l *RPCListener) sendData(conn net.Conn, payload []byte) error {
	resMsg := protocol.NewRPCMsg()
	resMsg.SetVersion(config.Protocol_MsgVersion)
	resMsg.SetMsgType(protocol.Response)
	resMsg.SetCompressType(protocol.None)
	resMsg.SetSerializeType(protocol.Gob)
	resMsg.Payload = payload
	return resMsg.Send(conn)
}
