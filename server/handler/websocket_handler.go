package handler

import (
	"app/model"
	"app/repository"
	"fmt"
	"log"
	"strconv"

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

// 接続されるクライアント
var roomClients = make(map[*websocket.Conn]int)

// メッセージブロードキャストチャネル
var roomBroadcast = make(chan Message)

func handleRoomMessages() {
	for {
		// broadcastチャネルからメッセージを受け取る
		message := <-roomBroadcast

		// 接続中の全クライアントにメッセージを送る
		for client := range roomClients {
			err := websocket.JSON.Send(client, message)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(roomClients, client)
			}
		}
	}
}

func RoomWebSocket(c echo.Context) error {
	go handleRoomMessages()
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()

		token := c.Get("token").(*auth.Token)

		user, _ := repository.UserGetByFirebaseUid(token.UID)

		roomId, _ := strconv.Atoi(c.Param("id"))

		roomUser, _ := repository.RoomUserGetByUserIdAndRoomId(user.ID, roomId)

		if roomUser != nil {
			roomClients[ws] = user.ID

			message := &Message{
				Username: user.Name,
				Message:  "チャットルームに参加しました",
			}

			// 初回のメッセージを送信
			roomBroadcast <- *message

			for {
				msg := ""
				err := websocket.Message.Receive(ws, &msg)
				if err != nil {
					ws.Close()
					delete(roomClients, ws)
					break
				} else {
					m := &Message{
						Username: user.Name,
						Message:  msg,
					}
					// 受け取ったメッセージをbroadcastチャネルに送る
					roomBroadcast <- *m
				}
			}
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}
