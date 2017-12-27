package core

import (
	"log"
	"strings"
)

var mapProcessor map[string]Processor

// Processor processor
type Processor func(ws *WsHandler)

func processRegister(ws *WsHandler) {
	/*
	  ---------------Message Format-------------------
	  command : register \r\n
	  \r\n
	  clientid : string
	  ------------------------------------------------
	*/
	log.Println("process register")

	strList := strings.Split(ws.message.body, ":")
	if len(strList) == 2 {
		ws.clientID = strings.TrimSpace(strList[1])
		ws.online()
	}
}

func processSubscribe(ws *WsHandler) {
	/*
	  ---------------Message Format-------------------
	  command : subscribe \r\n
	  \r\n
	  channels : channel1,channel2,...
	  ------------------------------------------------
	*/
	log.Println("process subscribe")

	strList := strings.Split(ws.message.body, ":")
	if len(strList) == 2 {
		channels := strings.TrimSpace(strList[1])
		ws.channels = strings.Split(channels, ",")

		for _, v := range ws.channels {
			Subscribe(v, ws)
		}
	}
}

func processUnsubscribe(ws *WsHandler) {
	/*
	  ---------------Message Format-------------------
	  command : unsubscribe \r\n
	  \r\n
	  channels : channel1,channel2,...
	  ------------------------------------------------
	*/
	log.Println("process unsubscribe")

	strList := strings.Split(ws.message.body, ":")
	if len(strList) == 2 {
		channels := strings.TrimSpace(strList[1])
		ws.channels = strings.Split(channels, ",")

		for _, v := range ws.channels {
			Unsubscribe(v, ws)
		}
	}
}

func processPublish(ws *WsHandler) {
	/*
		  ---------------Message Format-------------------
		  command : publish \r\n
		  \r\n
		  channels : channel1,channel2,... \r\n
		  	----------------------------------------------------------------------------------------------------------------------
			|											PROTOCOL HEADER										|   PROTOCOL BODY	 |
			----------------------------------------------------------------------------------------------------------------------
			| FLAG | LENGTH | CHECKSUM | VERSION | COMMANDCODE | ERRORCODE | TEXTDATALENGTH | BINDATALENGTH | TEXTDATA | BINDATA |
			----------------------------------------------------------------------------------------------------------------------
			|  4B  |   4B   |    4B    |    4B   |     4B      |     4B    |       4B       |      4B       |  Unknown | Unknown |
			----------------------------------------------------------------------------------------------------------------------
		  ------------------------------------------------
	*/
	i := strings.Index(ws.message.body, "\r\n")
	channelsLine := ws.message.body[:i]

	strList := strings.Split(channelsLine, ":")

	if len(strList) == 2 {
		channels := strings.TrimSpace(strList[1])
		channelList := strings.Split(channels, ",")

		message := ws.message.body[i+2:]

		for _, v := range channelList {
			Publish(v, message)
		}
	}
}

func getProcessor(command string) Processor {
	if elem, ok := mapProcessor[command]; ok {
		return elem
	}
	return nil
}

func init() {
	mapProcessor = make(map[string]Processor)
	mapProcessor["register"] = processRegister
	mapProcessor["subscribe"] = processSubscribe
	mapProcessor["unsubscribe"] = processUnsubscribe
	mapProcessor["publish"] = processPublish
}
