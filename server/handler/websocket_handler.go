package handler

import (
	"app/model"
	"app/repository"
	"fmt"

	"firebase.google.com/go/auth"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
)

// メッセージ用構造体
type Message struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

// メッセージブロードキャストチャネル
var broadcast = make(chan Message)

// 接続されるクライアント
var clients = make(map[*websocket.Conn]bool)

func WebSocketHandler(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()

		token := c.Get("token").(*auth.Token)
		user, _ := repository.UserGetByFirebaseUid(token.UID)

		clients[ws] = true

		message := &Message{
			Username: user.Name,
			Message:  user.Name + "がチャットルームに参加しました",
		}

		// 初回のメッセージを送信
		broadcast <- *message

		listenForWs(ws, user)
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}

func listenForWs(ws *websocket.Conn, user *model.User) {
	for {
		msg := ""
		err := websocket.Message.Receive(ws, &msg)
		if err != nil {
			ws.Close()
			delete(clients, ws)
			break
		} else {
			m := &Message{
				Username: user.Name,
				Message:  msg,
			}
			// 受け取ったメッセージをbroadcastチャネルに送る
			broadcast <- *m
		}
	}
}

func HandleMessages() {
	for {
		// broadcastチャネルからメッセージを受け取る
		message := <-broadcast

		// 接続中の全クライアントにメッセージを送る
		for client := range clients {
			err := websocket.JSON.Send(client, message)
			if err != nil {
				fmt.Println(err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
