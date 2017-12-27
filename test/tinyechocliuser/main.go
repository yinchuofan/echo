package main

import (
	"bufio"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/websocket"
)

// Message :
var Message = `{"code":"MsgCodeAlarm","content":"alarm","publisher":"user","time":"2016-01-01 10:10:10"}`

func getMessage(input string) string {
	cmdList := strings.Split(input, " ")
	if len(cmdList) != 2 {
		return ""
	}

	data := ""
	if cmdList[0] == "sub" {
		channels := cmdList[1]
		data = "command:subscribe\r\n\r\nchannels:" + channels
	} else if cmdList[0] == "unsub" {
		channels := cmdList[1]
		data = "command:unsubscribe\r\n\r\nchannels:" + channels
	} else if cmdList[0] == "pub" {
		channel := cmdList[1]
		data = "command:publish\r\n\r\nchannel:" + channel + "\r\nmessage:" + Message
	} else {
		return ""
	}

	return data
}

var (
	cliType = "app"
	cliID   = "4ftsc8bgeird6sb8"
	cliSec  = MD5("5ece183a0b274cb9992b503183e6494d")
)

var MessageHeaderLength int = 32

type Header struct {
	Flag          uint32
	Length        uint32
	Checksum      uint32
	Version       uint32
	CommandCode   uint32
	ErrorCode     uint32
	TxtDataLength uint32
	BinDataLength uint32
}

func NewHeader(buf []byte) *Header {
	header := new(Header)

	header.Flag = binary.BigEndian.Uint32(buf[0:4])
	header.Length = binary.BigEndian.Uint32(buf[4:8])
	header.Checksum = binary.BigEndian.Uint32(buf[8:12])
	header.Version = binary.BigEndian.Uint32(buf[12:16])
	header.CommandCode = binary.BigEndian.Uint32(buf[16:20])
	header.ErrorCode = binary.BigEndian.Uint32(buf[20:24])
	header.TxtDataLength = binary.BigEndian.Uint32(buf[24:28])
	header.BinDataLength = binary.BigEndian.Uint32(buf[28:32])

	return header
}

// Now :
func Now() time.Time {
	now := time.Now()
	_, s := now.Zone()
	now = now.Add(time.Second * time.Duration(s))
	return now
}

func readMessage(ws *websocket.Conn) bool {
	log.Println("read data begin")

	//read header
	var headBuf = make([]byte, MessageHeaderLength)
	if false == readData(ws, MessageHeaderLength, headBuf, len(headBuf)) {
		log.Println("read header data failed")
		return false
	}

	log.Printf("headBuf: %v, size: %v", headBuf, len(headBuf))

	header := NewHeader(headBuf)
	log.Printf("flag: %v, length: %v, checksum: %v, version: %v, command: %v, error: %v, txtdatalen: %v, bindatalen: %v",
		header.Flag, header.Length, header.Checksum, header.Version, header.CommandCode, header.ErrorCode, header.TxtDataLength, header.BinDataLength)

	//read text data
	var txtBuf = make([]byte, header.TxtDataLength)
	if false == readData(ws, int(header.TxtDataLength), txtBuf, len(txtBuf)) {
		log.Println("read text data failed")
		return false
	}

	log.Println("text data: \r\n", string(txtBuf))

	//read binary data
	var binBuf = make([]byte, header.BinDataLength)
	if false == readData(ws, int(header.BinDataLength), binBuf, len(binBuf)) {
		log.Println("read binary data failed")
		return false
	}

	var jpgName = "F:/Data/Picture/720p-" + strconv.FormatInt(Now().UnixNano(), 10) + ".jpg"
	err := ioutil.WriteFile(jpgName, binBuf, 0666)
	if err != nil {
		log.Println("write jpg file failed")
	}

	log.Println("read data succeed")
	return true
}

func readData(ws *websocket.Conn, readSize int, buffer []byte, bufferSize int) bool {
	log.Printf("readSize: %v, bufferSize: %v", readSize, bufferSize)
	if bufferSize < readSize {
		log.Println("buffer is not enough")
		return false
	}

	var ReadSizeOnce = 1024
	var remainSize = readSize
	var bufPos = 0
	for remainSize > 0 {
		var onceSize = ReadSizeOnce
		if remainSize < onceSize {
			onceSize = remainSize
		}

		n, err := ws.Read(buffer[bufPos : bufPos+onceSize])
		if err != nil {
			return false
		}
		if n <= 0 {
			return false
		}

		bufPos += n
		remainSize -= n
	}
	//log.Printf("buffer: %v, size: %v", buffer, len(buffer))
	return true
}

func main() {
	origin := "http://localhost"
	url := "ws://localhost:9000/echo"
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	go func() {
		//message := []string{}
		for {
			readMessage(ws)
		}
	}()

	inputReader := bufio.NewReader(os.Stdin)

	for {
		log.Println("input command")

		input, _, err := inputReader.ReadLine()
		if err != nil {
			break
		}

		if string(input) == "exit" {
			log.Println("exit")
			break
		}

		message := getMessage(string(input))
		if message == "" {
			log.Println("invalid message")
			continue
		}

		if _, err := ws.Write([]byte(message)); err != nil {
			log.Fatal(err)
			continue
		}
		log.Printf("sent message:\n%v", message)
	}

}

// MD5 : MD5
func MD5(input string) string {
	dataMD5 := md5.Sum([]byte(input))
	dataMD5Slice := []byte{}
	for _, v := range dataMD5 {
		dataMD5Slice = append(dataMD5Slice, v)
	}
	output := hex.EncodeToString(dataMD5Slice)
	return output
}
