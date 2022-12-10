package codec

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
)

// Codec 封装接口
type Codec interface {
	Encode(i interface{}) ([]byte, error)
	Decode(data []byte, i interface{}) error
}

//将JSON、GOB格式的编解码函数进行封装，称为一个编解码器

// JSONCodec JSON格式
type JSONCodec struct{}

func (c JSONCodec) Encode(i interface{}) ([]byte, error) {
	return json.Marshal(i)
}

func (c JSONCodec) Decode(data []byte, i interface{}) error {
	decode := json.NewDecoder(bytes.NewBuffer(data))
	return decode.Decode(i)
}

// GobCodec Gob 格式
type GobCodec struct{}

func (c GobCodec) Encode(i interface{}) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(i); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
func (c GobCodec) Decode(data []byte, i interface{}) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(i)
}
