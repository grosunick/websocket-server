package main

import (
	"github.com/grosunick/go-server"
)

// Слушатель канала брокера подписки, по протоколу websocket
type WebsocketListener struct {
	conn server.IConnection // объект подключения по протоколу websocket
}

// отправляет сообщение подписавшемуся клиенту
func (this *WebsocketListener) Send(message interface{}) {
	this.conn.Write(message.(string))
}