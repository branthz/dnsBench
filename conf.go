package main

import (
	"dnsBench/log"
	"flag"
	"os"

	"github.com/miekg/dns"
)

var domainSet []string
var P *params
var mlog *log.Logger

type params struct {
	clientsNum int
	recycle     int
	server     string
	port       int
	filePath   string
	DstHost    string
	tp         string
}

func NewParams() *params {
	p := new(params)
	return p
}

var dnstps map[string]uint16

func setdnstps() {
	dnstps = make(map[string]uint16)
	dnstps["A"] = dns.TypeA
	dnstps["CNAME"] = dns.TypeCNAME
}

func init() {
	P = NewParams()
	flag.StringVar(&P.server, "s", "127.0.0.1", "serverAddress")
	flag.IntVar(&P.port, "p", 53, "server port")
	flag.StringVar(&P.filePath, "path", "./domains", "file path for domain sets")
	flag.IntVar(&P.recycle, "n", 1, "run through domains N times")
	flag.IntVar(&P.clientsNum, "c", 1, "client counts")
	flag.StringVar(&P.tp, "t", "A", "dns type")

	var err error
	mlog, err = log.New("", "ERROR")
	if err != nil {
		os.Exit(-1)
	}
	mlog.Infoln("hello world")
	State = NewState()
	setdnstps()
}
