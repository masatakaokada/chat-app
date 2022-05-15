package main

import (
	"app/handler"
	"app/model"
	"app/repository"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"

	_ "github.com/go-sql-driver/mysql" // Using MySQL driver
	"github.com/jmoiron/sqlx"
)

var e = createMux()
var db *sqlx.DB

func main() {
	db = connectDB()
	repository.SetDB(db)

	e.GET("/", hello)
	e.GET("/ws", handler.WebSocket, firebaseMiddleware())
	e.GET("/ws/rooms/:id", handler.RoomWebSocket, firebaseMiddleware())
	e.GET("/room-creation-users", handler.RoomCreationUsers, firebaseMiddleware())
	e.POST("/users", user, firebaseMiddleware())
	e.GET("/rooms", handler.RoomIndex, firebaseMiddleware())
	e.POST("/rooms", handler.RoomCreate, firebaseMiddleware())

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!\n")
}

func createMux() *echo.Echo {
	e := echo.New()

	e.Use(middleware.CORS())
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.Gzip())

	return e
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

			if idToken == "" {
				idToken = c.QueryParam("token")
			}

			// JWT の検証
			token, err := client.VerifyIDToken(context.Background(), idToken)
			if err != nil {
				fmt.Printf("error verifying ID token: %v\n", err)
				return c.String(http.StatusUnauthorized, "トークンの有効期限切れです")
			}

			c.Set("token", token)
			return next(c)
		}
	}
}

func user(c echo.Context) error {
	token := c.Get("token").(*auth.Token)

	user, _ := repository.UserGetByFirebaseUid(token.UID)

	if user == nil {
		user := &model.User{
			Email:       token.Claims["email"].(string),
			FirebaseUid: token.UID,
		}

		if err := c.Bind(&user); err != nil {
			c.Logger().Error(err.Error())

			return c.NoContent(http.StatusBadRequest)
		}

		res, err := repository.UserCreate(user)
		if err != nil {
			c.Logger().Error(err.Error())

			return c.NoContent(http.StatusInternalServerError)
		}

		id, _ := res.LastInsertId()
		fmt.Printf("ユーザーの作成に成功しました。 id: %v\n", id)
	}

	return c.NoContent(http.StatusOK)
}
