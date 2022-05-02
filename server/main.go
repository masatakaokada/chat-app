package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/websocket"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"

	_ "github.com/go-sql-driver/mysql" // Using MySQL driver
	"github.com/jmoiron/sqlx"
)

var e = createMux()
var db *sqlx.DB

func main() {
	db = connectDB()

	e.GET("/", hello)
	e.GET("/public", public)
	e.GET("/private", private, firebaseMiddleware())
	e.GET("/ws", handleWebSocket)
	go handleMessages()

	e.Logger.Fatal(e.Start(":8082"))
}

func createMux() *echo.Echo {
	e := echo.New()

	e.Use(middleware.CORS())
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.Gzip())

	return e
}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!\n")
}

func public(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, public!\n")
}

func private(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, private!\n")
}

func connectDB() *sqlx.DB {
	dsn := os.Getenv("DSN")
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		e.Logger.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		e.Logger.Fatal(err)
	}
	log.Println("db connection succeeded")
	return db
}

// JWTを検証する
func firebaseMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Firebase SDK のセットアップ
			opt := option.WithCredentialsFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
			app, err := firebase.NewApp(context.Background(), nil, opt)
			if err != nil {
				fmt.Printf("error: %v\n", err)
				os.Exit(1)
			}

			client, err := app.Auth(context.Background())
			if err != nil {
				fmt.Printf("error: %v\n", err)
				os.Exit(1)
			}

			// クライアントから送られてきた JWT 取得
			auth := c.Request().Header.Get("Authorization")
			idToken := strings.Replace(auth, "Bearer ", "", 1)

			// JWT の検証
			token, err := client.VerifyIDToken(context.Background(), idToken)
			if err != nil {
				fmt.Printf("error verifying ID token: %v\n", err)
				return err
			}

			c.Set("token", token)
			return next(c)
		}
	}
}

// 接続されるクライアント
var clients = make(map[*websocket.Conn]bool)

// メッセージブロードキャストチャネル
var broadcast = make(chan Message)

// メッセージ用構造体
type Message struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

func handleMessages() {
	for {
		// broadcastチャネルからメッセージを受け取る
		message := <-broadcast

		// 接続中の全クライアントにメッセージを送る
		for client := range clients {
			err := websocket.Message.Send(client, message.Message)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func handleWebSocket(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()

		// クライアントを登録
		clients[ws] = true

		message := &Message{
			Username: "サーバーサイド",
			Message:  "チャットルームに参加しました",
		}

		// 初回のメッセージを送信
		broadcast <- *message

		for {
			msg := ""
			err := websocket.Message.Receive(ws, &msg)
			if err != nil {
				ws.Close()
				c.Logger().Error(err)
			}

			m := &Message{
				Username: "ユーザー",
				Message:  msg,
			}

			// 受け取ったメッセージをbroadcastチャネルに送る
			broadcast <- *m
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}
