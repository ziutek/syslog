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
