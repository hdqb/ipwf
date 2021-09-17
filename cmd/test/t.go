package main

import (
	"fmt"
	"net"
	"time"

	"github.com/miekg/dns"
)

func main() {

	m1 := new(dns.Msg)
	m1.Id = dns.Id()
	m1.RecursionDesired = true
	m1.Question = make([]dns.Question, 1)
	m1.Question[0] = dns.Question{"2d1b385858c0f3aa26cd0a55dea00ce2a6cc56aa639a025bc9175ab50634b5c.13eff45d1fa8f61b4390ad84bb248151239a46ccbb33f6370.c.vimmo.app.", dns.TypeTXT, dns.ClassINET}
	c := new(dns.Client)
	laddr := net.UDPAddr{
		IP:   net.ParseIP("[::1]"),
		Port: 12345,
		Zone: "",
	}

	c.Dialer = &net.Dialer{
		Timeout:   200 * time.Millisecond,
		LocalAddr: &laddr,
	}
	in, rtt, err := c.Exchange(m1, "8.8.8.8:53")

	fmt.Println(in, rtt, err)

}
