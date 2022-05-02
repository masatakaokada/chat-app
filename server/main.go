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

func handleWebSocket(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()

		// 初回のメッセージを送信
		err := websocket.Message.Send(ws, "Server: Hello, Client!")
		if err != nil {
			c.Logger().Error(err)
		}

		for {
			// Client からのメッセージを読み込む
			msg := ""
			err = websocket.Message.Receive(ws, &msg)
			if err != nil {
				c.Logger().Error(err)
			}

			// Client からのメッセージを元に返すメッセージを作成し送信する
			err := websocket.Message.Send(ws, fmt.Sprintf("サーバー: \"%s\" received!", msg))
			if err != nil {
				c.Logger().Error(err)
				break
			}
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}
