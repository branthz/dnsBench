package main

import (
	"net"
	"time"

	"strconv"
	"strings"

	"github.com/miekg/dns"
)

const (
	year68     = 1 << 31 // For RFC1982 (Serial Arithmetic) calculations in 32 bits.
	defaultTtl = 3600    // Default internal TTL.

	// DefaultMsgSize is the standard default for messages larger than 512 bytes.
	DefaultMsgSize = 4096
	// MinMsgSize is the minimal size of a DNS packet.
	MinMsgSize = 512
	// MaxMsgSize is the largest possible DNS packet.
	MaxMsgSize = 65535
)

type client struct {
	conn    *net.UDPConn
	Net     string      // if "tcp" or "tcp-tls" (DNS over TLS) a TCP query will be initiated, otherwise an UDP one (default is "" for UDP)
	UDPSize uint16      // minimum receive buffer for UDP messages
	Dialer  *net.Dialer // a net.Dialer used to set local address, timeouts and more
	// Timeout is a cumulative timeout for dial, write and read, defaults to 0 (disabled) - overrides DialTimeout, ReadTimeout,
	// WriteTimeout when non-zero. Can be overridden with net.Dialer.Timeout (see Client.ExchangeWithDialer and
	// Client.Dialer) or context.Context.Deadline (see the deprecated ExchangeContext)
	Timeout      time.Duration
	DialTimeout  time.Duration     // net.DialTimeout, defaults to 2 seconds, or net.Dialer.Timeout if expiring earlier - overridden by Timeout when that value is non-zero
	ReadTimeout  time.Duration     // net.Conn.SetReadTimeout value for connections, defaults to 2 seconds - overridden by Timeout when that value is non-zero
	WriteTimeout time.Duration     // net.Conn.SetWriteTimeout value for connections, defaults to 2 seconds - overridden by Timeout when that value is non-zero
	TsigSecret   map[string]string // secret(s) for Tsig map[<zonename>]<base64 secret>, zonename must be in canonical form (lowercase, fqdn, see RFC 4034 Section 6.2)
}

func NewClient() (*client, error) {
	c := new(client)
	host := P.server + ":" + strconv.Itoa(P.port)
	dst, err := net.ResolveUDPAddr("udp4", host)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialUDP("udp4", nil, dst)
	if err != nil {
		mlog.Error("conn failed:%v", err)
		return nil, err
	}
	//conn.SetReadDeadline(time.Now().Add(time.Second*3))
	//conn.SetWriteDeadline(time.Now().Add(time.Second*3))
	c.UDPSize = 1024
	c.conn = conn
	return c, nil
}

func (c *client) Query(d []string) error {
	var dl = len(d)
	var data []byte
	var err error
	mlog.Debug("remote addr:%s", c.conn.RemoteAddr().String())
	for j := 0; j < P.recycle; j++ {
		for i := 0; i < dl; i++ {
			State.sendNum++
			m := NewMsg(d[i])
			data, err = m.Pack()
			if err != nil {
				mlog.Error("pack failed:%v", err)
				break
			}
			_, err = c.conn.Write(data)
			if err != nil {
				mlog.Error("write failed:%v\n", err)
			}
		}
	}
	time.Sleep(time.Second * waitForClose)
	c.conn.Close()
	mlog.Warnln("query finished..............\n")
	return nil
}

func (c *client) Response() {
	var data = make([]byte, 1460)
	var rn int
	var err error
	for {
		rn, err = c.conn.Read(data)
		if err != nil {
			if strings.Contains(err.Error(), "closed") {
				mlog.Warnln("connection closed!")
				break
			}
			mlog.Error(err.Error())
			continue
		}
		State.recvNum++
		//deal data
		m := new(dns.Msg)
		err = m.Unpack(data[:rn])
		if err != nil {
			mlog.Error("parse respon failed:%v", err)
			continue
		} else {
			mlog.Debug("get reponse:\n%v", m)
		}
		if m.Rcode == dns.RcodeSuccess {
			State.succNum++
		} else {
			mlog.Warn(" *** invalid answer name after MX query for %s\n", m.Question[0].Name)
		}
	}
	return
}
