package client

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"reflect"
	"time"
	"ya-rpc/config"
	"ya-rpc/protocol"
)

type Client interface {
	Connect(string) error
	Invoke(context.Context, *Service, interface{}, ...interface{}) (interface{}, error)
	Close()
}

// Option 定义连接参数，设置重试次数、超时时间、序列化协议、压缩类型等。
type Option struct {
	Retries           int                    //重试次数
	ConnectionTimeout time.Duration          //超时时间
	SerializeType     protocol.SerializeType //序列化协议
	CompressType      protocol.CompressType  //压缩类型
}

var DefaultOption = Option{
	Retries:           3,
	ConnectionTimeout: 5 * time.Second,
	SerializeType:     protocol.Gob,
	CompressType:      protocol.None,
} //默认参数

type RPCClient struct {
	conn   net.Conn
	option Option
}

func NewClient(option Option) Client {
	return &RPCClient{option: option}
}
func (cli *RPCClient) Connect(addr string) error {
	conn, err := net.DialTimeout(config.NET_TRANS_PROTOCOL, addr, cli.option.ConnectionTimeout)
	if err != nil {
		return err
	}
	cli.conn = conn
	return nil
}
func (cli *RPCClient) Invoke(ctx context.Context, service *Service, stub interface{}, params ...interface{}) (interface{}, error) {
	cli.makeCall(service, stub)
	return cli.wrapCall(ctx, stub, params...)
}
func (cli *RPCClient) Close() {
	if cli.conn != nil {
		err := cli.conn.Close()
		if err != nil {
			return
		}
	}
}
func (cli *RPCClient) makeCall(service *Service, methodPtr interface{}) {
	container := reflect.ValueOf(methodPtr).Elem()
	coder := config.Codecs[cli.option.SerializeType]

	handler := func(req []reflect.Value) []reflect.Value {
		numOut := container.Type().NumOut()
		errorHandler := func(err error) []reflect.Value {
			outArgs := make([]reflect.Value, numOut)
			for i := 0; i < len(outArgs)-1; i++ {
				outArgs[i] = reflect.Zero(container.Type().Out(i))
			}
			outArgs[len(outArgs)-1] = reflect.ValueOf(&err).Elem()
			return outArgs
		}
		inArgs := make([]interface{}, 0, len(req))
		for _, arg := range req {
			inArgs = append(inArgs, arg.Interface())
		}
		payload, err := coder.Encode(inArgs) //[]byte
		if err != nil {
			log.Printf("encode err:%v\n", err)
			return errorHandler(err)
		}
		msg := protocol.NewRPCMsg()
		msg.SetVersion(config.Protocol_MsgVersion)
		msg.SetMsgType(protocol.Request)
		msg.SetCompressType(cli.option.CompressType)
		msg.SetSerializeType(cli.option.SerializeType)
		msg.ServiceClass = service.Class
		msg.ServiceMethod = service.Method
		msg.Payload = payload
		err = msg.Send(cli.conn)
		if err != nil {
			log.Printf("send err:%v\n", err)
			return errorHandler(err)
		}
		respMsg, err := protocol.Read(cli.conn)
		if err != nil {
			return errorHandler(err)
		}
		respDecode := make([]interface{}, 0)
		err = coder.Decode(respMsg.Payload, &respDecode)
		if err != nil {
			log.Printf("decode err:%v\n", err)
			return errorHandler(err)
		}
		if len(respDecode) == 0 {
			respDecode = make([]interface{}, numOut)
		}
		outArgs := make([]reflect.Value, numOut)
		for i := 0; i < numOut; i++ {
			if i != numOut {
				if respDecode[i] == nil {
					outArgs[i] = reflect.Zero(container.Type().Out(i))
				} else {
					outArgs[i] = reflect.ValueOf(respDecode[i])
				}
			} else {
				outArgs[i] = reflect.Zero(container.Type().Out(i))
			}
		}
		return outArgs
	}
	container.Set(reflect.MakeFunc(container.Type(), handler))
}
func (cli *RPCClient) wrapCall(ctx context.Context, stub interface{}, params ...interface{}) (interface{}, error) {
	f := reflect.ValueOf(stub).Elem()
	if len(params) != f.Type().NumIn() {
		return nil, errors.New(fmt.Sprintf("params not adapted: %d-%d", len(params), f.Type().NumIn()))
	}
	in := make([]reflect.Value, len(params))
	for idx, param := range params {
		in[idx] = reflect.ValueOf(param)
	}
	result := f.Call(in)
	return result, nil
}
