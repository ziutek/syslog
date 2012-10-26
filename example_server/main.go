package main

import (
	"fmt"
	"github.com/ziutek/syslog"
	"os"
	"os/signal"
	"syscall"
)

type handler struct {
	*syslog.BaseHandler
}

func filter(m *syslog.Message) bool {
	return m.Tag == "named" || m.Tag == "bind"
}

func newHandler() *handler {
	h := handler{syslog.NewBaseHandler(5, filter, false)}
	go h.mainLoop()
	return &h
}

func (h *handler) mainLoop() {
	for {
		m := h.Get()
		if m == nil {
			break
		}
		fmt.Println(m)
	}
	fmt.Println("Exit handler")
	h.End()
}

func main() {
	var s syslog.Server
	s.AddHandler(newHandler())
	s.Listen("0.0.0.0:1514")

	sc := make(chan os.Signal, 2)
	signal.Notify(sc, syscall.SIGTERM, syscall.SIGINT)
	<-sc

	fmt.Println("Shutdown the server...")
	s.Shutdown()
	fmt.Println("Server is down")
}
