// Syslog server library
// Currently it handles only RFC 3164 packets
package syslog

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
	"unicode"
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

type Server struct {
	q chan Message
}

func NewServer(bufLen int) *Server {
	return &Server{q: make(chan Message, bufLen)}
}

// Listen adds to the server next listen address which can be a path for unix
// domain socket or host:port for UDP socket.
func (s *Server) Listen(addr string) error {
	var c *net.UDPConn
	if strings.IndexRune(addr, ':') != -1 {
		a, err := net.ResolveUDPAddr("udp", addr)
		if err != nil {
			return err
		}
		c, err = net.ListenUDP("udp", a)
		if err != nil {
			return err
		}
	} else {
		a, err := net.ResolveUnixAddr("unixgram", addr)
		if err != nil {
			return err
		}
		c, err = net.ListenUnixgram("unixgram", a)
		if err != nil {
			return err
		}
	}
	go s.receiver(c)
	return nil
}

func isNotAlnum(r rune) bool {
	return !(unicode.IsLetter(r) || unicode.IsNumber(r))
}

func isNulCrLf(r rune) bool {
	return r == 0 || r == '\r' || r == '\n'
}

func (s *Server) receiver(c *net.UDPConn) {
	//q := (chan<- Message)(s.q)
	buf := make([]byte, 1024)
	for {
		var m Message
		n, addr, err := c.ReadFrom(buf)
		if err != nil {
			log.Println("Read error:", err)
			return
		}
		pkt := buf[:n]

		m.Source = addr
		m.Time = time.Now()

		// Parse priority (if exists)
		prio := 13 // default priority
		hasPrio := false
		if pkt[0] == '<' {
			n = 1 + bytes.IndexByte(pkt[1:], '>')
			if n > 1 && n < 5 {
				p, err := strconv.Atoi(string(pkt[1:n]))
				if err == nil && p >= 0 {
					hasPrio = true
					prio = p
					pkt = pkt[n+1:]
				}
			}
		}
		m.Severity = Severity(prio & 0x07)
		m.Facility = Facility(prio >> 3)

		// Parse header (if exists)
		if hasPrio && len(pkt) >= 16 && pkt[15] == ' ' {
			// Get timestamp
			layout := "Jan _2 15:04:05"
			ts, err := time.Parse(layout, string(pkt[:15]))
			if err == nil && !ts.IsZero() {
				// Get hostname
				n = 16 + bytes.IndexByte(pkt[16:], ' ')
				if n != 15 {
					m.Timestamp = ts
					m.Hostname = string(pkt[16:n])
					pkt = pkt[n+1:]
				}
			}
			// TODO: check for version an new format of timestamp as
			// described in RFC 5424.
		}

		// Parse msg part
		pkt = bytes.TrimRightFunc(pkt, isNulCrLf)
		n = bytes.IndexFunc(pkt, isNotAlnum)
		if n != -1 {
			m.Tag = string(pkt[:n])
			pkt = pkt[n:]
		}
		m.Content = string(pkt)

		fmt.Println(m.String())
	}
}
