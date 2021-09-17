package transport

import (
	"fmt"
	"net"
	"time"

	"github.com/miekg/dns"
)

func sendDNSQuery(data []byte, target string) (responses []string, err error) {
	// We use TXT requests to tunnel data. Feel free to implement your own method.
	send := fmt.Sprintf("%s.%s", data, target)
	// fmt.Println("send : ", send)
	responses, err = net.LookupTXT(send)
	// fmt.Println("res : ", responses)

	return
}

func LookupTXT(host string) ([]string, error) {

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
