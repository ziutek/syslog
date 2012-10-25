package main

import (
	"github.com/ziutek/syslog"
	"time"
)

func main() {
	s := syslog.NewServer(5)
	s.Listen("0.0.0.0:1514")
	time.Sleep(time.Hour)
}
