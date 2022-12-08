package config

import (
	"errors"
	"ya-rpc/codec"
	"ya-rpc/protocol"
)

const (
	TRANS_TYPE = "gob" //json
	HEADER_LEN = 4     //header length
)

type CodecMode int

const (
	CODEC_GOB CodecMode = iota
	CODEC_JSON
	NET_TRANS_PROTOCOL = "TCP"
)

const (
	Protocol_MsgVersion = 1
)

var Codecs = map[protocol.SerializeType]codec.Codec{
	protocol.JSON: &codec.JSONCodec{},
	protocol.Gob:  &codec.GobCodec{},
}
var ParamErr = errors.New("params not adapted")
