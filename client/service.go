package client

import (
	"errors"
	"strings"
)

type Service struct {
	AppId  string
	Class  string
	Method string
	Addrs  []string
}

// demo: UserService.user.GetUser
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
	return "ip:8811"
}
