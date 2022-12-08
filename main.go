package ya_rpc

import (
	"encoding/gob"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"ya-rpc/server"
)

func main() {
	flag.Parse()
	if ip == "" || port == 0 {
		panic("init ip and port error")
	}
	srv := server.NewRPCServer(ip, port)
	srv.RegisterName("User", &UserHandler{})
	srv.RegisterName("Test", &TestHandler{})
	gob.Register(User{})
	go srv.Run()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	<-quit
	srv.Close()
}
