package main

import (
	"os"
	"io/ioutil"
	"strings"
	"errors"
)

type Result []string

func readFile(path string)(Result,error){
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	content, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return parseFile(string(content))
}

func parseFile(data string) (Result,error){
	var line string
	var rst Result
	lines:=strings.Split(data,"\n")
	for i:=0;i<len(lines);i++{
		line = strings.Trim(lines[i], "\r\t\n ")
		if line == "" || line[0] == '#' {
			continue
		}
		rst=append(rst,line)
	}
	if len(rst) <1 {
		return nil,errors.New("not domains find")
	}
	return rst,nil
}
