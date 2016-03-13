package main

import (
	"github.com/grosunick/go-pubsub"
	"github.com/grosunick/go-server"

	"net/http"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	broker := createBroker()
	broker.Run()

	go runTcpServer(broker)
	runWebsocketServer(broker)
}

// PubSub queues brocker factory
func createBroker() pubsub.IBroker {
	return pubsub.NewBroker(
		&pubsub.BrokerProps{
			// период проверки каналов на возможность удаления
			CleanPeriod: pubsub.CHANNELS_WITHOUT_SUBSCRIBERS_CHECK_PERIOD,
			// без подписчиков канал будет жить 3 часа
			LifeTimeWithoutSubscribers: 3600 * 3,
		},
	)
}

// Запускает TcpIp сервер
func runTcpServer(broker pubsub.IBroker) {
	// определяем обработчики json команд, поступающих на вход сервера
	router := server.NewServerRouter()
	router.AddRoute(
		"publish",
		&PublishCommandHandler{broker},
	)

	tcpServer := server.NewSocketServerTcp(
		&server.SocketServerConfig{
			Addr: ":5050",
			MaxWorkerAmount: 200,
			MaxUnhandledRequests: 500,
			ReadTimeOut: 500,
			WriteTimeOut: 500,
		},
		&router,
	)

	tcpServer.Listen()
}

// Запускаем websocket сервер
func runWebsocketServer(broker pubsub.IBroker) {
	// определяем обработчики json команд, поступающих на вход сервера
	router := server.NewServerRouter()
	router.AddRoute(
		"subscribe",
		&SubscribeCommandHandler{broker},
	)

	webSockServer := server.NewWebSocketServer(
		&http.Server{Handler: nil},
		&server.WebSocketServerConfig{
			Addr: ":8080",
			Url: "/ws",
			ReadTimeOut: 60000 * 10, // 10 minutes
			WriteTimeOut: 500,
			ReadBufferSize: 1024,
			WriteBufferSize: 1024,
			CheckOrigin: false,
			ShowLog: true,
		},
		&router,
	)

	webSockServer.Listen()
}
