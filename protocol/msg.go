package protocol

import (
	"encoding/binary"
	"fmt"
	"io"
)

type RPCMsg struct {
	*Header
	ServiceClass  string
	ServiceMethod string
	Payload       []byte
}

const SplitLen = 4

func NewRPCMsg() *RPCMsg {
	header := Header([HeaderLen]byte{})
	header[0] = magicNumber
	return &RPCMsg{
		Header: &header,
	}
}

func (msg *RPCMsg) Send(writer io.Writer) error {
	//发送报文头
	_, err := writer.Write(msg.Header[:])
	if err != nil {
		return err
	}

	//发送报文总长度
	dataLen := SplitLen + len(msg.ServiceClass) + SplitLen + len(msg.ServiceMethod) + SplitLen + len(msg.Payload)
	err = binary.Write(writer, binary.BigEndian, uint32(dataLen)) //4
	if err != nil {
		return err
	}

	//发送要申请的方法的类的长度 4byte
	err = binary.Write(writer, binary.BigEndian, uint32(len(msg.ServiceClass)))
	if err != nil {
		return err
	}

	//发送要申请的方法的类
	err = binary.Write(writer, binary.BigEndian, StringToByte(msg.ServiceClass))
	if err != nil {
		return err
	}

	//发送要申请的方法的长度 4byte
	err = binary.Write(writer, binary.BigEndian, uint32(len(msg.ServiceMethod)))
	if err != nil {
		return err
	}

	//发送要申请的方法
	err = binary.Write(writer, binary.BigEndian, StringToByte(msg.ServiceMethod))
	if err != nil {
		return err
	}

	//发送函数参数的长度 4byte
	err = binary.Write(writer, binary.BigEndian, uint32(len(msg.Payload)))
	if err != nil {
		return err
	}

	//发送函数参数
	_, err = writer.Write(msg.Payload)
	if err != nil {
		return err
	}
	return nil
}

// Decode 对收到的报文进行解码
func (msg *RPCMsg) Decode(r io.Reader) error {

	_, err := io.ReadFull(r, msg.Header[:])
	if !msg.Header.CheckMagicNumber() { //magicNumber
		return fmt.Errorf("magic number error,data is wrong: %v", msg.Header[0])
	}

	headerByte := make([]byte, 4)
	_, err = io.ReadFull(r, headerByte)
	if err != nil {
		return err
	}
	bodyLen := binary.BigEndian.Uint32(headerByte)

	data := make([]byte, bodyLen)
	_, err = io.ReadFull(r, data)

	start := 0
	end := start + SplitLen
	classLen := binary.BigEndian.Uint32(data[start:end]) //0,4

	start = end
	end = start + int(classLen)
	msg.ServiceClass = ByteToString(data[start:end]) //4,x

	start = end
	end = start + SplitLen
	methodLen := binary.BigEndian.Uint32(data[start:end]) //x,x+4

	start = end
	end = start + int(methodLen)
	msg.ServiceMethod = ByteToString(data[start:end]) //x+4, x+4+y

	start = end
	end = start + SplitLen
	binary.BigEndian.Uint32(data[start:end]) //x+4+y, x+y+8 payloadLen

	start = end
	msg.Payload = data[start:]
	return nil
}

func Read(r io.Reader) (*RPCMsg, error) {
	msg := NewRPCMsg()
	err := msg.Decode(r)
	if err != nil {
		return nil, err
	}
	return msg, nil
}
