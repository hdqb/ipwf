package main

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/miekg/dns"
)

func GetFreePort() (port int, err error) {
	var a *net.TCPAddr
	if a, err = net.ResolveTCPAddr("udp", "localhost:0"); err == nil {
		var l *net.TCPListener
		if l, err = net.ListenTCP("udp", a); err == nil {
			defer l.Close()
			return l.Addr().(*net.TCPAddr).Port, nil
		}
	}
	return
}
func main() {
	// var txts []string
	m1 := new(dns.Msg)
	m1.Id = dns.Id()
	m1.RecursionDesired = true
	m1.Question = make([]dns.Question, 1)
	m1.Question[0] = dns.Question{"api.vimmo.app.", dns.TypeTXT, dns.ClassINET}
	c := new(dns.Client)
	fp, err := GetFreePort()

	laddr := net.UDPAddr{
		IP:   net.ParseIP("[::1]"),
		Port: fp,
		Zone: "",
	}

	c.Dialer = &net.Dialer{
		Timeout:   200 * time.Millisecond,
		LocalAddr: &laddr,
	}
	txt, rtt, err := c.Exchange(m1, "8.8.8.8:53")

	if err != nil {
		// return nil, &DNSError{
		// 	Err:    "cannot unmarshal DNS message",
		// 	Name:   name,
		// 	Server: server,
		// }
	}
	// Multiple strings in one TXT record need to be
	// concatenated without separator to be consistent
	// with previous Go resolver.
	// n := 0
	for _, s := range txt.Answer {
		fmt.Printf("%v %v ", s, rtt)
		// n += len(s)
	}
	// txtJoin := make([]byte, 0, n)
	// for _, s := range txt.Answer {
	// 	// txtJoin = append(txtJoin, s...)
	// }
	// if len(txts) == 0 {
	// 	txts = make([]string, 0, 1)
	// }
	// txts = append(txts, string(txtJoin))

}

func mains() {

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
	fp, err := GetFreePort()
	if err != nil {
		// return nil, err
	}
	laddr := net.UDPAddr{
		IP:   net.ParseIP("[::1]"),
		Port: fp,
		Zone: "",
	}

	d := &net.Dialer{
		Timeout:   200 * time.Millisecond,
		LocalAddr: &laddr,
	}
	return d.DialContext(ctx, "udp", "8.8.8.8:53")
}
