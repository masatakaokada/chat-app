package handler

import (
	"app/model"
	"app/repository"
	"net/http"
	"strconv"

	"firebase.google.com/go/auth"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

// メッセージ用構造体
type Message struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

// メッセージブロードキャストチャネル
var broadcast = make(chan Message)
var roomBroadcast = make(chan Message)

// 接続されるクライアント
var clients = make(map[*websocket.Conn]interface{})
var roomClients = make(map[*websocket.Conn]interface{})

// wsコネクションの基本設定
var upgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func readMessage(c echo.Context, clients map[*websocket.Conn]interface{}, ws *websocket.Conn, user *model.User, broadcast chan<- Message) {
	defer func() {
		delete(clients, ws)
		ws.Close()
	}()

	for {
		_, byteMessage, err := ws.ReadMessage()
		if err != nil {
			c.Logger().Error(err.Error())
			break
		}
		message := &Message{
			Username: user.Name,
			Message:  string(byteMessage),
		}
		// 受け取ったメッセージをbroadcastチャネルに送る
		broadcast <- *message
	}
}

func writeMessage(clients map[*websocket.Conn]interface{}, broadcast <-chan Message) {
	for {
		// broadcastチャネルからメッセージを受け取る
		message := <-broadcast

		// 接続中の全クライアントにメッセージを送る
		for client := range clients {
			err := client.WriteJSON(message)
			if err != nil {
				delete(clients, client)
				client.Close()
			}
		}
	}
}

func WebSocket(c echo.Context) error {
	// WebSocketsプロトコルにアップグレード
	ws, err := upgradeConnection.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		c.Logger().Error(err.Error())
		return nil
	}

	token := c.Get("token").(*auth.Token)
	user, _ := repository.UserGetByFirebaseUid(token.UID)

	// クライアント登録
	clients[ws] = true

	go writeMessage(clients, broadcast)
	go readMessage(c, clients, ws, user, broadcast)

	return nil
}

func RoomWebSocket(c echo.Context) error {
	// WebSocketsプロトコルにアップグレード
	ws, err := upgradeConnection.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		c.Logger().Error(err.Error())
		return nil
	}

	token := c.Get("token").(*auth.Token)
	user, _ := repository.UserGetByFirebaseUid(token.UID)

	roomId, _ := strconv.Atoi(c.Param("id"))
	roomUser, _ := repository.RoomUserGetByUserIdAndRoomId(user.ID, roomId)

	if roomUser == nil {
		c.Logger().Error(err.Error())
		return nil
	}

	// クライアント登録
	roomClients[ws] = user.ID

	go writeMessage(roomClients, roomBroadcast)
	go readMessage(c, roomClients, ws, user, roomBroadcast)

	return nil
}
