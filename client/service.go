package client

import (
	"errors"
	"strings"
	"ya-rpc/config"
)

type Service struct {
	AppId  string   //客户端ID
	Class  string   //类名
	Method string   //方法名
	Addrs  []string //服务器地址
}

// demo: UserService.Test.sum
func NewService(servicePath string) (*Service, error) {
	arr := strings.Split(servicePath, ".")
	service := &Service{}
	if len(arr) != 3 {
		return service, errors.New("service path illegal")
	}
	service.AppId = arr[0]
	service.Class = arr[1]
	service.Method = arr[2]
	return service, nil
}
func (service *Service) SelectAddr() string {
	//todo 服務中心
	return config.ADDR
	//return "127.0.0.1:4545"
}
