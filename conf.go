package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/branthz/utarrow/lib/log"

	"github.com/miekg/dns"
)

var domainSet []string
var P *params
var mlog *log.Logger

type params struct {
	clientsNum int
	recycle    int
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

func (p *params) Parse() {
	flag.StringVar(&P.server, "s", "127.0.0.1", "serverAddress")
	flag.IntVar(&P.port, "p", 53, "server port")
	flag.StringVar(&P.filePath, "path", "./domains", "file path for domain sets")
	flag.IntVar(&P.recycle, "n", 1, "run through domains N times")
	//flag.IntVar(&P.clientsNum, "c", 1, "client counts")
	flag.StringVar(&P.tp, "t", "A", "dns type,support A/CNAME")
	return
}

func (p *params) CheckInput() error {
	if _, ok := dnstps[p.tp]; !ok {
		return fmt.Errorf("not support type:%s", p.tp)
	}
	return nil
}

var dnstps map[string]uint16

func Initdnstps() {
	dnstps = make(map[string]uint16)
	dnstps["A"] = dns.TypeA
	dnstps["CNAME"] = dns.TypeCNAME
}

func init() {
	P = NewParams()
	P.Parse()
	if len(os.Args) < 2 {
		flag.PrintDefaults()
		os.Exit(-1)
	}
	flag.Parse()
	Initdnstps()

	var err error
	if err = P.CheckInput(); err != nil {
		fmt.Printf("input params not fit:%v\n", err)
		os.Exit(-1)
	}

	mlog, err = log.New("", "ERROR")
	if err != nil {
		os.Exit(-1)
	}
	mlog.Infoln("dnsBench start ...")
	State = NewState()
}
