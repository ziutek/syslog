package syslog

import (
	"fmt"
	"net"
	"time"
)

type Message struct {
	Time   time.Time
	Source net.Addr
	Facility
	Severity
	Timestamp time.Time // optional
	Hostname  string    // optional
	Tag       string
	Content   string
}

// NetAddr only network part of addr as string (IP for UDP or Name for UDS)
func (m *Message) NetAddr() string {
	switch a := m.Source.(type) {
	case *net.UDPAddr:
		return a.IP.String()
	case *net.UnixAddr:
		return a.Name
	case *net.TCPAddr:
		return a.IP.String()
	}
	// Unknown type
	return m.Source.String()
}

func (m *Message) String() string {
	timeLayout := "2006-01-02 15:04:05"
	timestampLayout := "01-02 15:04:05"
	return fmt.Sprintf(
		"%s %s <%s,%s> (%s '%s') [%s] %s",
		m.Time.Format(timeLayout), m.Source,
		m.Facility, m.Severity,
		m.Timestamp.Format(timestampLayout), m.Hostname,
		m.Tag, m.Content,
	)
}
