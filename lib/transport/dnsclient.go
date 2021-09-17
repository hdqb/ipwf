package transport

import (
	"context"
	"fmt"
	"net"
)

func sendDNSQuery(data []byte, target string) (responses []string, err error) {
	// We use TXT requests to tunnel data. Feel free to implement your own method.
	send := fmt.Sprintf("%s.%s", data, target)
	// fmt.Println("send : ", send)
	responses, err = LookupTXT(send)
	// fmt.Println("res : ", responses)

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
		panic(err)
	}
	fmt.Println("DNS Result", ipaddr)
	return ipaddr, nil
}

func GoogleDNSDialer(ctx context.Context, network, address string) (net.Conn, error) {
	d := net.Dialer{}
	return d.DialContext(ctx, "udp", "8.8.8.8:53")
}
