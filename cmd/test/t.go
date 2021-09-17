package main

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/miekg/dns"
)

func LookupTXT(host string) {

	m1 := new(dns.Msg)
	m1.Id = dns.Id()
	m1.RecursionDesired = true
	m1.Question = make([]dns.Question, 1)
	m1.Question[0] = dns.Question{host, dns.TypeTXT, dns.ClassINET}
	c := new(dns.Client)
	laddr := net.UDPAddr{
		IP:   net.ParseIP("[::1]"),
		Port: 0,
		Zone: "",
	}

	c.Dialer = &net.Dialer{
		Timeout:   200 * time.Millisecond,
		LocalAddr: &laddr,
	}
	in, rtt, err := c.Exchange(m1, "8.8.8.8:53")

	fmt.Println(in, rtt, err)

}

func main() {

	r := net.Resolver{
		PreferGo: true,

		Dial: GoogleDNSDialer,
	}
	ctx := context.Background()
	ipaddr, err := r.LookupTXT(ctx, "www.example.com")
	if err != nil {
		panic(err)
	}
	fmt.Println("DNS Result", ipaddr)
}
func GoogleDNSDialer(ctx context.Context, network, address string) (net.Conn, error) {
	d := net.Dialer{}
	return d.DialContext(ctx, "udp", "8.8.8.8:53")
}
