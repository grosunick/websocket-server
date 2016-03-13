package main

import (
	"encoding/json"
	jsonStructs "github.com/grosunick/websocket-server/json"

	"github.com/grosunick/go-pubsub"
	"github.com/grosunick/go-server"
)

/*
	Пример валидных json комманд обработчика

	{"method": "publish", "params": {"message": "Hello everybody!!!", "channel": "channel_my"}}
	{"method": "publish", "params": {"message": "Hello everybody!!!", "channel": "channel_my", "maxQueueSize": 150}}

	Пример невалидных json комманд обработчика

	{"method": "publish", "params": {"message": "Hello everybody!!!"}}
	{"method": "publish", "params": {"channel": "channel_my"}}
*/

// Обработчик публикации сообщения в канал
type PublishCommandHandler struct {
	broker pubsub.IBroker // брокер подписки
}

// Метод публикации сообщения
func (this *PublishCommandHandler) Handle(params string, connection server.IConnection) {
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

	// публикуем сообщение в канал
	channel.Publish(handlerParams.Message)
	connection.Write(`{"response": "success"}` + "\n")
	connection.Close()
}

// Возвращает структуру параметров обработчика, инициализированную из json строкой комманды
func (this *PublishCommandHandler) GetHandlerParams(
	params string,
	connection server.IConnection,
) (handlerParams jsonStructs.PublishCommandParams, res bool) {
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

	if handlerParams.Message == "" {
		connection.Write("Invalid message 'param' \n")
		return
	}

	return handlerParams, true
}
