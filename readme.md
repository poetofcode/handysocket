# HandySocket - удобный адаптер для websocket'ов #

Этот адаптер служит для организации *слушающего* сокета в серверном Go-приложении. В основе данной обёртки лежит библиотека [gorilla-websocket](https://github.com/gorilla/websocket).

## Пример использования ##

Код для сервера `main.go`:

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
		// HandySocket, через которую осуществляется вся работа с веб-сокетом
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
			log.Println("Пришло текстовое сообщение: ", data)
		})

		hs.OnBinaryMessage(func(data []byte) {
			log.Println("Пришло бинарное сообщение размером:", len(data))
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

Код для клиента на JavaScript `index.html`:

```javascript
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>WebSocket Client</title>
</head>
<body>

<h1>WebSocket Client</h1>

<script type="text/javascript">

	webSocket = new WebSocket('ws://localhost:8080/ws');

    webSocket.onopen = function(event) {
        console.log('Сокет открыт');

        // Шлем серверу 'Пинг' каждые 3 секунды
        setInterval(function() {
        	webSocket.send('Пинг');
        }, 3000)
    };

    webSocket.onmessage = function(event) {
        console.log('Сообщение от сервера:', event.data);
    };

    webSocket.onclose = function(event) {
        console.log('Сокет закрыт');
    };

</script>

</body>
</html>
```

Для запуска выполните в терминале:

```bash
go get
go run main.go
```

Для проверки соединений от клиента сохраните приведённую выше страницу `index.html` где-то на диске и запустите в браузере. Если всё работает правильно, то вы увидите сообщения работы сокетов в консоли браузера и в окне терминала.