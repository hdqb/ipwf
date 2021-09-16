package transport

import (
	"fmt"
	"net"
)

func sendDNSQuery(data []byte, target string) (responses []string, err error) {
	// We use TXT requests to tunnel data. Feel free to implement your own method.
	send := fmt.Sprintf("%s.%s", data, target)
	print(send)
	responses, err = net.LookupTXT(send)
	return
}
