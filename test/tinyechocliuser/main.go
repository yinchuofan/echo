package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"log"
	"os"
	"strings"

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
			var buffer = make([]byte, 512)
			var n int
			if n, err = ws.Read(buffer); err != nil {
				log.Fatal(err)
				continue
			}
			log.Printf("received buffer: %s. size: %v\n", buffer[:n], n)
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
