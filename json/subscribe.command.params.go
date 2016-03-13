package json

// Структура, описывающая параметры для обработчика SubscribeCommandHandler
type SubscribeCommandParams struct {
	// название канала
	Channel string `json:"channel"`
	// Максимальный размер очереди не отправленных сообщений в канале
	MaxQueueSize uint8 `json:"maxQueueSize,omitempty"`
}