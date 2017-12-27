package core

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/websocket"
)

var (
	mapOnlineClients = make(map[string]time.Time)
	mapCliMutex      = new(sync.Mutex)
)

// WsHandler websocket handler
type WsHandler struct {
	ws *websocket.Conn

	channels []string
	clientID string

	message *WsMessage
}

func (ws *WsHandler) online() {
	if len(ws.clientID) > 0 {
		log.Printf("online client id: %v", ws.clientID)
		mapCliMutex.Lock()
		mapOnlineClients[ws.clientID] = time.Now()
		mapCliMutex.Unlock()
	}
}

func (ws *WsHandler) offline() {
	if len(ws.clientID) > 0 {
		log.Printf("offline client id: %v", ws.clientID)
		mapCliMutex.Lock()
		if _, ok := mapOnlineClients[ws.clientID]; ok {
			delete(mapOnlineClients, ws.clientID)
		}
		mapCliMutex.Unlock()
	}
}

// Notify notify client
func (ws *WsHandler) Notify(message string) {
	log.Printf("notify message to client: %v", ws)
	if nil != websocket.Message.Send(ws.ws, message) {
		log.Printf("notify message failed. client: %v", ws)
	}
}

// WsMessage ws message
type WsMessage struct {
	rawdata string

	command string
	body    string
}

func parseMessage(message string) (wsMsg *WsMessage) {
	i := strings.Index(message, "\r\n\r\n")

	if i < 0 {
		return nil
	}

	wsMsg = &WsMessage{rawdata: message}

	command := message[0:i]
	strList := strings.Split(command, ":")
	if len(strList) == 2 {
		wsMsg.command = strings.TrimSpace(strList[1])
	} else {
		return nil
	}

	wsMsg.body = message[i+4:]

	return wsMsg
}

// ServeWS start server
func ServeWS(ws *websocket.Conn) {
	log.Println("connection: ", ws.RemoteAddr().String())

	handler := &WsHandler{ws: ws}

	for {
		var data string
		err := websocket.Message.Receive(ws, &data)
		if err != nil {
			log.Println("peer error or close websocket")
			break
		}

		log.Print("receive message:\n", data)

		wsMsg := parseMessage(data)
		if wsMsg == nil {
			log.Printf("parse message error. message: %v", data)
			continue
		}

		handler.message = wsMsg

		processor := getProcessor(handler.message.command)
		if processor != nil {
			go processor(handler)
		}
	}

	handler.offline()

	for _, v := range handler.channels {
		Unsubscribe(v, handler)
	}

	log.Println("server close websocket")
	ws.Close()
}

// QueryOnlineClients query online clients
func QueryOnlineClients(w http.ResponseWriter, r *http.Request) {
	clientIDListString := r.URL.Query()["clientid"]

	log.Println("http query online clients: ", clientIDListString)

	onlineClientIDList := []string{}
	mapCliMutex.Lock()
	for _, v := range clientIDListString {
		if _, ok := mapOnlineClients[v]; ok {
			onlineClientIDList = append(onlineClientIDList, v)
		}
	}
	mapCliMutex.Unlock()

	onlineClientIDListString := strings.Join(onlineClientIDList, ",")
	log.Println("online client list: ", onlineClientIDListString)
	w.Write([]byte(onlineClientIDListString))
}

// PublishMessage publish message
func PublishMessage(w http.ResponseWriter, r *http.Request) {
	/*
		URL: /publish?channel=channel1&channel=channel2
		  ---------------HTTP Body Format-------------------
		  	----------------------------------------------------------------------------------------------------------------------
			|											PROTOCOL HEADER										|   PROTOCOL BODY	 |
			----------------------------------------------------------------------------------------------------------------------
			| FLAG | LENGTH | CHECKSUM | VERSION | COMMANDCODE | ERRORCODE | TEXTDATALENGTH | BINDATALENGTH | TEXTDATA | BINDATA |
			----------------------------------------------------------------------------------------------------------------------
			|  4B  |   4B   |    4B    |    4B   |     4B      |     4B    |       4B       |      4B       |  Unknown | Unknown |
			----------------------------------------------------------------------------------------------------------------------
		  ------------------------------------------------
	*/
	channels := r.URL.Query()["channel"]
	log.Print("http publish message to channels:", channels)
	if len(channels) == 0 {
		w.Write([]byte("error"))
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("error"))
	}

	for _, v := range channels {
		v = strings.TrimSpace(v)
		Publish(v, string(body))
	}

	w.Write([]byte("success"))
}
