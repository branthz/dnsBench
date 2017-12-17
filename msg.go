package main

import "github.com/miekg/dns"

func NewMsg(site string)(m *dns.Msg){
	m=new(dns.Msg)
	m.SetQuestion(dns.Fqdn(site), dns.TypeA)
	m.RecursionDesired = true
	return
}

