package main

import (
	"flag"
	"log"
	"net/http"
	"tinyecho/core"

	"golang.org/x/net/websocket"
)

var (
	pPtr    = flag.String("p", "", "listen port")
	portPtr = flag.String("port", "", "listen port")
)

func main() {
	flag.Parse()

	port := ""

	switch true {
	case *pPtr != "":
		port = *pPtr
	case *portPtr != "":
		port = *portPtr
	default:
		port = "9000"
	}

	http.Handle("/echo", websocket.Handler(core.ServeWS))
	http.Handle("/", websocket.Handler(core.ServeWS))
	http.HandleFunc("/online", core.QueryOnlineClients)
	http.HandleFunc("/publish", core.PublishMessage)
	log.Printf("Listening and Serving on :%v", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
