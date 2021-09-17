package transport

import (
	"context"
	"fmt"
	"net"
	"time"
)

func sendDNSQuery(data []byte, target string) (responses []string, err error) {
	// We use TXT requests to tunnel data. Feel free to implement your own method.
	send := fmt.Sprintf("%s.%s", data, target)
	fmt.Println("send : ", send)
	responses, err = LookupTXT(send)
	fmt.Println("res : ", responses)

	return
}

func LookupTXT(send string) ([]string, error) {

	r := net.Resolver{
		PreferGo: true,

		Dial: GoogleDNSDialer,
	}
	ctx := context.Background()
	ipaddr, err := r.LookupTXT(ctx, send)
	if err != nil {
		// panic(err)

	}
	// fmt.Println("DNS Result", ipaddr)
	return ipaddr, nil
}

func GetFreePort() (port int, err error) {
	addr, err := net.ResolveTCPAddr("udp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("udp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

func GoogleDNSDialer(ctx context.Context, network, address string) (net.Conn, error) {
	// fp, err := GetFreePort()
	// if err != nil {
	// return nil, err
	// }
	laddr := net.UDPAddr{
		IP:   net.ParseIP("[::1]"),
		Port: 24123,
		Zone: "",
	}

	d := &net.Dialer{
		Timeout:   200,
		Deadline:  time.Now(),
		LocalAddr: &laddr,
		// KeepAlive: time.Duration(864000) * time.Millisecond,
	}
	return d.DialContext(ctx, "udp", "8.8.8.8:53")
}
