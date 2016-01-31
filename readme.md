# HandySocket - удобный адаптер для websocket'ов #

В основе данной обёртки лежит библиотека [gorilla-websocket](https://github.com/gorilla/websocket)

## Пример использования ##

```go
package main

import (
	"github.com/poetofcode/handysocket"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/ws", func(w http.ResponseWriter, req *http.Request) {
		//
		// Для создания просто вызываем метод New, который создаёт указатель на структуру
		// HandySocket, через которой осуществляется вся работа с веб-сокетом
		//
		hs := handysocket.New(w, req)

		//
		// Для каждого события или действия в handysocket предусмотрены callback-методы
		// Использование их приведено в коде ниже 
		//
		hs.OnOpen(func() {
			log.Println("Сокет открыт")
			hs.Send("Привет от сервера!")
		})

		hs.OnTextMessage(func(data string) {
			log.Println("Пришло текстовое сообщение: " + data)
		})

		hs.OnBinaryMessage(func(data []byte) {
			log.Println("Пришло бинарное сообщение размером:" + len(data))
		})

		hs.OnError(func(err error) {
			log.Println("Произошёл сбой", err)
		})

		hs.OnClose(func() {
			log.Println("Сокет закрыт")
		})

		hs.Run()
	})

	panic(http.ListenAndServe(":8080", nil))
}
```