package json

// Структура, описывающая параметры для обработчика PublishCommandHandler
type PublishCommandParams struct {
	// данные, которые необходимо передать в канал
	Message string `json:"message"`

	//////////// настройки канала ////////////////////////////
	// название канала
	Channel      string `json:"channel"`
	// Максимальный размер очереди не отправленных сообщений в канале
	MaxQueueSize uint8 `json:"maxQueueSize,omitempty"`
}