package main

import (
	"flag"
	"fmt"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"strings"
)

var (
	hostname = flag.String("h", "", "The RDP server we will connect to")
	username = flag.String("u", "", "Username for the RDP server")
	password = flag.String("p", "", "Password for the RDP server")
)

func getResolution(ws *websocket.Conn) (width int64, height int64) {
	request := ws.Request()
	dtsize := request.FormValue("dtsize")

	if !strings.Contains(dtsize, "x") {
		width = 800
		height = 600
	} else {
		sizeparts := strings.Split(dtsize, "x")

		width, _ = strconv.ParseInt(sizeparts[0], 10, 32)
		height, _ = strconv.ParseInt(sizeparts[1], 10, 32)

		if width < 400 {
			width = 400
		} else if width > 1920 {
			width = 1920
		}

		if height < 300 {
			height = 300
		} else if height > 1080 {
			height = 1080
		}
	}

	return width, height
}

func processSendQ(ws *websocket.Conn, sendq chan []byte) {
	for {
		buf := <-sendq
		err := websocket.Message.Send(ws, buf)
		if err != nil {
			panic("ListenAndServe: " + err.Error())
		}
	}
}

func initSocket(ws *websocket.Conn) {
	sendq := make(chan []byte, 100)
	recvq := make(chan []byte, 5)

	width, height := getResolution(ws)
	fmt.Printf("User requested size %d x %d\n", width, height)

	// Get connection parameters from WebSocket request
	request := ws.Request()
	host := request.FormValue("host")
	user := request.FormValue("user")
	pass := request.FormValue("pass")
	portStr := request.FormValue("port")

	// Use command line parameters as fallback
	if host == "" {
		host = *hostname
	}
	if user == "" {
		user = *username
	}
	if pass == "" {
		pass = *password
	}

	// Parse port, default to 3389
	port := 3389
	if portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			port = p
		}
	}

	fmt.Printf("Connecting to %s:%d as %s\n", host, port, user)

	settings := &rdpConnectionSettings{
		&host,
		&user,
		&pass,
		int(width),
		int(height),
		port,
	}

	go rdpconnect(sendq, recvq, settings)
	go processSendQ(ws, sendq)

	read := make([]byte, 1024, 1024)
	for {
		_, err := ws.Read(read)
		if err != nil {
			recvq <- []byte("1")
		}

		recvq <- read
		log.Println(string(read))
	}
}

func main() {
	flag.Parse()

	runtime.GOMAXPROCS(runtime.NumCPU())

	// WebSocket handler for RDP connection
	http.Handle("/ws", websocket.Handler(initSocket))

	// Static file server for webroot
	http.Handle("/", http.FileServer(http.Dir("webroot")))

	fmt.Printf("请访问: http://localhost:%d/index-debug.html\n", 4455)
	err := http.ListenAndServe(":4455", nil)
	if err != nil {
		panic("ListenANdServe: " + err.Error())
	}
}
