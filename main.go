package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"ya-rpc/client"
	"ya-rpc/config"
	"ya-rpc/server"
)

var (
	f bool
)

func init() {
	flag.BoolVar(&f, "f", true, "is server or client")
}

// 一个demo
func main() {
	flag.Parse()
	if f {
		//服务器代码
		srv := server.NewRPCServer("127.0.0.1", config.PORT)
		srv.RegisterName("Test", &SumHandler{})
		srv.RegisterName("String", &UppercaseHandler{}) //注册函数sum 和 uppercase
		go srv.Run()
		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
		<-quit
		srv.Close()
	} else {
		//客户端代码
		cli := client.NewClientProxy(client.DefaultOption)

		ctx := context.Background()
		var sum func(a, b float32) float32
		var (
			a float32 = 2.0165
			b float32 = 3.45
			c float32 = 4.4545
		)

		_, err := cli.Call(ctx, "Userservice.Test.Sum", &sum, a, b)
		if err != nil {
			log.Println("sum func remote call error:", err)
		}

		u := sum(a, c)
		v := sum(a, b)
		log.Println(u, v)
		cls := client.NewClientProxy(client.DefaultOption)
		cts := context.Background()
		var uppercase func(s string) string
		var str = "hello world!"
		r, err := cls.Call(cts, "Userservice.String.Uppercase", &uppercase, str)
		r = uppercase(str)
		t := uppercase("what? fuck")
		log.Println("uppercase func demo result:", r, t)
	}
}

type SumHandler struct{}

func (s *SumHandler) Sum(a, b float32) float32 {
	return a + b
}

type UppercaseHandler struct{}

func (u *UppercaseHandler) Uppercase(s string) string {
	return strings.ToUpper(s)
}
