package main

import (
	"encoding/json"
	jsonStructs "github.com/grosunick/websocket-server/json"

	"github.com/grosunick/go-pubsub"
	"github.com/grosunick/go-server"
)

/*
	Пример валидных json комманд обработчика
	{"method": "publish", "params": {"channel": "channel_my"}}

	Пример невалидных json комманд обработчика
	{"method": "publish", "params": {}}
*/

// Обработчик подписки на канал
type SubscribeCommandHandler struct {
	broker pubsub.IBroker // брокер подписки
}

// Метод подписки на канал
func (this *SubscribeCommandHandler) Handle(params string, connection server.IConnection) {
	handlerParams, ok := this.GetHandlerParams(params, connection)
	if !ok {
		connection.Close()
		return
	}

	// получаем канал
	channel, ok := this.broker.GetChannel(handlerParams.Channel)
	if !ok {
		channel = this.broker.CreateChannel(
			handlerParams.Channel,
			&pubsub.ChannelProps{
				handlerParams.MaxQueueSize,
			},
		)
	}

	listener := &WebsocketListener{connection}
	channel.GetSubscribers().Add(listener)

	// при закрытии коннекта, удаляем подписку
	connection.(*server.WebSocketConnection).AddOnCloseFunction(
		func () {
			channel.GetSubscribers().Remove(listener)
		},
	)
}

// Возвращает структуру параметров обработчика, инициализированную из json строкой комманды
func (this *SubscribeCommandHandler) GetHandlerParams(
	params string,
	connection server.IConnection,
) (handlerParams jsonStructs.SubscribeCommandParams, res bool) {
	res = false

	err := json.Unmarshal([]byte(params), &handlerParams)
	if err != nil {
		connection.Write("Invalid json for method params \n")
		connection.Write(err.Error())
		return
	}

	if handlerParams.Channel == "" {
		connection.Write("Invalid channel 'param' \n")
		return
	}

	return handlerParams, true
}