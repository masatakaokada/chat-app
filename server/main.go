package main

import (
	"app/auth"
	"app/handler"
	"app/repository"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gopkg.in/go-playground/validator.v9"

	_ "github.com/go-sql-driver/mysql" // Using MySQL driver
	"github.com/jmoiron/sqlx"
)

var e = createMux()
var db *sqlx.DB

func main() {
	db = connectDB()
	repository.SetDB(db)

	e.GET("/", hello)
	e.GET("/ws", handler.WebSocket, auth.FirebaseMiddleware())
	e.GET("/ws/rooms/:id", handler.RoomWebSocket, auth.FirebaseMiddleware())
	e.GET("/users", handler.UserIndex, auth.FirebaseMiddleware())
	e.GET("/user", handler.UserShow, auth.FirebaseMiddleware())
	e.POST("/users", handler.UserCreate, auth.FirebaseMiddleware())
	e.GET("/rooms", handler.RoomIndex, auth.FirebaseMiddleware())
	e.POST("/rooms", handler.RoomCreate, auth.FirebaseMiddleware())

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

	e.Validator = &CustomValidator{validator: validator.New()}

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

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}
