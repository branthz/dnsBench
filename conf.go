package main

import (
	"fmt"
	"os"

	"github.com/branthz/utarrow/lib/log"
	"github.com/urfave/cli"

	"github.com/miekg/dns"
)

var app *cli.App
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
	app = cli.NewApp()
	app.Author = "brant"
	app.Name = "dnsBench"
	app.Version = "0.0.2"
	app.Description = "an dns server stree test tool"
	app.HideHelp = true

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "s",
			Value:       "8.8.8.8",
			Usage:       "dns server address",
			Destination: &P.server,
		},
		cli.IntFlag{
			Name:        "p",
			Value:       53,
			Usage:       "server port",
			Destination: &P.port,
		},
		cli.StringFlag{
			Name:        "path",
			Value:       "./domains",
			Usage:       "file path of  domain sets",
			Destination: &P.filePath,
		},
		cli.StringFlag{
			Name:        "t",
			Value:       "A",
			Usage:       "msg type:(A,CNAME)",
			Destination: &P.tp,
		},
		cli.IntFlag{
			Name:        "n",
			Value:       1,
			Usage:       "run through domains N times",
			Destination: &P.recycle,
		},
	}
	app.Action = func(c *cli.Context) error {
		fmt.Printf("app run on version:%s\n", app.Version)
		return nil
	}
	var err error
	if err = app.Run(os.Args); err != nil {
		os.Exit(-1)
	}

	Initdnstps()

	if err = P.CheckInput(); err != nil {
		fmt.Printf("input params not fit:%v\n", err)
		os.Exit(-1)
	}

	mlog, err = log.New("", log.Error)
	if err != nil {
		os.Exit(-1)
	}
	mlog.Infoln("dnsBench start ...")
	State = NewState()

	domainSet, err = readFile(P.filePath)
	if err != nil || len(domainSet) == 0 {
		mlog.Errorln(err)
		os.Exit(-1)
	}
	mlog.Debug("%+v\n", domainSet[0])
}
