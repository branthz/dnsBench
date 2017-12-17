package main

import (
	"flag"
	"os"
)

func main(){
	if len(os.Args)<2 {
		flag.PrintDefaults()
		return
	}
	flag.Parse()

	var err error
	domainSet,err=readFile(P.filePath)
	if err!=nil||len(domainSet)==0{
		mlog.Errorln(err)
		return
	}
	mlog.Debug("%+v\n",domainSet[0])

	Client,err:=NewClient()
	if err !=nil{
		mlog.Errorln(err)
		return
	}
	State.Start()
	go Client.Query(domainSet)
	Client.Response()
	State.Show()
}