package handysocket

import (
	"github.com/gorilla/websocket"
	"io"
	"net/http"
)

type OnOpenFunc func()
type OnCloseFunc func()
type OnErrorFunc func(err error)
type OnTextMessageFunc func(data string)
type OnBinaryMessageFunc func(data []byte)

type HandySocket struct {
	writer  http.ResponseWriter
	request *http.Request

	// Callbacks
	openCallback          OnOpenFunc
	closeCallback         OnCloseFunc
	errorCallback         OnErrorFunc
	textMessageCallback   OnTextMessageFunc
	binaryMessageCallback OnBinaryMessageFunc

	// Web socket connection
	conn *websocket.Conn

	// Channel for sending msgs to the client
	send chan string
}

func New(w http.ResponseWriter, req *http.Request) *HandySocket {
	hs := &HandySocket{writer: w, request: req}
	return hs
}

func (hs *HandySocket) OnOpen(cb OnOpenFunc) {
	hs.openCallback = cb
}

func (hs *HandySocket) OnTextMessage(cb OnTextMessageFunc) {
	hs.textMessageCallback = cb
}

func (hs *HandySocket) OnBinaryMessage(cb OnBinaryMessageFunc) {
	hs.binaryMessageCallback = cb
}

func (hs *HandySocket) OnClose(cb OnCloseFunc) {
	hs.closeCallback = cb
}

func (hs *HandySocket) OnError(cb OnErrorFunc) {
	hs.errorCallback = cb
}

func (hs *HandySocket) Send(data string) {
	hs.send <- data
}

func (hs *HandySocket) Close() {
	hs.conn.Close()
}

func (hs *HandySocket) Run() {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	conn, err := upgrader.Upgrade(hs.writer, hs.request, nil)
	if err != nil {
		if hs.errorCallback != nil {
			hs.errorCallback(err)
		}
		return
	}
	hs.conn = conn

	// Running go-routine for sending messages to the client
	hs.send = make(chan string)
	go func() {
		for {
			msg := <-hs.send
			hs.conn.WriteMessage(websocket.TextMessage, []byte(msg))
		}
	}()

	if hs.openCallback != nil {
		hs.openCallback()
	}

	// Handling incoming messages from client
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			conn.Close()
			if err == io.EOF && hs.closeCallback != nil {
				hs.closeCallback()
			}

			if err != io.EOF && hs.errorCallback != nil {
				hs.errorCallback(err)
			}
			return
		}
		if messageType == websocket.TextMessage && hs.textMessageCallback != nil {
			hs.textMessageCallback(string(p))
		}
		if messageType == websocket.BinaryMessage && hs.binaryMessageCallback != nil {
			hs.binaryMessageCallback(p)
		}
	}
}
