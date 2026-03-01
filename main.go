package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"strings"

	"golang.org/x/net/websocket"
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
	host := "10.88.16.102"
	user := "administrator"
	pass := "abc@123ABCDE"
	port := 53389

	fmt.Printf("Connecting to %s:%d as %s\n", host, port, user)

	settings := &rdpConnectionSettings{
		&host,
		&user,
		&pass,
		int(width),
		int(height),
		port,
	}

	inputq := make(chan inputEvent, 50)
	go rdpconnect(sendq, recvq, inputq, settings)
	go processSendQ(ws, sendq)

	read := make([]byte, 1024)
	for {
		n, err := ws.Read(read)
		if err != nil {
			recvq <- []byte("1")
			return
		}
		if n >= 12 {
			op := binary.LittleEndian.Uint32(read[0:4])
			a := binary.LittleEndian.Uint32(read[4:8])
			b := binary.LittleEndian.Uint32(read[8:12])
			var c uint32
			if n >= 16 {
				c = binary.LittleEndian.Uint32(read[12:16])
			}
			select {
			case inputq <- inputEvent{op, a, b, c}:
			default:
			}
		}
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
